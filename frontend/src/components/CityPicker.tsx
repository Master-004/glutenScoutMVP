import { forwardRef } from "react";
import { cities } from "../lib/data";
import { useStore } from "../lib/store";

export const CityPicker = forwardRef<HTMLSelectElement, { error?: string }>(function CityPicker({ error }, ref) {
  const { citySlug, setCity } = useStore();
  return (
    <div className="city-picker">
      <label htmlFor="city">Your city</label>
      <select
        id="city" ref={ref} value={citySlug ?? ""} onChange={(e) => setCity(e.target.value)}
        aria-invalid={error ? true : undefined} aria-describedby={error ? "city-error" : undefined}
      >
        <option value="" disabled>Select a city…</option>
        {cities.map((c) => <option key={c.slug} value={c.slug}>{c.name}, {c.state}</option>)}
      </select>
      {error && <p className="field-error" id="city-error" role="alert">{error}</p>}
    </div>
  );
});
