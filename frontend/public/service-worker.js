const CACHE_NAME = "ghopper-v1";
const urlsToCache = [
  "/",
  "/index.html",
  "/manifest.json",
  "/app-icon.svg",
  // Add other static assets here
];

// This line is required for vite-plugin-pwa
// eslint-disable-next-line no-unused-vars
const manifestFiles = self.__WB_MANIFEST;

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(urlsToCache))
  );
});

self.addEventListener("fetch", (event) => {});
