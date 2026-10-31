import { Link } from "react-router-dom";

export function NotFound() {
  return (
    <>
      <h1 className="page-title">Page not found</h1>
      <p className="empty">We couldn't find that page. <Link to="/">Back to your dashboard</Link></p>
    </>
  );
}
