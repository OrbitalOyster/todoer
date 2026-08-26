"use strict"

const removeActiveTooltips = () => {
	document.querySelectorAll(".tooltip").forEach((e) => e.remove())
}

const updateTooltips = (el) => {
	/* Find all elements with tooltips */
	const tooltipTriggerList = [...el.querySelectorAll('[data-bs-toggle="tooltip"]')]
	/* Create tooltip */
	tooltipTriggerList.map((e) => new bootstrap.Tooltip(e))
}

const initTooltips = () => {
	updateTooltips(document);
	/* Update tooltips after any htmx DOM swap */
	document.addEventListener("htmx:beforeRequest", removeActiveTooltips)
	document.addEventListener("htmx:afterSettle", (e) => updateTooltips(e.detail.elt))
}
