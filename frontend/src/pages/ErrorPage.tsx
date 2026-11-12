import { useEffect } from "react";
import { useRouteError } from "react-router-dom";
import { monitoringEnabled, reportError } from "../lib/monitoring";

// Shown when a page crashes, instead of a blank screen. The error goes to
// Sentry so we hear about it without the tester having to tell us.
// standalone: the whole layout crashed, so there's no page wrapper around this.
export function ErrorPage({ standalone = false }: { standalone?: boolean }) {
  const error = useRouteError();
  useEffect(() => reportError(error), [error]);
  return (
    <div className={standalone ? "page" : undefined}>
      <h1 className="page-title">Something went wrong</h1>
      <p className="empty">
        This page hit a problem.{monitoringEnabled ? " We've been told about it." : ""}{" "}
        <a href="/">Go back to your dashboard</a>
      </p>
    </div>
  );
}
