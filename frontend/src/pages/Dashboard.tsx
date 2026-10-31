import { useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { CityPicker } from "../components/CityPicker";
import { RestaurantRow } from "../components/RestaurantRow";
import { cities, getRestaurant, search, TILES, type Tile } from "../lib/data";
import { trackEvent } from "../lib/clarity";
import { useStore } from "../lib/store";

export function Dashboard() {
  const { citySlug, savedRestaurants } = useStore();
  const navigate = useNavigate();
  const picker = useRef<HTMLSelectElement>(null);
  const [needCity, setNeedCity] = useState(false);
  const city = cities.find((c) => c.slug === citySlug);
  const saved = savedRestaurants.map(getRestaurant).filter((r) => r !== undefined);

  function openTile(tile: Tile) {
    if (!city) {
      setNeedCity(true);
      picker.current?.focus();
      return;
    }
    trackEvent(`tile_${tile.id}`);
    navigate(`/search?tile=${tile.id}`);
  }

  return (
    <>
      <section className="identity">
        <img className="identity-photo" src="/images/scout.jpg" alt="Scout, the GlutenScout guide" />
        <div>
          <h1 className="identity-title">Hi, I'm Scout.</h1>
          <p className="identity-text">Let's find a meal you can feel good about.</p>
        </div>
      </section>

      <CityPicker ref={picker} prominent error={needCity && !city ? "Choose your city first." : undefined} />

      <section aria-labelledby="tiles-heading">
        <h2 id="tiles-heading" className="section-header">What sounds good?</h2>
        {/* Quieter until there's a city: the picker above is the next step. */}
        <ul className={`tile-grid${city ? "" : " tile-grid-waiting"}`}>
          {TILES.map((tile) => (
            <li key={tile.id}>
              <button type="button" className="tile" onClick={() => openTile(tile)}>{tile.label}</button>
            </li>
          ))}
        </ul>
      </section>

      <section aria-labelledby="saved-heading">
        <h2 id="saved-heading" className="section-header">Your trusted places</h2>
        {saved.length ? (
          <ul className="row-list">
            {saved.map((r) => (
              <li key={r.id}>
                <RestaurantRow result={search({ cityId: r.cityId }).find((x) => x.restaurant.id === r.id)!} />
              </li>
            ))}
          </ul>
        ) : (
          <p className="empty">
            Nothing saved yet. <Link to="/search">Find a place you trust</Link> and save it to keep it here.
          </p>
        )}
      </section>
    </>
  );
}
