const CACHE_NAME = 'jiddat3d-cache-v3';

const STATIC_ASSETS = [
	'/static/css/output.css?v=3',
	'/static/js/htmx.min.js',
	'/static/js/alpine.min.js',
	'/static/js/scroll-reveal.js',
	'/static/img/logo/JiddatBrightTrasnparent.png'
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
	
	// API and Admin UI should never be cached
	if (url.pathname.startsWith('/api/') || url.pathname.startsWith('/_/')) {
		return;
	}

	// Static assets: Cache-first
	if (url.pathname.startsWith('/static/')) {
		event.respondWith(
			caches.match(event.request).then((cachedResponse) => {
				if (cachedResponse) {
					return cachedResponse;
				}
				return fetch(event.request).then((networkResponse) => {
					return caches.open(CACHE_NAME).then((cache) => {
						cache.put(event.request, networkResponse.clone());
						return networkResponse;
					});
				});
			})
		);
		return;
	}

	// HTML Pages: Network-first, fallback to cache
	event.respondWith(
		fetch(event.request)
			.then((networkResponse) => {
				return caches.open(CACHE_NAME).then((cache) => {
					cache.put(event.request, networkResponse.clone());
					return networkResponse;
				});
			})
			.catch(() => {
				return caches.match(event.request);
			})
	);
});
