"use strict";

const initPopover = (hostSelector, contentSelector) => {
	const hostEl = document.querySelector(hostSelector);
	if (!hostEl) throw new Error("Missing host element");
	const contentEl = document.querySelector(contentSelector);
	if (!contentEl) throw new Error("Missing content element");

	const popover = new bootstrap.Popover(hostEl, {
		html: true,
		sanitize: false,
		content: () => contentEl.innerHTML,
	});

	hostEl.addEventListener('shown.bs.popover', () => {
		const content = document.querySelector(".popover")
		htmx.process(content)
	})
};
