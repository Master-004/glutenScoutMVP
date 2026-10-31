// Small shared pieces, ported from the Go templates' design-system components.
import { useEffect, useRef, type ReactNode } from "react";
import { Link } from "react-router-dom";
import { reviewerLabel, shortDate, summarize, type Evidence, type Meal, type Restaurant, type RatingSummary, type Review } from "../lib/data";
import { useStore } from "../lib/store";

const svg = (children: ReactNode, size = 20) => (
  <svg className="icon" width={size} height={size} viewBox="0 0 24 24" aria-hidden="true">{children}</svg>
);
export const Icons = {
  check: svg(<path d="M5 12.5l4.5 4.5L19 7.5" />, 16),
  conflict: svg(<><path d="M12 3l10 18H2z" /><path d="M12 10v5M12 18v.01" /></>, 16),
  unknown: svg(<><circle cx="12" cy="12" r="9" /><path d="M9.5 9.5a2.5 2.5 0 1 1 3.5 2.3c-.7.4-1 .9-1 1.7M12 17v.01" /></>, 16),
  clock: svg(<><circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" /></>, 16),
  x: svg(<><circle cx="12" cy="12" r="9" /><path d="M9 9l6 6M15 9l-6 6" /></>, 16),
  group: svg(<><circle cx="9" cy="8" r="3" /><circle cx="17" cy="9" r="2.5" /><path d="M3 20c0-3.3 2.7-6 6-6s6 2.7 6 6M15 14.5c3 0 5.5 2.2 5.5 5.5" /></>, 16),
  // Plain clipboard, no checkmark: a check would contradict "Reports disagree".
  clipboard: svg(<><rect x="5" y="4" width="14" height="17" rx="2" /><path d="M9 4V3h6v1M9 10h6M9 14h6M9 18h3" /></>, 26),
  dashboard: svg(<><rect x="4" y="4" width="7" height="7" rx="1" /><rect x="13" y="4" width="7" height="7" rx="1" /><rect x="4" y="13" width="7" height="7" rx="1" /><rect x="13" y="13" width="7" height="7" rx="1" /></>, 24),
  search: svg(<><circle cx="11" cy="11" r="7" /><path d="M16.5 16.5L21 21" /></>, 24),
  bookmark: svg(<path d="M6 3h12v18l-6-4-6 4z" />, 24),
  profile: svg(<><circle cx="12" cy="8" r="4" /><path d="M4 21c0-4.4 3.6-8 8-8s8 3.6 8 8" /></>, 24),
};

export function Stars({ value, label }: { value: number; label: string }) {
  // Half-star steps only (design system).
  const rounded = Math.round(value * 2) / 2;
  return (
    <span className="stars" role="img" aria-label={label}>
      {[1, 2, 3, 4, 5].map((n) => (
        <span key={n} className={`star ${rounded >= n ? "star-full" : rounded >= n - 0.5 ? "star-half" : "star-empty"}`} />
      ))}
    </span>
  );
}

// Never a bare or blank rating: an empty spot always says why.
export function Rating({ summary }: { summary: RatingSummary }) {
  if (!summary.ratingCount) {
    return <span className="rating-none">{summary.reviewCount ? "No ratings yet" : "No reviews yet"}</span>;
  }
  const number = summary.average.toFixed(1);
  return (
    <span className="rating">
      <Stars value={summary.average} label={`${number} out of 5 stars`} />
      <span className="rating-num">{number}</span>
      <span className="rating-count">· {summary.ratingCount} {summary.ratingCount === 1 ? "rating" : "ratings"}</span>
    </span>
  );
}

export function GFBadge({ meal }: { meal: Meal }) {
  if (meal.gfStatus === "cross_contact") return <span className="badge badge-conflict">{Icons.conflict}Cross-contact risk</span>;
  if (meal.gfStatus === "contains") return <span className="badge badge-contains">{Icons.x}Contains gluten</span>;
  return <span className="badge badge-confirmed">{Icons.check}GF option</span>;
}

export function EvidenceIcon({ level }: { level: Evidence }) {
  return level === "confirmed" ? Icons.check : level === "conflict" ? Icons.conflict : Icons.clock;
}

export function GroupsLabel({ restaurant }: { restaurant: Restaurant }) {
  if (!restaurant.goodForGroups) return null;
  return <span className="groups-label">{Icons.group}Good for mixed groups</span>;
}

const PLACEHOLDER = "/images/meals/placeholder.svg";

// A missing or broken image falls back to the placeholder, so layouts hold.
export function Img({ src, alt, className }: { src: string; alt: string; className?: string }) {
  return (
    <img
      className={className} src={src || PLACEHOLDER} alt={alt} loading="lazy"
      onError={(e) => { if (!e.currentTarget.src.endsWith(PLACEHOLDER)) e.currentTarget.src = PLACEHOLDER; }}
    />
  );
}

// A bottom sheet on the native <dialog>: the browser traps focus, closes on
// Esc and hands focus back to the button that opened it.
export function Sheet({ open, onClose, title, children, footer }: {
  open: boolean; onClose(): void; title: string; children: ReactNode; footer?: ReactNode;
}) {
  const ref = useRef<HTMLDialogElement>(null);
  useEffect(() => {
    const d = ref.current;
    if (!d) return;
    if (open && !d.open) d.showModal();
    if (!open && d.open) d.close();
  }, [open]);
  return (
    <dialog
      ref={ref} className="sheet" aria-label={title} onClose={onClose}
      // A click whose target is the dialog itself landed on the scrim.
      onClick={(e) => { if (e.target === ref.current) onClose(); }}
    >
      <div className="sheet-inner">
        <div className="sheet-head">
          <h2 className="section-header">{title}</h2>
          <button type="button" className="sheet-close" onClick={onClose} aria-label="Close">✕</button>
        </div>
        <div className="sheet-body">{open && children}</div>
        {footer && <div className="sheet-foot">{footer}</div>}
      </div>
    </dialog>
  );
}

export function MealRow({ meal, foodReviews }: { meal: Meal; foodReviews: Review[] }) {
  return (
    <Link className="meal-row" to={`/meals/${meal.id}`}>
      <Img className="meal-thumb" src={meal.photo} alt={meal.name} />
      <span className="meal-row-body">
        <span className="meal-row-name">{meal.name}</span>
        <Rating summary={summarize(foodReviews)} />
        <GFBadge meal={meal} />
      </span>
      <span className="chevron" aria-hidden="true">›</span>
    </Link>
  );
}

export function ReviewItem({ review }: { review: Review }) {
  const tag = reviewerLabel(review.reviewerType);
  return (
    <li className="review" data-clarity-mask="true">
      <span className="avatar" aria-hidden="true">{review.author.charAt(0).toUpperCase()}</span>
      <div>
        <p className="review-meta"><strong>{review.author}</strong>{tag && <span className="reviewer-tag">{tag}</span>}</p>
        <p className="review-sub">
          {review.rating > 0 && <Stars value={review.rating} label={`Food: ${review.rating} out of 5 stars`} />}
          <time dateTime={review.createdAt}>{shortDate(review.createdAt, true)}</time>
        </p>
        <p className="review-body">{review.body}</p>
      </div>
    </li>
  );
}

// Shared hook: every meal's reviews for a restaurant's rating.
export function useRestaurantSummary(mealIds: number[]): RatingSummary {
  const { reviewsFor } = useStore();
  return summarize(mealIds.flatMap((id) => reviewsFor(id)));
}
