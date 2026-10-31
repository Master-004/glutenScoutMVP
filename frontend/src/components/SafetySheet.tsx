import { shortDate, precautionsFor, safetyFor, safetyStatus, type Restaurant } from "../lib/data";
import { EvidenceIcon, Icons, Sheet } from "./ui";

const LEVEL_WORDS = { confirmed: "Confirmed", conflict: "Reports disagree", unverified: "Unverified" };

// The button shows the status so warnings stay on the page, even for someone
// who never opens the sheet. The green fill makes it the page's next action.
export function SafetyButton({ restaurantId, onClick }: { restaurantId: number; onClick(): void }) {
  const status = safetyStatus(restaurantId);
  return (
    <button type="button" className="safety-button" onClick={onClick}>
      <span className="safety-button-icon">{Icons.clipboard}</span>
      <span className="safety-button-text">
        <span className="safety-button-title">Check kitchen safety</span>
        <span className="safety-button-status">
          {status.level === "conflict" ? Icons.conflict : status.level === "confirmed" ? Icons.check : Icons.clock}
          {status.text}
        </span>
      </span>
      <span className="chevron" aria-hidden="true">›</span>
    </button>
  );
}

export function SafetySheet({ restaurant, open, onClose }: { restaurant: Restaurant; open: boolean; onClose(): void }) {
  const profile = safetyFor(restaurant.id);
  const list = precautionsFor(restaurant.id);
  const toConfirm = list.filter((p) => p.status !== "confirmed");
  const digits = restaurant.phone.replace(/\D/g, "");
  return (
    <Sheet
      open={open} onClose={onClose} title="Kitchen safety"
      footer={<button type="button" className="button button-full" onClick={onClose}>OK</button>}
    >
      <p>Based on how the kitchen handles prep, from the restaurant's answers and gluten-free diners' reports. It's evidence, not a guarantee: you decide.</p>

      {profile ? (
        <>
          <p className="verified">{Icons.clock} Last verified {shortDate(profile.verifiedOn, true)}</p>
          {profile.conflictNote && <p className="conflict-note">{Icons.conflict}<span>{profile.conflictNote}</span></p>}
          {list.length > 0 && (
            <ul className="badges">
              {list.map((p) => (
                <li key={p.label} className={`badge badge-${p.status}`}>
                  <EvidenceIcon level={p.status} />
                  <span className="visually-hidden">{LEVEL_WORDS[p.status]}: </span>{p.label}
                </li>
              ))}
            </ul>
          )}
          <p className="meta"><strong>Restaurant answers</strong> · {profile.answersUpdatedOn ? `updated ${shortDate(profile.answersUpdatedOn, true)}` : "not answered yet"}</p>
        </>
      ) : (
        <p className="verified">{Icons.clock} Nobody has reported on this kitchen yet. Call ahead and ask before you order.</p>
      )}

      <details className="sheet-details">
        <summary>What do these labels mean?</summary>
        <dl className="levels">
          <dt><span className="badge badge-confirmed">{Icons.check}Confirmed</span></dt>
          <dd>The restaurant and recent gluten-free diners agree.</dd>
          <dt><span className="badge badge-conflict">{Icons.conflict}Reports disagree</span></dt>
          <dd>Reports conflict, or nobody has backed up the claim yet. Worth asking about.</dd>
          <dt><span className="badge badge-unverified">{Icons.clock}Unverified</span></dt>
          <dd>Only the restaurant has said so, or the information is out of date.</dd>
        </dl>
        <p><strong>Stars rate the food,</strong> not gluten safety.</p>
      </details>

      <details className="sheet-details">
        <summary>Call ahead: questions to ask</summary>
        <ol className="script">
          <li><strong>Meals:</strong> Which dishes can you make gluten-free today?</li>
          <li><strong>Staff practices:</strong> Do cooks change gloves and use clean utensils for gluten-free orders?</li>
          <li><strong>Food prep:</strong> Is there a separate prep area? Are fryers, ovens or grills shared?</li>
          {toConfirm.length > 0 && <li><strong>For this kitchen:</strong> Ask about {toConfirm.map((p) => p.label.toLowerCase()).join(", ")}.</li>}
        </ol>
        {digits
          ? <a className="button-secondary" href={`tel:+1${digits}`}>Call {restaurant.phone}</a>
          : <p className="meta">We don't have a phone number for this restaurant yet.</p>}
      </details>
    </Sheet>
  );
}
