# GlutenScout React prototype

A clickable, mobile-first prototype for the next round of user testing. It runs on sample data in `src/data/sample.json` (exported from the Go app's fictional seed data), with no backend. Saved places and new reviews stay in the visitor's browser.

```sh
npm install
npm run dev     # http://localhost:5173
npm test
npm run build
```

- `?videos=off` on any test link turns Scout's videos off for that visitor, for comparing versions. Profile has the same switch.
- Set `VITE_CLARITY_ID` in `.env` to turn on Microsoft Clarity. Tell testers they're being recorded first.
- `node scripts/extend-sample.mjs` rebuilds the extra fictional sample data. It's only in this prototype, not the Go seed.
- Dashboard tiles are the `TILES` list in `src/lib/data.ts`. They're a starting guess to validate.

## Known issues and later ideas

- **Desktop nav:** the top nav on wider screens looks odd. Mobile is the priority, so fix this later.
- **Photos and accounts:** once accounts exist, tie each photo to the account that added it. Restaurants hold restaurant, prep area and menu photos; meals hold diners' food photos.
- **Uploads:** "Add a photo" only shows a coming-soon message for now.
