"use strict";

const initTogglePasswordButton = (btnSelector, inputSelector) => {
	const toggleBtn = document.querySelector(btnSelector);
	if (!toggleBtn) throw new Error("Missing element");
	/* Need actual function here, to preserve 'this' */
	toggleBtn.addEventListener("click", function () {
		const input = document.querySelector(inputSelector);
		if (!input) throw new Error("Missing element");
		if (input.type === "password") input.type = "text";
		else input.type = "password";
	});
};
