# GlutenScout React prototype

A clickable, mobile-first prototype for the next round of user testing. It runs on sample data in `src/data/sample.json` (exported from the Go app's fictional seed data), with no backend. Saved places and new reviews stay in the visitor's browser.

```sh
npm install
npm run dev     # http://localhost:5173
npm test
npm run build
```

- `?videos=off` on any test link turns Scout's videos off for that visitor, for comparing versions. Profile has the same switch.
- `node scripts/extend-sample.mjs` rebuilds the extra fictional sample data. It's only in this prototype, not the Go seed.
- Dashboard tiles are the `TILES` list in `src/lib/data.ts`. They're a starting guess to validate.

## Deploying on Vercel

Set the project's root directory to `frontend`. `vercel.json` sends every page to the app, so links like `/restaurants/2` work on refresh.

## Checks

Every pull request runs `.github/workflows/frontend-ci.yml`. Run the same checks locally:

```sh
npm run lint        # ESLint
npm run typecheck   # TypeScript
npm test            # Vitest unit tests
npm run build
npm run smoke       # Playwright: the first-visit journey on a phone-sized screen
```

The first time, run `npx playwright install chromium` before `npm run smoke`.

## Error monitoring (Sentry)

Sentry stays off until `VITE_SENTRY_DSN` is set. When it's on:

- Crashes show a "Something went wrong" page and send a report to the `glutenscout-frontend` project.
- Reports leave out user details, cookies, IP addresses and search text (`?q=`). There's no Sentry Session Replay.
- An email alert goes out for new or returning errors in production. An uptime check looks at the live site every 5 minutes.

**In Vercel:**

1. Add `VITE_SENTRY_DSN` (copy it from Sentry: Settings → Projects → glutenscout-frontend → Client Keys). It isn't a secret.
2. For readable stack traces, add `SENTRY_AUTH_TOKEN`, an organization token from Sentry (Settings → Auth Tokens). Mark it as sensitive. Builds that have it upload source maps to Sentry, then delete them so they're never published.
3. Redeploy. Vercel's commit and environment are attached to each report automatically.

For local testing, put the DSN in `frontend/.env.local`. That file isn't committed.

## Turning on Clarity

Clarity stays off until `VITE_CLARITY_ID` is set. With no ID, the script never loads and no notice appears.

**Before turning it on:**

1. Read the [Clarity terms](https://clarity.microsoft.com/terms). This app asks about celiac disease, which is health-related, so confirm it's allowed.
2. Review the draft wording in `src/pages/Privacy.tsx` and the notices on the start screen and footer.
3. Make sure testers are 18 or older. Clarity shouldn't be used on anything aimed at under-18s.
4. Tell testers they're being recorded when you invite them.

**To turn it on:**

1. Create a project at [clarity.microsoft.com](https://clarity.microsoft.com) and copy its project ID.
2. In the Clarity project's masking settings, keep **Balanced** or stricter. The code masks reviews, forms and dietary choices on top of that.
3. In Vercel, add `VITE_CLARITY_ID` as an environment variable, then redeploy. The ID is built in at deploy time.
4. Run one test session and watch the recording. Confirm that review text, the review form and the Profile dietary choice are hidden. Masking changes can't be applied to past recordings, so do this before inviting testers.

**What it records:**

- It starts only after the start screen, so the dietary choice there is never recorded. People can turn it off in Profile.
- Tag: `videos` (on or off), to compare versions with `?videos=off`.
- Events: `tile_<id>`, `meal_opened`, `safety_sheet_opened`, `call_tapped`, `photos_opened`, `restaurant_saved`, `review_posted`, `video_skipped`.
- Never add a reviewer type to a tag or event.
- Page addresses are recorded, and they include search text (`?q=`).

## Known issues and later ideas

- **Desktop nav:** the top nav on wider screens looks odd. Mobile is the priority, so fix this later.
- **Photos and accounts:** once accounts exist, tie each photo to the account that added it. Restaurants hold restaurant, prep area and menu photos; meals hold diners' food photos.
- **Uploads:** "Add a photo" only shows a coming-soon message for now.
