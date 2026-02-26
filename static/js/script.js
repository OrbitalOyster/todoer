'use strict'

let confirmMsg = null,
  htmxConfirmMsg = null,
  addToast = null;

(function () {

  /* Set up bootstrap tooltips */
  const tooltipTriggerList = 
    document.querySelectorAll('[data-bs-toggle="tooltip"]');
  [...tooltipTriggerList].map(
    tooltipTriggerEl => new bootstrap.Tooltip(tooltipTriggerEl)
  )

  /* Set up bootstrap toasts */
  const toastElList = document.querySelectorAll('.toast');
  [...toastElList].map(toastEl => new bootstrap.Toast(toastEl))

  /* Toggle password buttons */
  const togglePasswordBtns = 
    document.getElementsByClassName('toggle-password-btn')
  for (let i = 0; i < togglePasswordBtns.length; i++)
    togglePasswordBtns[i].addEventListener('click', togglePassword)

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
      modalEl = document.getElementById('confirmModal'),
      titleEl = modalDiv.querySelector('#modalTitle'),
      contentEl = modalDiv.querySelector('#modalContent'),
      okBtn = modalEl.querySelector('#modalOkBtn'),
      cancelBtn = modalDiv.querySelector('#modalCancelBtn')
    titleEl.textContent = title
    contentEl.textContent = content
    modal.show()
    return new Promise((resolve) => {
      const ok = () => {
          resolve(true)
          cleanUp()
        },
        cancel = () => cleanUp(),
        cleanUp = () => {
          okBtn.removeEventListener('click', ok)
          cancelBtn.removeEventListener('click', cleanUp)
          modal.hide()
        }
      okBtn.addEventListener('click', ok)
      cancelBtn.addEventListener('click', cancel)
    })
  }

  /* Confirm modal for htmx events */
  htmxConfirmMsg = (el, title, content) => 
    confirmMsg(title, content)
      .then(res => res && htmx.trigger(el, 'confirmed'))

  addToast = (title, msg) => {
    const toastEl = document.createElement('div'),
      toastHeader = document.createElement('div'),
      toastTitle = document.createElement('strong'),
      toastCloseBtn = document.createElement('button'),
      toastBody = document.createElement('div')

    toastEl.className = 'toast'
    toastEl.dataset.bsAutohide = 'false'
    toastHeader.className = 'toast-header'
    toastTitle.textContent = title
    toastTitle.className = 'me-auto'
    toastCloseBtn.type = 'button'
    toastCloseBtn.className = 'btn-close'
    toastCloseBtn.dataset.bsDismiss = 'toast'
    toastBody.className = 'toast-body'
    toastBody.textContent = msg

    toastEl.appendChild(toastHeader)
    toastHeader.appendChild(toastTitle)
    toastHeader.appendChild(toastCloseBtn)
    toastEl.appendChild(toastBody)

    document.querySelector('.toast-container').appendChild(toastEl)
    new bootstrap.Toast(toastEl).show()
    toastEl.addEventListener('hidden.bs.toast', () => toastEl.remove())
  }

})()
