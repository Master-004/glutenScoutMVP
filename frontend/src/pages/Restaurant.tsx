import { useState } from "react";
import { Link, useNavigate, useParams, useSearchParams, useLocation } from "react-router-dom";
import { SafetyButton, SafetySheet } from "../components/SafetySheet";
import { GroupsLabel, Img, MealRow, Rating, Sheet, useRestaurantSummary } from "../components/ui";
import { getRestaurant, mealsFor, photosFor, shortDate, type Meal, type Restaurant as R } from "../lib/data";
import { trackEvent } from "../lib/clarity";
import { useStore } from "../lib/store";
import { NotFound } from "./NotFound";

export function BackLink({ fallback }: { fallback: string }) {
  const navigate = useNavigate();
  const location = useLocation();
  // "default" means this is the first page of the visit, so there's no Back.
  return (
    <button
      type="button" className="back-link"
      onClick={() => (location.key === "default" ? navigate(fallback) : navigate(-1))}
    >
      ← Back
    </button>
  );
}

const PHOTO_GROUPS = [
  { kind: "restaurant", label: "Restaurant" },
  { kind: "prep", label: "Prep areas" },
  { kind: "menu", label: "Menus" },
] as const;

function PhotoSheet({ restaurant, open, onClose }: { restaurant: R; open: boolean; onClose(): void }) {
  const [comingSoon, setComingSoon] = useState(false);
  const photos = photosFor(restaurant.id);
  return (
    <Sheet open={open} onClose={onClose} title="Photos">
      <ul className="photo-groups">
        {PHOTO_GROUPS.map((g) => {
          const photo = photos.find((p) => p.kind === g.kind);
          return (
            <li key={g.kind}>
              <h3 className="photo-label">{g.label}</h3>
              {photo ? <Img src={photo.url} alt={photo.alt} /> : <div className="photo-missing">No photo yet</div>}
              {/* Menus change, so a menu photo always says how old it is. */}
              {photo?.addedOn && <p className="photo-date">Added {shortDate(photo.addedOn)}</p>}
            </li>
          );
        })}
      </ul>
      <p className="meta">Illustrations for now. Real photos are coming.</p>
      <button type="button" className="button-secondary" onClick={() => setComingSoon(true)}>Add a photo</button>
      {comingSoon && <p role="status" className="meta">Coming soon: you'll be able to add photos of the restaurant, prep areas and menus.</p>}
    </Sheet>
  );
}

function MealList({ meals }: { meals: Meal[] }) {
  const { reviewsFor } = useStore();
  return (
    <ul className="meal-list">
      {meals.map((m) => <li key={m.id}><MealRow meal={m} foodReviews={reviewsFor(m.id)} /></li>)}
    </ul>
  );
}

const NOUNS: Record<string, string> = { dessert: "desserts" };

export function Restaurant() {
  const { id } = useParams();
  const [params, setParams] = useSearchParams();
  const { savedRestaurants, toggleSaved } = useStore();
  const [sheet, setSheet] = useState<"safety" | "photos" | null>(null);
  const restaurant = getRestaurant(Number(id));
  const all = restaurant ? mealsFor(restaurant.id) : [];
  const summary = useRestaurantSummary(all.map((m) => m.id));
  if (!restaurant) return <NotFound />;

  const food = params.get("food");
  const meals = food ? all.filter((m) => m.foodType === food) : all;
  const gf = meals.filter((m) => m.gfStatus === "gf");
  const cross = meals.filter((m) => m.gfStatus === "cross_contact");
  const contains = meals.filter((m) => m.gfStatus === "contains");
  const photo = photosFor(restaurant.id).find((p) => p.kind === "restaurant");
  const saved = savedRestaurants.includes(restaurant.id);

  return (
    <>
      <BackLink fallback="/search" />
      <section className="identity">
        <button type="button" className="identity-photo-button" onClick={() => { trackEvent("photos_opened"); setSheet("photos"); }} aria-label={`See photos of ${restaurant.name}`}>
          <Img className="identity-photo" src={photo?.url ?? ""} alt="" />
          <span className="identity-photo-label">Photos</span>
        </button>
        <div className="identity-main">
          <h1 className="identity-title">{restaurant.name}</h1>
          <Rating summary={summary} />
          <p className="meta">{restaurant.cuisine} · {restaurant.address}</p>
          <GroupsLabel restaurant={restaurant} />
          <button type="button" className="save-button" aria-pressed={saved} onClick={() => { if (!saved) trackEvent("restaurant_saved"); toggleSaved(restaurant.id); }}>
            {saved ? "Saved" : "Save"}
          </button>
        </div>
      </section>
      {summary.reviewCount === 0 && <p className="empty-note">No reviews yet. Be the first to review a meal here.</p>}

      <SafetyButton restaurantId={restaurant.id} onClick={() => { trackEvent("safety_sheet_opened"); setSheet("safety"); }} />
      <SafetySheet restaurant={restaurant} open={sheet === "safety"} onClose={() => setSheet(null)} />
      <PhotoSheet restaurant={restaurant} open={sheet === "photos"} onClose={() => setSheet(null)} />

      <section aria-labelledby="gf-heading">
        <h2 id="gf-heading" className="section-header">{food ? `Gluten-free ${NOUNS[food] ?? food}` : "Gluten-free options"}</h2>
        {food && (
          <p className="meta filter-note">
            Showing {NOUNS[food] ?? food} only. <button type="button" className="link-button" onClick={() => setParams({}, { replace: true })}>Show all meals</button>
          </p>
        )}
        {gf.length ? <MealList meals={gf} /> : <p className="empty">No gluten-free options listed here yet.</p>}
      </section>

      {cross.length > 0 && (
        <section aria-labelledby="cross-heading">
          <h2 id="cross-heading" className="section-header">Cross-contact risk</h2>
          <p className="meta">Made without gluten ingredients, but prepared with shared equipment. Ask before ordering.</p>
          <MealList meals={cross} />
        </section>
      )}

      {contains.length > 0 && (
        // Folded, but still there for the non-GF people in the group.
        <details className="contains-group">
          <summary className="section-header">Contains gluten ({contains.length})</summary>
          <MealList meals={contains} />
        </details>
      )}

      {!food && all.length === 0 && <p className="empty"><Link to="/search">Back to search</Link></p>}
    </>
  );
}
