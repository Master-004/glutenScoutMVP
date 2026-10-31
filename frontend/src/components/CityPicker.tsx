import { forwardRef } from "react";
import { cities } from "../lib/data";
import { useStore } from "../lib/store";

// "prominent" makes choosing a city the green next step until it's done, then
// it shrinks out of the way so the page's attention moves to what comes next.
export const CityPicker = forwardRef<HTMLSelectElement, { error?: string; prominent?: boolean }>(
  function CityPicker({ error, prominent }, ref) {
    const { citySlug, setCity } = useStore();
    const style = citySlug ? "city-compact" : prominent ? "city-primary" : "city-plain";
    return (
      <div className={`city-picker ${style}`}>
        <label htmlFor="city" className="visually-hidden">Your city</label>
        <div className="select-wrap">
          <select
            id="city" ref={ref} value={citySlug ?? ""} onChange={(e) => setCity(e.target.value)}
            aria-invalid={error ? true : undefined} aria-describedby={error ? "city-error" : undefined}
          >
            <option value="" disabled>Choose your city</option>
            {cities.map((c) => <option key={c.slug} value={c.slug}>{c.name}, {c.state}</option>)}
          </select>
        </div>
        {error && <p className="field-error" id="city-error" role="alert">{error}</p>}
      </div>
    );
  },
);
