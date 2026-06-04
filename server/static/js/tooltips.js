"use strict";

const removeActiveTooltips = () => {
	const activeTooltips = document.querySelectorAll(".tooltip");
	activeTooltips.forEach((t) => t.remove());
};

const updateTooltips = (el) => {
	const tooltipTriggerList = el.querySelectorAll('[data-bs-toggle="tooltip"]');
	[...tooltipTriggerList].map(
		(tooltipTriggerEl) => new bootstrap.Tooltip(tooltipTriggerEl),
	);
};

const initTooltips = () => {
	/* Update tooltips after any htmx DOM swap */
	document.addEventListener("htmx:beforeRequest", removeActiveTooltips);
	document.addEventListener("htmx:afterSettle", (e) =>
		updateTooltips(e.detail.elt),
	);
	updateTooltips(document);
};
