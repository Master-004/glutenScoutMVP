// Microsoft Clarity session recordings. Off unless VITE_CLARITY_ID is set, so
// nobody is recorded by accident. See "Turning on Clarity" in the README.
//
// Never send health details here: no reviewer type in tags or events. Review
// text and forms are hidden in the pages with data-clarity-mask.
type ClarityFn = ((...args: unknown[]) => void) & { q?: unknown[] };
declare global {
  interface Window { clarity?: ClarityFn }
}

const CLARITY_ID = import.meta.env.VITE_CLARITY_ID as string | undefined;
export const clarityEnabled = Boolean(CLARITY_ID);

export function startClarity() {
  if (!CLARITY_ID || window.clarity) return;
  // Microsoft's own snippet: calls queue up until the script arrives. It expects
  // each entry to be an "arguments" object, so this can't use ...args.
  const queue: ClarityFn = function () {
    // eslint-disable-next-line prefer-rest-params
    (queue.q = queue.q ?? []).push(arguments);
  };
  window.clarity = queue;
  const script = document.createElement("script");
  script.async = true;
  script.src = `https://www.clarity.ms/tag/${encodeURIComponent(CLARITY_ID)}`;
  document.head.appendChild(script);
  // Analytics only: this app doesn't advertise.
  window.clarity("consentv2", { ad_Storage: "denied", analytics_Storage: "granted" });
}

// Clears Clarity's cookies and stops tracking, for someone who opts out.
export function stopClarity() {
  window.clarity?.("consent", false);
}

// Labels a session so recordings can be compared, e.g. videos on vs off.
export function tagSession(key: string, value: string) {
  window.clarity?.("set", key, value);
}

// Marks a moment in a recording, e.g. opening the safety sheet.
export function trackEvent(name: string) {
  window.clarity?.("event", name);
}
