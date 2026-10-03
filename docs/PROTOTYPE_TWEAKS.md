# GlutenScout Prototype Tweaks (Round 1 Feedback)

This file is a brief for an AI coding agent. It lists six changes to make to the GlutenScout prototype, based on testing with three potential customers.

## Context

GlutenScout is a clickable mock-up that helps gluten-free eaters find a specific meal they can trust. It runs on sample data with no backend. The goal of this round is to fix what confused or slowed down testers, not to add real functionality.

Test results that drive these changes:

- Two of three testers needed help finding their way through the app.
- One tester opened a restaurant with no reviews and did not understand the empty screen.
- One tester hesitated at the restaurant safety rating because its meaning was unclear.
- Two of three asked for photos (the food, the storefront, the prep area).
- Testers wanted to sort by type of food, and desserts are hard to find.
- All three eat out with people who are not gluten free.

## Ground rules

1. Read the existing code first. List the screens, components, and sample data files before changing anything, and reuse what is there.
2. Keep the current visual style: same colors, fonts, spacing, and badge treatment. Do not redesign.
3. Keep everything on sample data. No backend, no real uploads, no new accounts or logins.
4. Do not remove or rename existing features.
5. Keep it working at mobile width (375px) and keep visible keyboard focus states.
6. Make one commit per tweak, on a new branch, and open a pull request when done.
7. If something in this brief does not match the code, stop and ask instead of guessing.

## Tweaks, in priority order

### 1. Make "Find restaurants" obvious

**Why:** Two of three testers got stuck and did not know what to do next.

**Build:**
- Make "Find restaurants" the clear primary action on the first screen (largest, highest-contrast button, visible without scrolling).
- Give every screen one obvious next step, and a clear way back to the restaurant list.

**Done when:** A first-time user can go from the first screen to a chosen meal without any explanation.

### 2. Show "No reviews yet" on empty restaurants and meals

**Why:** A tester picked a restaurant with no reviews and was confused by the blank screen.

**Build:**
- In the restaurant list, show a "No reviews yet" label on any restaurant with zero reviews.
- On the detail screen, replace the blank area with a short message, for example: "No reviews yet. Be the first to review a meal here."
- Do the same for individual meals with zero reviews.
- Make sure the sample data includes at least one restaurant with no reviews so this state can be tested.

**Done when:** No screen shows an empty review area without an explanation.

### 3. Explain the safety rating

**Why:** A tester hesitated at the restaurant safety rating.

**Build:**
- Add a one-line explanation next to or under the rating wherever it appears.
- Base the wording on how the rating is actually defined in the code or sample data. If there is no definition, use: "Based on gluten-free diner reviews and how the kitchen handles prep."
- Optional: an info icon that opens a short description of what each rating level means.

**Done when:** A user can tell what the rating means without leaving the screen.

### 4. Add photos

**Why:** Two of three testers asked for photos. One wants candid photos of the real food. One judges safety by how clean the storefront and prep area look.

**Build:**
- Add a photo to each meal in the sample data.
- Add a small photo section on the restaurant detail screen with three labeled groups: Food, Storefront, Prep area.
- Use local placeholder images stored in the repo. Do not hotlink images from other websites.
- Add an "Add a photo" button that shows a "Coming soon" message. No real upload.
- Give every image alt text. Lay out the page so it still looks right when a photo is missing.

**Done when:** Each sample restaurant shows photos in all three groups, and each meal has a photo or a clean fallback.

### 5. Add filters by food type, including desserts

**Why:** One tester asked to separate salad-based from grain-based options. Another said good desserts are the hardest thing to find.

**Build:**
- Add a food type field to each meal in the sample data. Suggested values: Salads, Grain-based, Entrees, Desserts.
- Add filter chips to the restaurant list and the meal list. Selecting a chip shows only matching results.
- Include an "All" option that clears the filter.
- Show a short message when a filter has no results.
- Make sure the sample data has at least two desserts.

**Done when:** Tapping "Desserts" shows only desserts, and "All" brings everything back.

### 6. Label restaurants that work for mixed groups

**Why:** All three testers eat out with people who are not gluten free. One liked seeing options for the rest of the group.

**Build:**
- Add a "good for groups" true/false field to each restaurant in the sample data.
- Show a "Good for mixed groups" label on the restaurant card and the detail screen when true.
- Keep the existing section that shows options for people who are not gluten free.

**Done when:** A user can spot group-friendly restaurants from the list without opening each one.

## Out of scope

- Real photo uploads or user accounts
- Real restaurant data or map integration
- Payments or pricing screens
- Any change to the overall visual design

## When finished

Reply with:

1. A short list of what changed for each tweak, with the files touched.
2. Anything in this brief that could not be done, and why.
3. Steps to check each tweak by hand in the browser.
