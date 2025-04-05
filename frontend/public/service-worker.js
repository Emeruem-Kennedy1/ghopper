const CACHE_NAME = "ghopper-v1";
const urlsToCache = [
  "/",
  "/index.html",
  "/manifest.json",
  "/app-icon.svg",
  // Add other static assets here
];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(urlsToCache))
  );
});

self.addEventListener("fetch", (event) => {
  // Network-first strategy
  event.respondWith(
    fetch(event.request).catch(() => caches.match(event.request))
  );
});
