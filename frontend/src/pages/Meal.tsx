import { useCallback, useState, type FormEvent } from "react";
import { Link, useParams } from "react-router-dom";
import { ScoutVideo } from "../components/ScoutVideo";
import { GFBadge, Img, Rating, ReviewItem, Sheet } from "../components/ui";
import { getMeal, getRestaurant, REVIEWER_TYPES, summarize, validateReview, type ReviewErrors } from "../lib/data";
import { useStore } from "../lib/store";
import { NotFound } from "./NotFound";
import { BackLink } from "./Restaurant";

function ReviewForm({ onSubmit }: { onSubmit(author: string, reviewerType: string, rating: number, body: string): void }) {
  const [author, setAuthor] = useState("");
  const [reviewerType, setReviewerType] = useState("");
  const [rating, setRating] = useState(0);
  const [body, setBody] = useState("");
  const [errors, setErrors] = useState<ReviewErrors>({});

  function submit(e: FormEvent) {
    e.preventDefault();
    const found = validateReview({ author, reviewerType, rating, body });
    setErrors(found);
    if (Object.keys(found).length === 0) onSubmit(author, reviewerType, rating, body);
  }

  return (
    <form className="review-form" onSubmit={submit} noValidate data-clarity-mask="true">
      <div className="field">
        <label htmlFor="author">Display name <span className="optional">(optional)</span></label>
        <input id="author" maxLength={40} placeholder="Anonymous" value={author} onChange={(e) => setAuthor(e.target.value)} />
      </div>
      <fieldset className="field">
        <legend>Which describes you?</legend>
        <div className="choices">
          {REVIEWER_TYPES.map((t) => (
            <label key={t.value} className="choice">
              <input type="radio" name="reviewerType" value={t.value} checked={reviewerType === t.value} onChange={() => setReviewerType(t.value)} />
              <span>{t.label}</span>
            </label>
          ))}
        </div>
        {errors.reviewerType && <p className="field-error" role="alert">{errors.reviewerType}</p>}
      </fieldset>
      <fieldset className="field">
        <legend>Food rating</legend>
        <div className="choices">
          {[1, 2, 3, 4, 5].map((n) => (
            <label key={n} className="choice">
              <input type="radio" name="rating" value={n} checked={rating === n} onChange={() => setRating(n)} />
              <span><span aria-hidden="true">★</span> {n}<span className="visually-hidden"> {n === 1 ? "star" : "stars"}</span></span>
            </label>
          ))}
        </div>
        <p className="hint">Rate the food, not gluten safety. Skip this if you work at the restaurant.</p>
        {errors.rating && <p className="field-error" role="alert">{errors.rating}</p>}
      </fieldset>
      <div className="field">
        <label htmlFor="body">Your review</label>
        <textarea id="body" rows={4} placeholder="How did it go? Any tips for ordering?" value={body} onChange={(e) => setBody(e.target.value)} />
        {errors.body && <p className="field-error" role="alert">{errors.body}</p>}
      </div>
      <button type="submit" className="button button-full">Post review</button>
    </form>
  );
}

export function Meal() {
  const { id } = useParams();
  const { reviewsFor, addReview, upvotedMeals, toggleUpvote, videos } = useStore();
  const [writing, setWriting] = useState(false);
  const [celebrating, setCelebrating] = useState(false);
  const [posted, setPosted] = useState(false);
  const finish = useCallback(() => setCelebrating(false), []);
  const meal = getMeal(Number(id));
  const restaurant = meal && getRestaurant(meal.restaurantId);
  if (!meal || !restaurant) return <NotFound />;

  const reviews = reviewsFor(meal.id);
  const upvoted = upvotedMeals.includes(meal.id);

  return (
    <>
      <BackLink fallback={`/restaurants/${restaurant.id}`} />
      <Img className="meal-hero" src={meal.photo} alt={meal.name} />
      <h1 className="identity-title">{meal.name}</h1>
      <Rating summary={summarize(reviews)} />
      {reviews.some((r) => r.rating > 0) && <p className="hint">Stars rate the food, not gluten safety.</p>}
      <p className="meta">at <Link to={`/restaurants/${restaurant.id}`}>{restaurant.name}</Link></p>
      <p className="meal-badge"><GFBadge meal={meal} /></p>
      {meal.description && <p>{meal.description}</p>}
      {meal.gfNote && <p className="gf-note">{meal.gfNote}</p>}
      <div className="meal-actions">
        <button
          type="button" className={`upvote-btn${upvoted ? " is-upvoted" : ""}`} aria-pressed={upvoted}
          onClick={() => toggleUpvote(meal.id)}
        >
          <span aria-hidden="true">▲</span> {upvoted ? 1 : 0}<span className="visually-hidden"> upvotes</span>
        </button>
        {/* Copy rule: an upvote means it went well for one diner, not "safe". */}
        <span className="hint">Upvote if you'd order this again</span>
      </div>

      <section id="reviews" aria-labelledby="reviews-heading">
        <div className="section-head">
          <h2 id="reviews-heading" className="section-header">Reviews ({reviews.length})</h2>
          <button type="button" className="button" onClick={() => setWriting(true)}>Write a review</button>
        </div>
        {posted && <p className="form-success" role="status">✓ Thanks! Your review was posted.</p>}
        {reviews.length ? (
          <ol className="review-list">{reviews.map((r) => <ReviewItem key={r.id} review={r} />)}</ol>
        ) : (
          <p className="empty">No reviews yet. Be the first to share how it went.</p>
        )}
      </section>

      <Sheet open={writing} onClose={() => setWriting(false)} title="Write a review">
        <ReviewForm
          onSubmit={(author, reviewerType, rating, body) => {
            addReview(meal.id, { author, reviewerType, rating, body });
            setWriting(false);
            setPosted(true);
            // The reward for reviewing, so nobody is made to wait for it.
            if (videos) setCelebrating(true);
          }}
        />
      </Sheet>
      {celebrating && <ScoutVideo src="/videos/scout-badge.mp4" caption="Review badge earned!" onDone={finish} />}
    </>
  );
}
