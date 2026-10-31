import { useEffect, useRef } from "react";
import { REVIEWER_TYPES } from "../lib/data";
import { useStore } from "../lib/store";

// Restaurant staff can pick their option in Profile.
const DINER_TYPES = REVIEWER_TYPES.filter((t) => t.value !== "restaurant");

// The first screen, right after the intro video. There's no real login yet: the
// choice stays in this browser. It's health-related, so it's masked in Clarity.
export function Start({ ready }: { ready: boolean }) {
  const { finishStart } = useStore();
  const heading = useRef<HTMLHeadingElement>(null);

  // Once the video is gone, keyboard and screen reader users start here.
  useEffect(() => {
    if (ready) heading.current?.focus();
  }, [ready]);

  return (
    <main className="page start" data-clarity-mask="true">
      <img className="identity-photo" src="/images/scout.jpg" alt="Scout, the GlutenScout guide" />
      <h1 ref={heading} tabIndex={-1} className="identity-title">Which describes you?</h1>
      <p className="start-text">We'll show this on your reviews, so other diners know how to weigh them. You can change it anytime in Profile.</p>
      <ul className="start-choices">
        {DINER_TYPES.map((t) => (
          <li key={t.value}>
            <button type="button" className="start-choice" onClick={() => finishStart(t.value)}>{t.label}</button>
          </li>
        ))}
      </ul>
      <button type="button" className="link-button start-skip" onClick={() => finishStart(null)}>Skip for now</button>
    </main>
  );
}
