"use strict";

/* Set up bootstrap modals */
const modalEl = document.getElementById("modal"),
	modal = new bootstrap.Modal("#modal");
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
/* Confirm modal for htmx events */
const htmxConfirmMsg = (el, title, content) =>
	confirmMsg(title, content).then(
		(res) => res && htmx.trigger(el, "confirmed"),
	);
/* HTMX modal */
const showHTMXModal = () => {
	/* Remove previous HTMX content */
	const toRemoveQuery = "#modal > div:first-child > :not(.modal-placeholder)",
		toRemoveEls = document.querySelectorAll(toRemoveQuery);
	toRemoveEls.forEach((el) => el.remove());
	modal.show();
};

document.body.addEventListener("hideModal", () => modal.hide());
