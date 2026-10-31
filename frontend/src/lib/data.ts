// Sample data and the rules for reading it. No React here, so it's easy to test.
// sample.json is exported from the Go app's seed data (fictional restaurants).
import sample from "../data/sample.json";

export type GFStatus = "gf" | "cross_contact" | "contains";
export type Evidence = "confirmed" | "conflict" | "unverified";

export interface City { id: number; name: string; state: string; slug: string }
export interface Restaurant {
  id: number; cityId: number; name: string; address: string; cuisine: string;
  description: string; phone: string; goodForGroups: boolean;
}
export interface Meal {
  id: number; restaurantId: number; name: string; description: string; photo: string;
  gfStatus: GFStatus; gfNote: string; foodType: string;
}
export interface Review {
  id: number; mealId: number; author: string; body: string;
  reviewerType: string; rating: number; createdAt: string;
}
export interface SafetyProfile {
  restaurantId: number; verifiedOn: string; answersUpdatedOn: string | null; conflictNote: string;
}
export interface Precaution { restaurantId: number; label: string; status: Evidence }
// Food photos aren't here: those belong to meals, because they come from diners.
export interface Photo { restaurantId: number; kind: "restaurant" | "prep" | "menu"; url: string; alt: string; addedOn?: string }

export const cities = sample.cities as City[];
export const restaurants = sample.restaurants as Restaurant[];
export const meals = sample.meals as Meal[];
export const sampleReviews = sample.reviews as Review[];
const safetyProfiles = sample.safetyProfiles as SafetyProfile[];
const precautions = sample.precautions as Precaution[];
const photos = sample.photos as Photo[];

export const getRestaurant = (id: number) => restaurants.find((r) => r.id === id);
export const getMeal = (id: number) => meals.find((m) => m.id === id);
export const mealsFor = (restaurantId: number) => meals.filter((m) => m.restaurantId === restaurantId);
export const safetyFor = (restaurantId: number) => safetyProfiles.find((s) => s.restaurantId === restaurantId);
export const precautionsFor = (restaurantId: number) => precautions.filter((p) => p.restaurantId === restaurantId);
export const photosFor = (restaurantId: number) => photos.filter((p) => p.restaurantId === restaurantId);

export const REVIEWER_TYPES = [
  { value: "celiac", label: "Celiac" },
  { value: "strict", label: "Strict GF" },
  { value: "sensitive", label: "Gluten-sensitive" },
  { value: "restaurant", label: "Restaurant" },
];
export const reviewerLabel = (value: string) => REVIEWER_TYPES.find((t) => t.value === value)?.label ?? "";

// Dates are stored without a time, so read them as UTC: otherwise "Sep 12"
// can show as "Sep 11" west of Greenwich.
export function shortDate(iso: string, withYear = false): string {
  return new Date(iso).toLocaleDateString("en-US", {
    month: "short", day: "numeric", timeZone: "UTC", ...(withYear ? { year: "numeric" } : {}),
  });
}

// The status shown on the safety button, search results and the sheet. Warnings
// have to show on the page, because anyone can skip a button. Copy rule: name
// what's known, and never call a restaurant "safe".
export type SafetyLevel = Evidence | "none";
export interface SafetyStatus { level: SafetyLevel; label: string; summary: string; verifiedOn?: string }

export function safetyStatus(restaurantId: number): SafetyStatus {
  const profile = safetyFor(restaurantId);
  if (!profile) {
    return { level: "none", label: "No safety information yet", summary: "Nobody has reported on this kitchen yet. Call ahead and ask before you order." };
  }
  const list = precautionsFor(restaurantId);
  const { verifiedOn } = profile;
  if (profile.conflictNote || list.some((p) => p.status === "conflict")) {
    return { level: "conflict", label: "Reports disagree", summary: profile.conflictNote || "Reports don't agree about the practices below.", verifiedOn };
  }
  if (list.length === 0 || list.some((p) => p.status === "unverified")) {
    return { level: "unverified", label: "Some details unconfirmed", summary: "Some practices haven't been backed up by diners, or the information is out of date.", verifiedOn };
  }
  return { level: "confirmed", label: "Practices confirmed", summary: "The restaurant and recent gluten-free diners agree on the practices below.", verifiedOn };
}

// Dashboard tiles. A plain list so tiles can change after testing: the set
// is a guess until interviews and Clarity data back it up. A cuisine tile
// only groups restaurants; it says nothing about how safe that cuisine is.
export interface Tile {
  id: string;
  label: string;
  cuisines?: string[];
  foodType?: string;
  groups?: boolean;
}
export const TILES: Tile[] = [
  { id: "american", label: "American", cuisines: ["American", "Burgers", "Steakhouse", "Barbecue", "Fried chicken", "Cafe"] },
  { id: "italian", label: "Italian", cuisines: ["Italian", "Pizza"] },
  { id: "mexican", label: "Mexican", cuisines: ["Mexican"] },
  { id: "thai", label: "Thai", cuisines: ["Thai"] },
  { id: "desserts", label: "Desserts", foodType: "dessert" },
  { id: "groups", label: "Good for mixed groups", groups: true },
];
export const getTile = (id: string | null) => TILES.find((t) => t.id === id);

const gfMealsFor = (restaurantId: number) => mealsFor(restaurantId).filter((m) => m.gfStatus === "gf");

export interface SearchResult {
  restaurant: Restaurant;
  gfCount: number;      // GF meals, or only GF desserts under the Desserts tile
  countNoun: string;    // "meal" or "dessert"
  matchedMeals: string[]; // GF meals whose name matched the typed search
}

// Finds restaurants in a city (or everywhere when cityId is null) by tile
// and/or typed text. Typed text only matches meals that are GF options, so a
// search for "pizza" never leads with a meal the diner can't eat.
export function search(opts: { cityId: number | null; tile?: Tile; query?: string }): SearchResult[] {
  const q = (opts.query ?? "").trim().toLowerCase();
  const tile = opts.tile;
  const out: SearchResult[] = [];
  for (const r of restaurants) {
    if (opts.cityId !== null && r.cityId !== opts.cityId) continue;
    const gf = gfMealsFor(r.id);
    const typed = tile?.foodType ? gf.filter((m) => m.foodType === tile.foodType) : gf;
    if (tile?.cuisines && !tile.cuisines.includes(r.cuisine)) continue;
    if (tile?.groups && !r.goodForGroups) continue;
    if (tile?.foodType && typed.length === 0) continue;
    const matchedMeals = q ? typed.filter((m) => m.name.toLowerCase().includes(q)).map((m) => m.name) : [];
    if (q && matchedMeals.length === 0 &&
        !r.name.toLowerCase().includes(q) && !r.cuisine.toLowerCase().includes(q)) continue;
    out.push({ restaurant: r, gfCount: typed.length, countNoun: tile?.foodType ?? "meal", matchedMeals });
  }
  return out.sort((a, b) => a.restaurant.name.localeCompare(b.restaurant.name));
}

// How the search URL changes. Typing searches the whole city: someone who
// picked American and then types "tacos" has changed their mind, so the tile
// must not quietly narrow the search to nothing.
export function nextSearchParams(current: URLSearchParams, next: { tile?: string | null; q?: string }): URLSearchParams {
  const p = new URLSearchParams(current);
  if (next.q !== undefined) {
    if (next.q) {
      p.set("q", next.q);
      p.delete("tile");
    } else {
      p.delete("q");
    }
  }
  if (next.tile !== undefined) {
    if (next.tile) p.set("tile", next.tile);
    else p.delete("tile");
  }
  return p;
}

export interface RatingSummary { average: number; ratingCount: number; reviewCount: number }

// Restaurant-written reviews carry no rating (restaurants don't rate their own
// food), so they count as reviews but not ratings.
export function summarize(reviews: Review[]): RatingSummary {
  const rated = reviews.filter((r) => r.rating > 0);
  const total = rated.reduce((sum, r) => sum + r.rating, 0);
  return { average: rated.length ? total / rated.length : 0, ratingCount: rated.length, reviewCount: reviews.length };
}

export interface ReviewInput { author: string; reviewerType: string; rating: number; body: string }
export type ReviewErrors = Partial<Record<"reviewerType" | "rating" | "body", string>>;

export function validateReview(input: ReviewInput): ReviewErrors {
  const errors: ReviewErrors = {};
  const body = input.body.trim();
  if (!input.reviewerType) errors.reviewerType = "Choose the option that describes you.";
  if (input.reviewerType && input.reviewerType !== "restaurant" && !input.rating) {
    errors.rating = "Choose a food rating from 1 to 5.";
  }
  if (!body) errors.body = "Tell others how it went.";
  else if (body.length > 1000) errors.body = "Keep it to 1000 characters.";
  return errors;
}
