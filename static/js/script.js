'use strict'

let confirmMsg = null,
  htmxConfirmMsg = null,
  addToast = null,
  showEditTaskModal = null;

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
    const confirmModal = new bootstrap.Modal('#confirmModal'),
      modalEl = document.getElementById('confirmModal'),
      titleEl = modalEl.querySelector('#confirmModalTitle'),
      contentEl = modalEl.querySelector('#confirmModalContent'),
      okBtn = modalEl.querySelector('#confirmModalOkBtn'),
      cancelBtn = modalEl.querySelector('#confirmModalCancelBtn')
    titleEl.textContent = title
    contentEl.textContent = content
    confirmModal.show()
    return new Promise((resolve) => {
      const ok = () => {
          resolve(true)
          cleanUp()
        },
        cancel = () => cleanUp(),
        cleanUp = () => {
          okBtn.removeEventListener('click', ok)
          cancelBtn.removeEventListener('click', cleanUp)
          confirmModal.hide()
        }
      okBtn.addEventListener('click', ok)
      cancelBtn.addEventListener('click', cancel)
    })
  }

  /* Confirm modal for htmx events */
  htmxConfirmMsg = (el, title, content) => 
    confirmMsg(title, content)
      .then(res => res && htmx.trigger(el, 'confirmed'))

  /* Create new toast message */
  addToast = (type, title, msg) => {
    const toastEl = document.createElement('div'),
      toastHeader = document.createElement('div'),
      toastIcon = document.createElement('i'),
      toastTitle = document.createElement('strong'),
      toastTime = document.createElement('small'),
      toastCloseBtn = document.createElement('button'),
      toastBody = document.createElement('div')

    toastEl.className = 'toast'
    toastEl.dataset.bsAutohide = 'false'
    toastHeader.className = 'toast-header'
    toastIcon.classList.add('bi', 'me-2')
    toastTitle.textContent = title
    toastTitle.className = 'me-auto'
    toastTime.textContent = new Date().toLocaleString()
    toastCloseBtn.type = 'button'
    toastCloseBtn.className = 'btn-close'
    toastCloseBtn.dataset.bsDismiss = 'toast'
    toastBody.className = 'toast-body'
    toastBody.textContent = msg

    switch (type) {
      case 'info':
        toastEl.classList.add('border-info-subtle')
        toastHeader.classList.add('bg-info-subtle')
        toastIcon.classList.add('text-info', 'bi-info-circle-fill')
        break
      case 'warning':
        toastEl.classList.add('border-warning-subtle')
        toastHeader.classList.add('bg-warning-subtle')
        toastIcon.classList.add('text-warning', 'bi-exclamation-triangle-fill' )
        break
      case 'danger':
        toastEl.classList.add('border-danger-subtle')
        toastHeader.classList.add('bg-danger-subtle')
        toastIcon.classList.add('text-danger', 'bi-x-octagon-fill')
        break
      default:
        break
    }

    toastEl.appendChild(toastHeader)
    toastHeader.appendChild(toastIcon)
    toastHeader.appendChild(toastTitle)
    toastHeader.appendChild(toastTime)
    toastHeader.appendChild(toastCloseBtn)
    toastEl.appendChild(toastBody)

    document.querySelector('.toast-container').appendChild(toastEl)
    new bootstrap.Toast(toastEl).show()
    toastEl.addEventListener('hidden.bs.toast', () => toastEl.remove())
  }

  showEditTaskModal = (taskId) => {
    const editTaskModal = new bootstrap.Modal('#editTaskModal'),
      modalEl = document.getElementById('editTaskModal'),
      inputEl = modalEl.querySelector('#taskDescriptionInput')
    inputEl.textContent = taskId
      // okBtn = modalEl.querySelector('#confirmModalOkBtn'),
      // cancelBtn = modalEl.querySelector('#confirmModalCancelBtn')
    //contentEl.textContent = content
    editTaskModal.show()
  }

  /* HTMX triggers */
  function onToast(event) {
    const type = event.detail.type,
      title = event.detail.title,
      msg = event.detail.msg
    addToast(type, title, msg) 
  }

  document.body.addEventListener("toast", onToast)

})()
