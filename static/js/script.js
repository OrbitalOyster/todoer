'use strict'

let confirmMsg = null,
  htmxConfirmMsg = null;

(function () {

  /* Set up bootstrap tooltips */
  const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]');
  [...tooltipTriggerList].map(tooltipTriggerEl => new bootstrap.Tooltip(tooltipTriggerEl))

  /* Toggle password buttons */
  const togglePasswordButtons = document.getElementsByClassName('toggle-password-btn')
  for (let i = 0; i < togglePasswordButtons.length; i++)
    togglePasswordButtons[i].addEventListener('click', togglePassword)

  function togglePassword() {
    const input = this.previousElementSibling
    if (!input)
      return
    if (input.type === 'password')
      input.type = 'text'
    else
      input.type = 'password'
  }

  /* Confirm modal */
  confirmMsg = async (title, content) => {
    const modal = new bootstrap.Modal('#confirmModal'),
      modalDiv = document.getElementById('confirmModal'),
      titleH = modalDiv.querySelector('#modalTitle'),
      contentDiv = modalDiv.querySelector('#modalContent'),
      okButton = modalDiv.querySelector('#modalOkBtn'),
      cancelButton = modalDiv.querySelector('#modalCancelBtn')
    titleH.textContent = title
    contentDiv.textContent = content
    modal.show()
    return new Promise((resolve) => {
      const ok = () => {
          resolve(true)
          cleanUp()
        },
        cancel = () => cleanUp(),
        cleanUp = () => {
          okButton.removeEventListener('click', ok)
          cancelButton.removeEventListener('click', cleanUp)
          modal.hide()
        }
      okButton.addEventListener('click', ok)
      cancelButton.addEventListener('click', cancel)
    })
  }

  // htmxConfirmMsg = (title, content) => {
  //   confirmMsg(title, content).then(res => res && htmx.trigger(this, 'confirmed'))
  // }

})()
