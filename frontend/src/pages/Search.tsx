import { useSearchParams } from "react-router-dom";
import { CityPicker } from "../components/CityPicker";
import { RestaurantRow } from "../components/RestaurantRow";
import { cities, getTile, search, TILES } from "../lib/data";
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
    const p = new URLSearchParams(params);
    if (next.tile !== undefined) {
      if (next.tile) p.set("tile", next.tile); else p.delete("tile");
    }
    if (next.q !== undefined) {
      if (next.q) p.set("q", next.q); else p.delete("q");
    }
    setParams(p, { replace: next.q !== undefined });
  }

  const where = city ? ` in ${city.name}` : "";
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
      ) : (
        <p className="empty" aria-live="polite">
          Nothing matches{tile ? ` ${tile.label}` : ""}{query ? ` "${query}"` : ""}{where} yet.{" "}
          <button type="button" className="link-button" onClick={() => setParams({})}>Show all restaurants</button>
        </p>
      )}
    </>
  );
}
