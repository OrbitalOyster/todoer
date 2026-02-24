"use strict";
(function () {

  function togglePassword() {
    const input = this.previousElementSibling
    if (!input)
      return
    if (input.type === "password") {
      // this.classList.remove("bi-eye-fill")
      // this.classList.add("bi-eye-slash-fill")
      input.type = "text"
    }
    else
      input.type = "password"
  }

  const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]')
  const tooltipList = [...tooltipTriggerList].map(tooltipTriggerEl => new bootstrap.Tooltip(tooltipTriggerEl))

  /* Toggle password buttons */
  const togglePasswordButtons = document.getElementsByClassName("toggle-password-btn")
  for (let i = 0; i < togglePasswordButtons.length; i++)
    togglePasswordButtons[i].addEventListener("click", togglePassword)

})()
