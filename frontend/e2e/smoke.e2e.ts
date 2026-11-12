import { expect, test } from "@playwright/test";

// The journey a first-time tester takes: start, pick a city, choose a cuisine,
// open a restaurant and check its kitchen safety. If this breaks, nobody can
// use the prototype, so it blocks the merge.
test("a new visitor can find a restaurant and check kitchen safety", async ({ page }) => {
  const errors: string[] = [];
  page.on("pageerror", (e) => errors.push(e.message));

  await page.goto("/?videos=off");
  await page.getByRole("button", { name: "Celiac" }).click();

  await expect(page.getByRole("heading", { name: "Hi, I'm Scout." })).toBeVisible();
  await page.getByLabel("Your city").selectOption({ index: 1 });
  await page.getByRole("button", { name: "Mexican" }).click();

  await expect(page).toHaveURL(/\/search\?tile=mexican/);
  await expect(page.getByRole("heading", { name: /\d+ restaurants? in / })).toBeVisible();
  await page.locator(".restaurant-row").first().click();

  await expect(page).toHaveURL(/\/restaurants\/\d+/);
  await page.getByRole("button", { name: /Check kitchen safety/ }).click();
  const sheet = page.getByRole("dialog");
  await expect(sheet).toBeVisible();
  await expect(sheet.getByText("Call restaurant")).toBeVisible();
  await sheet.getByRole("button", { name: "OK" }).click();
  await expect(sheet).toBeHidden();

  expect(errors).toEqual([]);
});

test("an unknown address shows a page, not a blank screen", async ({ page }) => {
  await page.goto("/no-such-page?videos=off");
  await page.getByRole("button", { name: "Skip for now" }).click();
  await expect(page.getByRole("heading", { name: "Page not found" })).toBeVisible();
  await expect(page.getByRole("navigation").getByText("Search")).toBeVisible();
});
