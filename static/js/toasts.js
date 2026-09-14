"use strict"

/* On toast */
const onToast = () => {
	const toastContainer = document.getElementById("toast-container"),
		toastEls = [...toastContainer.querySelectorAll(".toast:not(.fade)")]
	/* For every new toast */
	toastEls.forEach((toastEl) => {
		/* Show it */
		new bootstrap.Toast(toastEl).show()
		/* Handle funky stuff */
		const progressBar = toastEl.querySelector(".progress-bar")
		if (progressBar) {
			progressBar.style = "width: 100%; transition: width linear 10s"
			/* Halt progress bar animation on mouse over */
			toastEl.addEventListener(
				"mouseover",
				() => (progressBar.style = "width: 0%; transition: none"),
			)
			/* Resume on mouse leave */
			toastEl.addEventListener(
				"mouseleave",
				() => (progressBar.style = "width: 100%; transition: width linear 10s"),
			)
		}
		/* Remove element after delay */
		toastEl.addEventListener("hidden.bs.toast", toastEl.remove)
	})
}

const initToasts = () => {
	document.body.addEventListener("toast", onToast)
}
