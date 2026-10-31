import { useSearchParams } from "react-router-dom";
import { CityPicker } from "../components/CityPicker";
import { RestaurantRow } from "../components/RestaurantRow";
import { cities, getTile, nextSearchParams, search, TILES } from "../lib/data";
import { useStore } from "../lib/store";

// The tile and search text live in the URL, so Back returns to these results.
export function Search() {
  const { citySlug } = useStore();
  const [params, setParams] = useSearchParams();
  const tile = getTile(params.get("tile"));
  const query = params.get("q") ?? "";
  const city = cities.find((c) => c.slug === citySlug);
  const results = search({ cityId: city?.id ?? null, tile, query });

  function change(next: { tile?: string | null; q?: string }) {
    setParams(nextSearchParams(params, next), { replace: next.q !== undefined });
  }

  const where = city ? ` in ${city.name}` : "";
  const everywhere = city ? `all of ${city.name}` : "every city";
  return (
    <>
      <h1 className="page-title">Search</h1>
      <div className="field">
        <label htmlFor="q" className="visually-hidden">Search restaurants, cuisines or meals</label>
        <input
          id="q" type="search" value={query} onChange={(e) => change({ q: e.target.value })}
          placeholder="Restaurant, cuisine or meal" autoComplete="off"
        />
      </div>
      <CityPicker />
      <nav className="filter-chips" aria-label="Filter results">
        <button type="button" className="filter-chip" aria-pressed={!tile} onClick={() => change({ tile: null })}>All</button>
        {TILES.map((t) => (
          <button key={t.id} type="button" className="filter-chip" aria-pressed={tile?.id === t.id} onClick={() => change({ tile: t.id })}>
            {t.label}
          </button>
        ))}
      </nav>

      {results.length > 0 ? (
        <>
          <h2 className="results-heading" aria-live="polite">
            {results.length} {results.length === 1 ? "restaurant" : "restaurants"}{where}
          </h2>
          <ul className="row-list">
            {results.map((r) => (
              <li key={r.restaurant.id}><RestaurantRow result={r} foodType={tile?.foodType} /></li>
            ))}
          </ul>
        </>
      ) : tile && query ? (
        // Both a tile and text: offer the way out instead of a dead end.
        <p className="empty" aria-live="polite">
          No "{query}" under {tile.label}{where}.{" "}
          <button type="button" className="link-button" onClick={() => change({ tile: null })}>Search {everywhere} for "{query}"</button>
        </p>
      ) : (
        <p className="empty" aria-live="polite">
          Nothing matches{tile ? ` ${tile.label}` : ""}{query ? ` "${query}"` : ""}{where} yet.{" "}
          <button type="button" className="link-button" onClick={() => setParams({})}>Show all restaurants</button>
        </p>
      )}
    </>
  );
}
