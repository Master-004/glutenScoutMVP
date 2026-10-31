// What the prototype remembers for one visitor: chosen city, saved places, new
// reviews, upvotes. It lives in the browser because there are no accounts yet,
// and starts empty, which is what a real first-time user sees.
import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from "react";
import { sampleReviews, type Review, type ReviewInput } from "./data";
import { clarityEnabled, startClarity, stopClarity, tagSession } from "./clarity";

export type Theme = "auto" | "light" | "dark";

interface Saved {
  citySlug: string | null;
  savedRestaurants: number[];
  upvotedMeals: number[];
  myReviews: Review[];
  theme: Theme;
  videos: boolean;
  introSeen: boolean;
  started: boolean; // has seen the first-time start screen
  reviewerType: string | null; // Celiac, Strict GF...: filled in on their reviews
  recording: boolean; // allows Clarity, once it's turned on for the app
}
const EMPTY: Saved = {
  citySlug: null, savedRestaurants: [], upvotedMeals: [], myReviews: [], theme: "auto", videos: true, introSeen: false,
  started: false, reviewerType: null, recording: true,
};
const KEY = "glutenscout-prototype";

// Storage can throw in private windows, so the app must still work without it.
function load(): Saved {
  try {
    const raw = localStorage.getItem(KEY);
    return raw ? { ...EMPTY, ...JSON.parse(raw) } : EMPTY;
  } catch {
    return EMPTY;
  }
}
function persist(state: Saved) {
  try {
    localStorage.setItem(KEY, JSON.stringify(state));
  } catch {
    // Nothing to do: the choice just won't survive a refresh.
  }
}

interface Store extends Saved {
  setCity(slug: string): void;
  toggleSaved(id: number): void;
  toggleUpvote(id: number): void;
  addReview(mealId: number, input: ReviewInput): void;
  setTheme(theme: Theme): void;
  setVideos(on: boolean): void;
  markIntroSeen(): void;
  finishStart(reviewerType: string | null): void;
  setReviewerType(reviewerType: string | null): void;
  restart(): void;
  setRecording(on: boolean): void;
  reset(): void;
  reviewsFor(mealId: number): Review[];
}
const Ctx = createContext<Store | null>(null);

const toggled = (list: number[], id: number) => (list.includes(id) ? list.filter((x) => x !== id) : [...list, id]);

export function StoreProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<Saved>(() => {
    const loaded = load();
    // ?videos=off in a test link sets this visitor's version of the app.
    const param = new URLSearchParams(window.location.search).get("videos");
    return param === "off" || param === "on" ? { ...loaded, videos: param === "on" } : loaded;
  });

  useEffect(() => persist(state), [state]);
  useEffect(() => {
    if (state.theme === "auto") document.documentElement.removeAttribute("data-theme");
    else document.documentElement.setAttribute("data-theme", state.theme);
  }, [state.theme]);
  // Recording starts only after the start screen, where people see the notice
  // and pick a dietary option, so that choice is never part of a recording.
  useEffect(() => {
    if (!state.started || !state.recording) return;
    startClarity();
    tagSession("videos", state.videos ? "on" : "off");
  }, [state.started, state.recording, state.videos]);

  const update = useCallback((change: (s: Saved) => Partial<Saved>) => {
    setState((s) => ({ ...s, ...change(s) }));
  }, []);

  const store = useMemo<Store>(() => ({
    ...state,
    setCity: (citySlug) => update(() => ({ citySlug })),
    toggleSaved: (id) => update((s) => ({ savedRestaurants: toggled(s.savedRestaurants, id) })),
    toggleUpvote: (id) => update((s) => ({ upvotedMeals: toggled(s.upvotedMeals, id) })),
    addReview: (mealId, input) => update((s) => ({
      myReviews: [...s.myReviews, {
        id: Date.now(), mealId, author: input.author.trim() || "Anonymous", body: input.body.trim(),
        reviewerType: input.reviewerType,
        rating: input.reviewerType === "restaurant" ? 0 : input.rating,
        createdAt: new Date().toISOString(),
      }],
    })),
    setTheme: (theme) => update(() => ({ theme })),
    setVideos: (videos) => update(() => ({ videos })),
    markIntroSeen: () => update(() => ({ introSeen: true })),
    finishStart: (reviewerType) => update(() => ({ started: true, reviewerType })),
    setReviewerType: (reviewerType) => update(() => ({ reviewerType })),
    restart: () => update(() => ({ introSeen: false, started: false })),
    setRecording: (recording) => {
      update(() => ({ recording }));
      if (!recording && clarityEnabled) {
        stopClarity();
        // A fresh page guarantees the script is gone, once the choice is saved.
        window.setTimeout(() => window.location.reload(), 50);
      }
    },
    reset: () => update((s) => ({ ...s, savedRestaurants: [], upvotedMeals: [], myReviews: [] })),
    // Newest first, with this visitor's own reviews included.
    reviewsFor: (mealId) => [...sampleReviews, ...state.myReviews]
      .filter((r) => r.mealId === mealId)
      .sort((a, b) => b.createdAt.localeCompare(a.createdAt)),
  }), [state, update]);

  return <Ctx.Provider value={store}>{children}</Ctx.Provider>;
}

export function useStore(): Store {
  const store = useContext(Ctx);
  if (!store) throw new Error("useStore must be used inside StoreProvider");
  return store;
}
