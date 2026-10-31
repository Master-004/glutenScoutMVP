import { Link } from "react-router-dom";
import { mealsFor, safetyStatus, shortDate, type SearchResult } from "../lib/data";
import { StatusIcon } from "./SafetySheet";
import { GroupsLabel, Rating, useRestaurantSummary } from "./ui";

export function RestaurantRow({ result, foodType }: { result: SearchResult; foodType?: string }) {
  const { restaurant: r, gfCount, countNoun, matchedMeals } = result;
  const summary = useRestaurantSummary(mealsFor(r.id).map((m) => m.id));
  const status = safetyStatus(r.id);
  return (
    <Link className="restaurant-row" to={`/restaurants/${r.id}${foodType ? `?food=${foodType}` : ""}`}>
      <h3>{r.name}</h3>
      <Rating summary={summary} />
      <p className="meta">{r.cuisine} · {r.address}</p>
      {/* The same status as the restaurant page: warnings are seen early. */}
      <p className="row-status">
        <span className={`status-pill status-${status.level}`}>
          <StatusIcon level={status.level} />{status.label}
          {status.verifiedOn && <span className="status-pill-date">· Verified {shortDate(status.verifiedOn)}</span>}
        </span>
      </p>
      {matchedMeals.length > 0 && <p className="row-match">Has {matchedMeals.join(", ")}</p>}
      <p className="row-labels">
        <span className="chip">{gfCount} GF {countNoun}{gfCount === 1 ? "" : "s"}</span>
        <GroupsLabel restaurant={r} />
      </p>
    </Link>
  );
}
