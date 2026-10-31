import { useCallback } from "react";
import { Link, NavLink, Outlet, ScrollRestoration, useLocation } from "react-router-dom";
import { clarityEnabled } from "../lib/clarity";
import { useStore } from "../lib/store";
import { ScoutVideo } from "./ScoutVideo";
import { Start } from "./Start";
import { Icons } from "./ui";

const TABS = [
  { to: "/", label: "Dashboard", icon: Icons.dashboard, end: true },
  { to: "/search", label: "Search", icon: Icons.search, end: false },
  { to: "/saved", label: "Saved", icon: Icons.bookmark, end: false },
  { to: "/profile", label: "Profile", icon: Icons.profile, end: false },
];

export function Layout() {
  const { videos, introSeen, started, markIntroSeen } = useStore();
  const finishIntro = useCallback(() => markIntroSeen(), [markIntroSeen]);
  const { pathname } = useLocation();
  // The privacy page has to be readable before someone gets past the start screen.
  const onPrivacy = pathname === "/privacy";
  const playingIntro = videos && !introSeen && !onPrivacy;
  return (
    <div className={`app${started ? "" : " app-start"}`}>
      {started || onPrivacy ? <>
      {started && <nav className="tab-bar" aria-label="Main">
        {TABS.map((t) => (
          <NavLink key={t.to} to={t.to} end={t.end} className="tab">
            {t.icon}<span>{t.label}</span>
          </NavLink>
        ))}
      </nav>}
      <main className="page">
        <Outlet />
      </main>
      </> : <Start ready={!playingIntro} />}
      <footer className="disclaimer">
        GlutenScout info comes from the community and may be out of date. Always confirm with the restaurant before ordering.
        {clarityEnabled && <> We use Microsoft Clarity to see how this prototype is used. <Link to="/privacy">Privacy</Link></>}
      </footer>
      <ScrollRestoration />
      {playingIntro && (
        <ScoutVideo src="/videos/scout-intro.mp4" caption="Let's find some gluten-free meals together" onDone={finishIntro} />
      )}
    </div>
  );
}
