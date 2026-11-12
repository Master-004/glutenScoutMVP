import * as Sentry from "@sentry/react";
import { useEffect } from "react";
import { createBrowserRouter, createRoutesFromChildren, matchRoutes, useLocation, useNavigationType } from "react-router-dom";

// Error monitoring with Sentry. It stays off unless VITE_SENTRY_DSN is set,
// the same way Clarity waits for VITE_CLARITY_ID.
const dsn = import.meta.env.VITE_SENTRY_DSN as string | undefined;
export const monitoringEnabled = Boolean(dsn);

if (dsn) {
  Sentry.init({
    dsn,
    environment: import.meta.env.VITE_ENVIRONMENT || import.meta.env.MODE,
    release: import.meta.env.VITE_RELEASE || undefined,
    // No user details or cookies. No search text either: the search page keeps
    // what someone typed in the URL (?q=).
    dataCollection: { userInfo: false, cookies: false, httpHeaders: false, urlQueryParams: false },
    integrations: [
      Sentry.reactRouterBrowserTracingIntegration({ useEffect, useLocation, useNavigationType, createRoutesFromChildren, matchRoutes }),
    ],
    // Every trace locally so a test shows up; a tenth in production.
    tracesSampleRate: import.meta.env.DEV ? 1.0 : 0.1,
    // No Session Replay: Clarity covers that, and only with consent.
    beforeBreadcrumb(crumb) {
      // Button and link text can include what someone typed or chose about their diet.
      if (crumb.category === "ui.click" || crumb.category === "ui.input") delete crumb.message;
      return crumb;
    },
  });
}

// Wrapping the router names traces by route (/restaurants/:id), not by URL.
export const createRouter: typeof createBrowserRouter = dsn ? Sentry.wrapCreateBrowserRouter(createBrowserRouter) : createBrowserRouter;

export function reportError(error: unknown) {
  if (dsn) Sentry.captureException(error);
}
