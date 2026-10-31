import { useCallback } from "react";
import { NavLink, Outlet, ScrollRestoration } from "react-router-dom";
import { useStore } from "../lib/store";
import { ScoutVideo } from "./ScoutVideo";
import { Icons } from "./ui";

const TABS = [
  { to: "/", label: "Dashboard", icon: Icons.dashboard, end: true },
  { to: "/search", label: "Search", icon: Icons.search, end: false },
  { to: "/saved", label: "Saved", icon: Icons.bookmark, end: false },
  { to: "/profile", label: "Profile", icon: Icons.profile, end: false },
];

export function Layout() {
  const { videos, introSeen, markIntroSeen } = useStore();
  const finishIntro = useCallback(() => markIntroSeen(), [markIntroSeen]);
  return (
    <div className="app">
      <nav className="tab-bar" aria-label="Main">
        {TABS.map((t) => (
          <NavLink key={t.to} to={t.to} end={t.end} className="tab">
            {t.icon}<span>{t.label}</span>
          </NavLink>
        ))}
      </nav>
      <main className="page">
        <Outlet />
      </main>
      <footer className="disclaimer">
        GlutenScout info comes from the community and may be out of date. Always confirm with the restaurant before ordering.
      </footer>
      <ScrollRestoration />
      {videos && !introSeen && (
        <ScoutVideo src="/videos/scout-intro.mp4" caption="Scout's on the way" onDone={finishIntro} />
      )}
    </div>
  );
}
