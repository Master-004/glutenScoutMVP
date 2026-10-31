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
- Dashboard tiles are the `TILES` list in `src/lib/data.ts`. They're a starting guess to validate.
