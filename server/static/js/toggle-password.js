"use strict";

const initTogglePasswordButtons = () => {
	const togglePasswordBtns = document.getElementsByClassName(
		"toggle-password-btn",
	);
	for (let i = 0; i < togglePasswordBtns.length; i++)
		/* Need actual function here, to preserve 'this' */
		togglePasswordBtns[i].addEventListener("click", function () {
			/* Toggle password must come after input */
			const input = document.getElementById("password")
			if (!input) return;
			if (input.type === "password") input.type = "text";
			else input.type = "password";
		});
};

initTogglePasswordButtons();
