import { Link } from "react-router-dom";
import { clarityEnabled } from "../lib/clarity";
import { RestaurantRow } from "../components/RestaurantRow";
import { getRestaurant, REVIEWER_TYPES, search } from "../lib/data";
import { useStore, type Theme } from "../lib/store";

export function Saved() {
  const { savedRestaurants } = useStore();
  const saved = savedRestaurants.map(getRestaurant).filter((r) => r !== undefined);
  return (
    <>
      <h1 className="page-title">Saved</h1>
      {saved.length ? (
        <ul className="row-list">
          {saved.map((r) => (
            <li key={r.id}><RestaurantRow result={search({ cityId: r.cityId }).find((x) => x.restaurant.id === r.id)!} /></li>
          ))}
        </ul>
      ) : (
        <p className="empty">Nothing saved yet. <Link to="/search">Find a place you trust</Link> and tap Save to keep it here.</p>
      )}
    </>
  );
}

const THEMES: { value: Theme; label: string }[] = [
  { value: "auto", label: "Auto" }, { value: "light", label: "Light" }, { value: "dark", label: "Dark" },
];

// A placeholder until accounts exist. It holds the settings the old header had
// and the switch for comparing the app with and without videos.
export function Profile() {
  const { theme, setTheme, videos, setVideos, myReviews, savedRestaurants, reset, reviewerType, setReviewerType, restart, recording, setRecording } = useStore();
  return (
    <>
      <h1 className="page-title">Profile</h1>
      <p>Accounts aren't part of this prototype. Everything you do here stays in this browser.</p>
      <p className="meta" data-clarity-mask="true">{savedRestaurants.length} saved · {myReviews.length} {myReviews.length === 1 ? "review" : "reviews"} written</p>

      <fieldset className="field" data-clarity-mask="true">
        <legend>Which describes you?</legend>
        <div className="choices">
          {REVIEWER_TYPES.map((t) => (
            <label key={t.value} className="choice">
              <input type="radio" name="profileType" checked={reviewerType === t.value} onChange={() => setReviewerType(t.value)} />
              <span>{t.label}</span>
            </label>
          ))}
        </div>
        <p className="hint">Shown on your reviews.</p>
      </fieldset>

      <fieldset className="field">
        <legend>Theme</legend>
        <div className="choices">
          {THEMES.map((t) => (
            <label key={t.value} className="choice">
              <input type="radio" name="theme" checked={theme === t.value} onChange={() => setTheme(t.value)} />
              <span>{t.label}</span>
            </label>
          ))}
        </div>
      </fieldset>
      <label className="toggle">
        <input type="checkbox" checked={videos} onChange={(e) => setVideos(e.target.checked)} />
        <span>Show Scout videos</span>
      </label>
      {clarityEnabled && (
        <label className="toggle">
          <input type="checkbox" checked={recording} onChange={(e) => setRecording(e.target.checked)} />
          <span>Let GlutenScout record how I use the app</span>
        </label>
      )}
      <div className="profile-actions">
        <button type="button" className="button-secondary" onClick={reset}>Clear my saved places and reviews</button>
        <button type="button" className="button-secondary" onClick={restart}>See the start screen again</button>
      </div>
      <p className="profile-privacy"><Link to="/privacy">Privacy</Link></p>
    </>
  );
}
