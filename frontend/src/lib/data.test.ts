import { describe, expect, it } from "vitest";
import { getTile, safetyStatus, search, summarize, TILES, validateReview, type Review } from "./data";

const PROVO = 1;
const SLC = 2;
const names = (r: ReturnType<typeof search>) => r.map((x) => x.restaurant.name);

describe("safetyStatus", () => {
  it("shows a disagreement on the button", () => {
    expect(safetyStatus(1)).toEqual({ level: "conflict", text: "Reports disagree · Verified Sep 12" });
  });
  it("says so when a restaurant has no profile", () => {
    expect(safetyStatus(10).level).toBe("none");
  });
  it("never uses the word safe", () => {
    for (let id = 1; id <= 13; id++) expect(safetyStatus(id).text).not.toMatch(/safe/i);
  });
});

describe("search", () => {
  it("counts only desserts under the Desserts tile", () => {
    const r = search({ cityId: PROVO, tile: getTile("desserts") });
    expect(names(r)).toEqual(["Maple & Sage Kitchen"]);
    expect(r[0]).toMatchObject({ gfCount: 1, countNoun: "dessert" });
  });
  it("leaves out a restaurant whose only dessert contains gluten", () => {
    expect(names(search({ cityId: PROVO, tile: getTile("desserts") }))).not.toContain("Cluck Shack");
  });
  it("groups pizza under Italian", () => {
    expect(names(search({ cityId: PROVO, tile: getTile("italian") }))).toEqual(["Canyon Hearth Pizza"]);
  });
  it("returns nothing for a tile the city doesn't have", () => {
    expect(search({ cityId: PROVO, tile: getTile("thai") })).toEqual([]);
    expect(names(search({ cityId: SLC, tile: getTile("thai") }))).toEqual(["Saltbox Thai"]);
  });
  it("shows only group-friendly restaurants under Mixed groups", () => {
    const r = search({ cityId: PROVO, tile: getTile("groups") });
    expect(r.length).toBeGreaterThan(0);
    expect(r.every((x) => x.restaurant.goodForGroups)).toBe(true);
  });
  it("finds typed text in names, cuisines and GF meals, and reports the meal", () => {
    expect(names(search({ cityId: PROVO, query: "lotus" }))).toEqual(["Little Lotus Pho"]);
    const cake = search({ cityId: null, query: "chocolate cake" });
    expect(cake[0].matchedMeals).toEqual(["Flourless Chocolate Cake"]);
  });
  it("doesn't match typed text against meals containing gluten", () => {
    expect(names(search({ cityId: PROVO, query: "chocolate chip" }))).toEqual([]);
  });
  it("has six tiles", () => expect(TILES).toHaveLength(6));
});

describe("summarize", () => {
  const review = (rating: number) => ({ rating }) as Review;
  it("leaves unrated reviews out of the average", () => {
    expect(summarize([review(4), review(5), review(0)])).toEqual({ average: 4.5, ratingCount: 2, reviewCount: 3 });
  });
});

describe("validateReview", () => {
  const good = { author: "", reviewerType: "celiac", rating: 5, body: "Great" };
  it("accepts a complete review", () => expect(validateReview(good)).toEqual({}));
  it("requires a reviewer type, rating and text", () => {
    expect(Object.keys(validateReview({ author: "", reviewerType: "", rating: 0, body: " " })).sort()).toEqual(["body", "reviewerType"]);
    expect(validateReview({ ...good, rating: 0 }).rating).toBeTruthy();
  });
  it("doesn't ask restaurants to rate their own food", () => {
    expect(validateReview({ ...good, reviewerType: "restaurant", rating: 0 })).toEqual({});
  });
  it("limits length", () => expect(validateReview({ ...good, body: "x".repeat(1001) }).body).toBeTruthy());
});
