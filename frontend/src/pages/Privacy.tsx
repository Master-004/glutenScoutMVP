import { Link } from "react-router-dom";
import { clarityEnabled } from "../lib/clarity";
import { useStore } from "../lib/store";
import { BackLink } from "./Restaurant";

// DRAFT wording, based on Microsoft's sample disclosure. Review it before
// turning Clarity on.
export function Privacy() {
  const { started } = useStore();
  return (
    <>
      <BackLink fallback="/" />
      <h1 className="page-title">Privacy</h1>
      <p>GlutenScout is a prototype for people 18 and older. It uses sample data about fictional restaurants.</p>
      <h2 className="section-header">What stays on your device</h2>
      <p>Your city, saved places, reviews and the option you chose about your diet stay in this browser. They aren't sent to GlutenScout.</p>
      <h2 className="section-header">Session recordings</h2>
      {clarityEnabled ? (
        <>
          <p>We use Microsoft Clarity to see how people use this prototype, so we can improve it. Clarity captures session replays, heatmaps, and how you tap and scroll, using cookies and similar technology. It's not used for advertising.</p>
          <p>We hide what you type, reviews, and your dietary choice from these recordings. Using this prototype means you agree that we and Microsoft can collect and use this data. You can turn recording off anytime{started ? <> in <Link to="/profile">Profile</Link></> : " in Profile"}.</p>
          <p>Learn how Microsoft uses data in the <a href="https://www.microsoft.com/privacy/privacystatement" target="_blank" rel="noreferrer">Microsoft Privacy Statement</a>.</p>
        </>
      ) : (
        <p>Session recording is off right now.</p>
      )}
    </>
  );
}
