function scrollToAnchor() {
	if (location && location.hash) {
		const el = document.getElementById(location.hash.replace('#', ''));
		if (el) {
			el.scrollIntoViewIfNeeded();
		}
	}
}

function registerObserver() {
	const io = new IntersectionObserver((entries) => {
		entries.forEach((entry) => {
			if (entry.intersectionRatio == 1) {
				for (const el of document.querySelectorAll(".active")) {
					el.classList.remove("active");
				}
				entry.target.classList.add("active");
				location.hash = entry.target.id;
			}
		});
	}, {
		root: document.querySelector("section"),
		threshold: 1,
	});

	// Declares what to observe, and observes its properties.
	const articles = document.querySelectorAll("section>article");
	articles.forEach((el) => {
		io.observe(el);
	});
}

function registerKeybindings() {
	document.addEventListener('keydown', (e) => {
		if (e.code === "Space" || e.code === "ArrowRight") {
			e.preventDefault();
			const next = document.querySelector(".active+article");
			if (next) {
				next.scrollIntoView({ behavior: "smooth" });
			}
		}
		if (e.code === "ArrowLeft") {
			e.preventDefault();
			const next = document.querySelector("article:has(+ .active)");
			if (next) {
				next.scrollIntoView({ behavior: "smooth" });
			}
		}
		if (e.code === "KeyF") {
			e.preventDefault();
			const section = document.querySelector("section");
			section.requestFullscreen();
		}
	});
}
