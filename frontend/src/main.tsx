// Monitoring first, so it can catch errors from everything loaded after it.
import { createRouter } from "./lib/monitoring";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import { Layout } from "./components/Layout";
import { StoreProvider } from "./lib/store";
import { Dashboard } from "./pages/Dashboard";
import { ErrorPage } from "./pages/ErrorPage";
import { Meal } from "./pages/Meal";
import { NotFound } from "./pages/NotFound";
import { Profile, Saved } from "./pages/Other";
import { Privacy } from "./pages/Privacy";
import { Restaurant } from "./pages/Restaurant";
import { Search } from "./pages/Search";
import "./styles.css";

const router = createRouter([
  {
    element: <Layout />,
    errorElement: <ErrorPage standalone />,
    children: [
      {
        // A crash inside a page keeps the tab bar, so people can move on.
        errorElement: <ErrorPage />,
        children: [
          { path: "/", element: <Dashboard /> },
          { path: "/search", element: <Search /> },
          { path: "/saved", element: <Saved /> },
          { path: "/profile", element: <Profile /> },
          { path: "/privacy", element: <Privacy /> },
          { path: "/restaurants/:id", element: <Restaurant /> },
          { path: "/meals/:id", element: <Meal /> },
          { path: "*", element: <NotFound /> },
        ],
      },
    ],
  },
]);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <StoreProvider>
      <RouterProvider router={router} />
    </StoreProvider>
  </StrictMode>,
);
