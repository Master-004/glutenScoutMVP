# GlutenScout Design System

> Source of truth: the GlutenScout design system (https://claude.ai/artifact/PWYaC62nFk9q7Ghd9rE5kM).
> This file is the in-repo copy for people and AI agents. If the two disagree, the design system wins — update this file to match.

GlutenScout helps gluten-sensitive and celiac diners find restaurants nearby, inspect dated safety evidence, and prepare questions before they visit. Every screen answers: **"Can I eat here safely, and how sure are we?"** — name, safety profile, last verified date, precaution badges, then ratings and reviews.

---

## Rules for AI agents (read first)

1. Use only the tokens below — never hard-code a hex value, font size or spacing that isn't listed. Reference tokens by name (e.g. `brand`, `space-4`).
2. Support light and dark themes on every screen.
3. Never write "Safe", "celiac-safe", "100% gluten-free" or "guaranteed" in UI copy. Show evidence, dates and sources instead; the diner decides.
4. Every safety claim shows a date. Every report shows its reviewer type.
5. Never signal safety with color alone — always icon + words + color.
6. Stars rate food and experience only, never gluten safety.
7. One green (`brand`) button per screen.
8. Text must meet 4.5:1 contrast (3:1 for 24px+ text, icons, borders).
9. Reuse the components below; don't invent new variants without updating the design system.

---

## Color tokens

| Token | Light | Dark | Use |
|---|---|---|---|
| `canvas` | #ffffff | #161616 | Page background (never pure black) |
| `surface` | #f5f5f5 | #1e1e1e | Search bar, pills, chips, sheets |
| `surface-raised` | #ffffff | #2a2a2a | Second elevation |
| `band` | #f0f0f0 | #0c0c0c | 8px section bands between page sections |
| `divider` | #e0e0e0 | #303030 | Hairlines, 1px borders |
| `pressed` | #ededed | #2a2a2a | Pressed background, secondary buttons |
| `ink` | #2b2b2b | #e6e6e6 | Primary text |
| `ink-secondary` | #6e6e6e | #a8a8a8 | Meta text, dates, inactive icons |
| `ink-tertiary` | #9a9a9a | #6e6e6e | Placeholders only (fails 4.5:1) |
| `on-brand` | #ffffff | #ffffff | Text/icons on `brand` |
| `brand` | #0f7b4a | #0f7b4a | Primary CTA fill, selected chip |
| `brand-pressed` | #0b5e38 | #0b5e38 | Pressed brand fill |
| `brand-deep` | #0a4d2f | #0a4d2f | Wordmark, avatar monogram, reviewer-type tag text (light) |
| `brand-ink` | #0f7b4a | #4cc38a | Green as text/icon: active tab, focus ring, chosen actions |
| `brand-tint` | #e3f4ea | #13301f | Quiet green background, reviewer-type tag |
| `brand-disabled` | #b9dec8 | #1e3a2a | Disabled primary button |
| `star` | #e06a00 | #e06a00 | Filled rating star (food rating only) |
| `star-empty` | #e3e3e0 | #3a3a3a | Empty star |
| `status-open` | = `brand-ink` | = `brand-ink` | "Open" |
| `status-closed` | #d62f2f | #f06a6a | "Closed" and errors — the only red |
| `status-warning` | #b45309 | #e08b2d | "Closes soon", warnings |
| `link` | #0073bb | #4fa3d9 | Text links |
| `evidence-confirmed` | = `brand-ink` | = `brand-ink` | Confirmed precaution (check icon) |
| `evidence-confirmed-bg` | = `brand-tint` | = `brand-tint` | Background for confirmed |
| `evidence-conflict` | = `status-warning` | = `status-warning` | Reports disagree (triangle icon) |
| `evidence-conflict-bg` | #fff6ec | #33240f | Background for conflict |
| `evidence-unverified` | = `ink-secondary` | = `ink-secondary` | Restaurant-only or out of date (clock icon) |
| `evidence-unverified-bg` | = `surface` | = `surface` | Background for unverified |
| `scrim` | rgba(0,0,0,0.4) | same | Behind sheets |
| `photo-button` | rgba(0,0,0,0.45) | same | Floating buttons on photos |

## Typography

Font: **Open Sans** (Google Fonts, weights 400/600/700/800). Fallback: `-apple-system, BlinkMacSystemFont, "SF Pro Text", "Segoe UI", Roboto, Helvetica, Arial, sans-serif`.

| Style | Size | Weight | Line height | Use |
|---|---|---|---|---|
| `screen-title` | 32 | 800 | 1.2 | Navigation titles |
| `business-name` | 24 | 800 | 1.2 | Restaurant name |
| `section-header` | 20 | 700 | 1.25 | "Safety profile", "Reports" |
| `card-header` | 16 | 800 | 1.3 | Sub-sections |
| `emphasis` | 15 | 700 | 1.3 | Ratings, key stats |
| `body` | 14 | 400 | 1.55 | Review text — never bold |
| `reviewer-name` | 14 | 700 | 1.3 | Report author |
| `meta` | 13 | 600 | 1.4 | Dates, locations, counts |
| `caption` | 12 | 600 | 1.4 | Peer-review actions |
| `button` | 15 | 700 | 1.0 | Primary buttons |
| `chip` | 13 | 600 | 1.0 | Chips, badges (fixed size) |
| `tab-label` | 10 | 600 | 1.0 | Tab bar (fixed size) |

## Spacing, radius, shadow

- Spacing (4px base): `space-1` 4 · `space-2` 8 · `space-3` 12 · `space-4` 16 · `space-5` 20 · `space-6` 24 · `space-8` 32 · `space-10` 40 · `space-12` 48 · `space-16` 64. Content inset = `space-4`.
- Radius: `radius-sm` 4 (photos) · `radius-md` 8 (buttons, inputs, safety profile) · `radius-lg` 12 (sheets) · `radius-pill` 500 (chips, badges) · `radius-circle` 50% (avatars).
- Shadow: flat by default. `shadow-raised` `0 1px 4px rgba(0,0,0,0.10)` (sticky search bar) · `shadow-sheet` `0 -4px 20px rgba(0,0,0,0.16)` (bottom sheets).
- Touch targets ≥ 44px. Motion 120–300ms ease-out; nothing flashes on safety information.

---

## Components

**Button**
- Primary: `brand` fill, `on-brand` text, `button` style, `radius-md`, 13×26 padding, min height 44. Pressed `brand-pressed` + scale 0.98. Disabled `brand-disabled`, text 60% opacity.
- Secondary: transparent, 1px `divider`, `ink` 14/600, 11×20 padding.
- Ghost icon (Save / Share / Directions): 56×44, 1px `divider`, `ink-secondary`; selected → `brand-ink`.

**PrecautionBadge** — a pill naming one kitchen practice and its evidence state. Confirmed = check + `evidence-confirmed` on `evidence-confirmed-bg`; Conflict = triangle + `evidence-conflict` on `evidence-conflict-bg`; Unverified = clock + `evidence-unverified` on `evidence-unverified-bg`. `chip` type, `radius-pill`, 6×12 padding. Labels name practices ("Separate prep area", "Gloves changed", "Shared fryer?"), never verdicts.

**SafetyProfile** — the combined summary on a restaurant page. `section-header` "Safety profile" → "Last verified {date}" in `meta` `ink-secondary` with clock (always shown) → PrecautionBadges row → source lines ("Restaurant answers · updated {date}", "{n} celiac / strict reports · newest {date}") → conflict notice on `evidence-conflict-bg` in words when reports disagree → secondary "Call ahead" button when anything is conflicting or unverified. 1px `divider`, `radius-md`, `space-4` padding.

**ReviewCard (report)** — no border; separated by 1px `divider`, `space-4` vertical padding. 40px circle avatar (monogram on `brand-deep`). Name `reviewer-name` + required reviewer-type tag (Celiac / Strict GF / Gluten-sensitive / Restaurant) in `brand-deep` on `brand-tint`. Location + report count and visit date in `caption` `ink-secondary`. Stars (food). Body in `body`, "Read more" after 4 lines. Peer actions: Helpful · Matches my visit · Report conflict (`caption`, chosen → `brand-ink`).

**StarRating** — five stars, `star` / `star-empty`, 2px gap, half-star steps only. Sizes 22 / 18 / 15 / 14. Always paired with the count: "4.5 · 3,812 reviews".

**FilterChip** — `surface` fill, 1px `divider`, `radius-pill`, `chip` text, 9×16 padding, optional 15px glyph in `brand-ink`. Selected: `brand` fill, `on-brand` text. Horizontal scroll, 10px gap. Examples: Celiac reports · Dedicated fryer · Within 2 mi · Verified this month.

**Page layout** — restaurant page order: photo header → name, stars, hours → `band` → SafetyProfile → `band` → call-ahead → `band` → reports. Bands (8px `band`), not bordered cards, separate sections.

---

## Voice

- Evidence, not guarantees: "Separate prep area — confirmed by 4 celiac reviewers."
- Dates on everything: "Last verified Sep 12, 2026."
- Conflicts in words: "2 reports disagree about the fryer."
- Calm, short, plain. No alarm words, no medical advice. When evidence is thin, point to the call-ahead script.
