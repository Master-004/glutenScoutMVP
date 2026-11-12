import { sentryVitePlugin } from "@sentry/bundler-plugins/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// Vercel sets these at build time. They tell Sentry which deploy and which
// environment (production or preview) an error came from.
const release = process.env.VERCEL_GIT_COMMIT_SHA ?? "";
const environment = process.env.VERCEL_ENV ?? "";
// Only builds that have a Sentry token (Vercel) upload source maps.
const authToken = process.env.SENTRY_AUTH_TOKEN;

export default defineConfig({
  // "hidden" makes the maps without pointing to them from the code. They're
  // uploaded to Sentry, then deleted so they're never published.
  build: { sourcemap: authToken ? "hidden" : false },
  plugins: [
    react(),
    // Must stay last.
    authToken && sentryVitePlugin({
      org: "brigham-young-university-r5",
      project: "glutenscout-frontend",
      authToken,
      release: release ? { name: release } : undefined,
      sourcemaps: { filesToDeleteAfterUpload: ["./dist/**/*.map"] },
      telemetry: false,
    }),
  ],
  define: {
    "import.meta.env.VITE_RELEASE": JSON.stringify(release),
    "import.meta.env.VITE_ENVIRONMENT": JSON.stringify(environment),
  },
});
