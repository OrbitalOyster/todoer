"use strict";

/* On toast */
const initToasts = () => {
	document.body.addEventListener("toast", () => {
		/* Get new toast */
		const toastEl = document.querySelector(".toast-container").lastElementChild;
		/* Show it */
		new bootstrap.Toast(toastEl).show();
		/* Handle funky stuff */
		const progressBar = toastEl.querySelector(".progress-bar");
		if (progressBar) {
			progressBar.style = "width: 100%; transition: width linear 10s";
			/* Halt progress bar animation on mouse over */
			toastEl.addEventListener(
				"mouseover",
				() => (progressBar.style = "width: 0%; transition: none"),
			);
			/* Resume on mouse leave */
			toastEl.addEventListener(
				"mouseleave",
				() => (progressBar.style = "width: 100%; transition: width linear 10s"),
			);
		}
		/* Remove element after delay */
		toastEl.addEventListener("hidden.bs.toast", toastEl.remove);
	});
};

initToasts();
