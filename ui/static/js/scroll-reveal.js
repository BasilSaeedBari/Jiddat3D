document.addEventListener("DOMContentLoaded", () => {
	const observer = new IntersectionObserver((entries) => {
		entries.forEach(entry => {
			if (entry.isIntersecting) {
				entry.target.classList.add('is-visible');
				observer.unobserve(entry.target);
			}
		});
	}, {
		root: null,
		rootMargin: '0px',
		threshold: 0.1
	});

	document.querySelectorAll('.reveal-on-scroll').forEach((el) => {
		observer.observe(el);
	});
});
