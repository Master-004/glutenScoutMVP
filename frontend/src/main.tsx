import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { Layout } from "./components/Layout";
import { startClarity } from "./lib/clarity";
import { StoreProvider } from "./lib/store";
import { Dashboard } from "./pages/Dashboard";
import { Meal } from "./pages/Meal";
import { NotFound } from "./pages/NotFound";
import { Profile, Saved } from "./pages/Other";
import { Restaurant } from "./pages/Restaurant";
import { Search } from "./pages/Search";
import "./styles.css";

startClarity();

const router = createBrowserRouter([
  {
    element: <Layout />,
    children: [
      { path: "/", element: <Dashboard /> },
      { path: "/search", element: <Search /> },
      { path: "/saved", element: <Saved /> },
      { path: "/profile", element: <Profile /> },
      { path: "/restaurants/:id", element: <Restaurant /> },
      { path: "/meals/:id", element: <Meal /> },
      { path: "*", element: <NotFound /> },
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
