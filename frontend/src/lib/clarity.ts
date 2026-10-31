// Microsoft Clarity session recordings. Off unless VITE_CLARITY_ID is set, so
// nobody is recorded by accident: testers must be told before it's turned on.
// Review text and forms are masked in the pages with data-clarity-mask.
type ClarityFn = ((...args: unknown[]) => void) & { q?: unknown[][] };
declare global {
  interface Window { clarity?: ClarityFn }
}

export function startClarity(id = import.meta.env.VITE_CLARITY_ID as string | undefined) {
  if (!id || window.clarity) return;
  const queue: ClarityFn = (...args) => { (queue.q = queue.q ?? []).push(args); };
  window.clarity = queue;
  const script = document.createElement("script");
  script.async = true;
  script.src = `https://www.clarity.ms/tag/${encodeURIComponent(id)}`;
  document.head.appendChild(script);
}

// Labels a session so recordings can be compared, e.g. videos on vs off.
export function tagSession(key: string, value: string) {
  window.clarity?.("set", key, value);
}
