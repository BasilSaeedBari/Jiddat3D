const CACHE_NAME = 'jiddat3d-cache-v8';

const STATIC_ASSETS = [
	'/static/css/output.css?v=9',
	'/static/js/htmx.min.js',
	'/static/js/alpine.min.js',
	'/static/js/scroll-reveal.js',
	'/static/img/logo/JiddatBrightTrasnparent-removebg-preview.png'
];

self.addEventListener('install', (event) => {
	event.waitUntil(
		caches.open(CACHE_NAME).then((cache) => {
			return cache.addAll(STATIC_ASSETS);
		}).then(() => self.skipWaiting())
	);
});

self.addEventListener('activate', (event) => {
	event.waitUntil(
		caches.keys().then((cacheNames) => {
			return Promise.all(
				cacheNames.map((cacheName) => {
					if (cacheName !== CACHE_NAME) {
						return caches.delete(cacheName);
					}
				})
			);
		}).then(() => self.clients.claim())
	);
});

self.addEventListener('fetch', (event) => {
	if (event.request.method !== 'GET') return;
	
	const url = new URL(event.request.url);

	// Ignore cross-origin requests (e.g. Google Fonts, CDNs)
	if (url.origin !== self.location.origin) {
		return;
	}
	
	// API and Admin UI should never be cached
	if (url.pathname.startsWith('/api/') || url.pathname.startsWith('/_/')) {
		return;
	}

	// Static assets: Cache-first (only cache successful responses)
	if (url.pathname.startsWith('/static/')) {
		event.respondWith(
			caches.match(event.request).then((cachedResponse) => {
				if (cachedResponse) {
					return cachedResponse;
				}
				return fetch(event.request).then((networkResponse) => {
					if (networkResponse && networkResponse.ok) {
						const copy = networkResponse.clone();
						caches.open(CACHE_NAME).then((cache) => cache.put(event.request, copy));
					}
					return networkResponse;
				});
			})
		);
		return;
	}

	// HTML Pages: Network-first, fallback to cache (only cache successful responses)
	event.respondWith(
		fetch(event.request)
			.then((networkResponse) => {
				if (networkResponse && networkResponse.ok) {
					const copy = networkResponse.clone();
					caches.open(CACHE_NAME).then((cache) => cache.put(event.request, copy));
				}
				return networkResponse;
			})
			.catch(() => {
				return caches.match(event.request);
			})
	);
});
