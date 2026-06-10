"use strict";

const createSelectAllGroup = (masterSelector, childrenSelector) => {
	const masterEl = document.querySelector(masterSelector);
	if (!masterEl) return;
	/* Convert whatever querySelectorAll returns to an actual array */
	const childrenEls = [...document.querySelectorAll(childrenSelector)],
		allSelected = () => childrenEls.every((e) => e.checked),
		noneSelected = () => childrenEls.every((e) => !e.checked),
		selectAll = () => childrenEls.forEach((e) => (e.checked = true)),
		selectNone = () => childrenEls.forEach((e) => (e.checked = false));
	/* On master change */
	masterEl.onchange = () => (masterEl.checked ? selectAll() : selectNone());
	/* On children change */
	const updateMaster = () => {
		masterEl.indeterminate = false;
		if (allSelected()) masterEl.checked = true;
		else if (noneSelected()) masterEl.checked = false;
		else masterEl.indeterminate = true;
	};
	childrenEls.forEach((el) => (el.onchange = updateMaster));
	updateMaster();
};
