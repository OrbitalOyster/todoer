"use strict";
(function () {

  function togglePassword() {
    const input = this.previousElementSibling
    if (!input)
      return
    if (input.type === "password")
      input.type = "text"
    else
      input.type = "password"
  }

  /* Toggle password buttons */
  const togglePasswordButtons = document.getElementsByClassName("toggle-password-btn")
  for (let i = 0; i < togglePasswordButtons.length; i++)
    togglePasswordButtons[i].addEventListener("click", togglePassword)

})()
