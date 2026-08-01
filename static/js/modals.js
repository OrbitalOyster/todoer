"use strict";

var modalEl, modal;

/* HTMX modal */
const showHTMXModal = () => {
	/* Remove previous HTMX content */
	const toRemoveQuery = "#modal > div:first-child > :not(.modal-placeholder)",
		toRemoveEls = document.querySelectorAll(toRemoveQuery);
	toRemoveEls.forEach((el) => el.remove());
	modal.show();
};

/* Confirm modal */
const confirmMsg = async (title, content) => {
	const confirmModal = new bootstrap.Modal("#confirm-modal"),
		modalEl = document.getElementById("confirm-modal"),
		titleEl = modalEl.querySelector("#confirm-modal-title"),
		contentEl = modalEl.querySelector("#confirm-modal-content"),
		okBtn = modalEl.querySelector("#confirm-modal-ok-btn");
	titleEl.textContent = title;
	contentEl.textContent = content;
	confirmModal.show();
	return new Promise((resolve) => {
		const confirmed = () => {
				resolve(true);
				cleanUp();
			},
			cleanUp = () => {
				okBtn.removeEventListener("click", confirmed);
				modalEl.removeEventListener("hidden.bs.modal", cleanUp);
				confirmModal.hide();
			};
		okBtn.addEventListener("click", confirmed);
		modalEl.addEventListener("hidden.bs.modal", cleanUp);
	});
};

/* For elements with "hx-confirm" tag */
const onHTMXConfirm = (event) => {
	/* No confirmation needed */
	if (!event.detail.question) return;
	/* Skip default action, show confirm modal */
	event.preventDefault();
	confirmMsg("Confirm action", event.detail.question)
		/* true to skip the built-in window.confirm() */
		.then((res) => res && event.detail.issueRequest(true));
};

/* Set up bootstrap modals */
const initModals = () => {
	modalEl = document.getElementById("modal");
	modal = new bootstrap.Modal("#modal");
	document.body.addEventListener("hideModal", () => modal.hide());
	document.addEventListener("htmx:confirm", onHTMXConfirm);
};
