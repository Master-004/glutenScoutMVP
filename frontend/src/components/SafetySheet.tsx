import { precautionsFor, safetyFor, safetyStatus, shortDate, type Restaurant, type SafetyLevel } from "../lib/data";
import { EvidenceIcon, Icons, Sheet } from "./ui";

// Status is never color alone: each level has its own icon and words too.
export function StatusIcon({ level }: { level: SafetyLevel }) {
  return level === "none" ? Icons.unknown : <EvidenceIcon level={level} />;
}

// The status sits beside the button, so warnings stay on the page even for
// someone who never opens the sheet. The button is green at every level,
// because checking is always the next step.
export function SafetyButton({ restaurantId, onClick }: { restaurantId: number; onClick(): void }) {
  const status = safetyStatus(restaurantId);
  return (
    <div className="safety-row">
      <div className={`safety-status status-${status.level}`}>
        <span className="safety-status-label"><StatusIcon level={status.level} />{status.label}</span>
        {status.verifiedOn && <span className="safety-status-date">Verified {shortDate(status.verifiedOn)}</span>}
      </div>
      <button type="button" className="safety-button" onClick={onClick}>
        {Icons.clipboard}
        <span className="safety-button-title">Check kitchen safety</span>
      </button>
    </div>
  );
}

const LEVEL_WORDS = { confirmed: "Confirmed", conflict: "Reports disagree", unverified: "Unverified" };

// Answers in order: what's the situation, what should I ask, what's confirmed?
// The explanations are folded so they're there for anyone who wants them.
export function SafetySheet({ restaurant, open, onClose }: { restaurant: Restaurant; open: boolean; onClose(): void }) {
  const profile = safetyFor(restaurant.id);
  const list = precautionsFor(restaurant.id);
  const status = safetyStatus(restaurant.id);
  const askAbout = list.filter((p) => p.status !== "confirmed");
  const confirmed = list.filter((p) => p.status === "confirmed");
  const digits = restaurant.phone.replace(/\D/g, "");
  const callInFooter = status.level !== "confirmed" && digits !== "";

  return (
    <Sheet
      open={open} onClose={onClose} title="Kitchen safety"
      footer={
        <div className="sheet-actions">
          {callInFooter && <a className="button-secondary" href={`tel:+1${digits}`}>Call restaurant</a>}
          <button type="button" className="button" onClick={onClose}>OK</button>
        </div>
      }
    >
      <div className={`status-banner status-${status.level}`}>
        <p className="status-banner-title"><StatusIcon level={status.level} />{status.label}</p>
        <p className="status-banner-text">{status.summary}</p>
      </div>

      {askAbout.length > 0 && (
        <section className="practice-group">
          <h3 className="practice-heading">Ask about</h3>
          <ul>
            {askAbout.map((p) => (
              <li key={p.label} className={`practice practice-${p.status}`}>
                <EvidenceIcon level={p.status} />
                <span className="visually-hidden">{LEVEL_WORDS[p.status]}: </span>{p.label}
              </li>
            ))}
          </ul>
        </section>
      )}

      {confirmed.length > 0 && (
        <section className="practice-group">
          <h3 className="practice-heading">Confirmed</h3>
          <ul>
            {confirmed.map((p) => (
              <li key={p.label} className="practice practice-confirmed">
                <EvidenceIcon level="confirmed" />{p.label}
              </li>
            ))}
          </ul>
        </section>
      )}

      {profile && (
        <p className="sheet-fineprint">
          Verified {shortDate(profile.verifiedOn, true)} · Evidence, not a guarantee<br />
          {profile.answersUpdatedOn ? `Restaurant answered ${shortDate(profile.answersUpdatedOn, true)}` : "Restaurant hasn't answered yet"}
        </p>
      )}

      <details className="sheet-details">
        <summary>How this works</summary>
        <p>Based on how the kitchen handles prep, from the restaurant's answers and gluten-free diners' reports.</p>
        <dl className="levels">
          <dt>{Icons.check} Confirmed</dt>
          <dd>The restaurant and recent gluten-free diners agree.</dd>
          <dt>{Icons.conflict} Reports disagree</dt>
          <dd>Reports conflict, or nobody has backed up the claim yet.</dd>
          <dt>{Icons.clock} Unverified</dt>
          <dd>Only the restaurant has said so, or the information is out of date.</dd>
        </dl>
        <p>Stars rate the food, not gluten safety.</p>
      </details>

      <details className="sheet-details">
        <summary>Questions to ask when you call</summary>
        <ol className="script">
          <li>Which dishes can you make gluten-free today?</li>
          <li>Do cooks change gloves and use clean utensils for gluten-free orders?</li>
          <li>Is there a separate prep area? Are fryers, ovens or grills shared?</li>
          {askAbout.length > 0 && <li>Can you tell me about: {askAbout.map((p) => p.label.toLowerCase()).join(", ")}</li>}
        </ol>
        {!callInFooter && digits && <a className="button-secondary" href={`tel:+1${digits}`}>Call {restaurant.phone}</a>}
        {!digits && <p className="meta">We don't have a phone number for this restaurant yet.</p>}
      </details>
    </Sheet>
  );
}
