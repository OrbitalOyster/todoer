"use strict"

const initTogglePasswordButton = (btnSelector, inputSelector) => {
	/* Toggle password button */
	const btn = document.querySelector(btnSelector)
	if (!btn) throw new Error("Missing button element")
	/* Need the actual function here, to preserve 'this' */
	btn.addEventListener("click", function () {
		const input = document.querySelector(inputSelector)
		if (!input)
			throw new Error("Missing input element")
		if (input.type === "password")
			input.type = "text"
		else
			input.type = "password"
	})
}
