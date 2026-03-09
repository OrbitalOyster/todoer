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
      toastProgress = document.createElement('div'),
      toastProgressBar = document.createElement('div'),
      toastCloseBtn = document.createElement('button'),
      toastBody = document.createElement('div')

    toastEl.className = 'toast'
    toastHeader.className = 'toast-header'
    toastIcon.classList.add('bi', 'me-2')
    toastTitle.textContent = title
    toastTitle.className = 'me-auto'
    toastTime.textContent = new Date().toLocaleString()
    
    toastProgress.className = 'progress'
    toastProgress.style = 'height: 3px'
    toastProgressBar.classList.add('progress-bar')
    toastProgressBar.style = 'width: 0%'

    toastCloseBtn.type = 'button'
    toastCloseBtn.className = 'btn-close'
    toastCloseBtn.dataset.bsDismiss = 'toast'
    toastBody.className = 'toast-body'
    toastBody.textContent = msg

    switch (type) {
      case 'success':
        toastEl.classList.add('border-success-subtle')
        toastHeader.classList.add('bg-success-subtle')
        toastIcon.classList.add('text-success', 'bi-hand-thumbs-up-fill')
        toastProgressBar.classList.add('bg-success')
        toastEl.dataset.bsAutohide = 'true'
        toastEl.dataset.bsDelay = '10000'
        break
      case 'info':
        toastEl.classList.add('border-info-subtle')
        toastHeader.classList.add('bg-info-subtle')
        toastIcon.classList.add('text-info', 'bi-info-circle-fill')
        toastProgressBar.classList.add('bg-info')
        toastEl.dataset.bsAutohide = 'true'
        toastEl.dataset.bsDelay = '10000'
        break
      case 'warning':
        toastEl.classList.add('border-warning-subtle')
        toastHeader.classList.add('bg-warning-subtle')
        toastIcon.classList.add('text-warning', 'bi-exclamation-triangle-fill' )
        toastEl.dataset.bsAutohide = 'false'
        toastProgress.className = 'd-none'
        break
      case 'danger':
        toastEl.classList.add('border-danger-subtle')
        toastHeader.classList.add('bg-danger-subtle')
        toastIcon.classList.add('text-danger', 'bi-x-octagon-fill')
        toastEl.dataset.bsAutohide = 'false'
        toastProgress.className = 'd-none'
        break
      default:
        break
    }

    toastEl.appendChild(toastHeader)
    toastHeader.appendChild(toastIcon)
    toastHeader.appendChild(toastTitle)
    toastHeader.appendChild(toastTime)
    toastProgress.appendChild(toastProgressBar)
    toastEl.appendChild(toastProgress)
    toastHeader.appendChild(toastCloseBtn)
    toastEl.appendChild(toastBody)

    toastEl.addEventListener(
      'mouseover', () => toastProgressBar.style = 'width: 0%; transition: none'
    )

    toastEl.addEventListener(
      'mouseleave', () => toastProgressBar.style = 'width: 100%; transition: width linear 10s'
    )

    toastEl.addEventListener('hidden.bs.toast', () => toastEl.remove())
    document.querySelector('.toast-container').appendChild(toastEl)
    new bootstrap.Toast(toastEl).show()
    toastProgressBar.style = 'width: 100%; transition: width linear 10s'
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

  /* On toast */
  function onToast() {
    const toastEl = document.querySelector('.toast-container').lastElementChild
    new bootstrap.Toast(toastEl).show()
    const progressBar = toastEl.querySelector('.progress-bar')
    if (progressBar)
      progressBar.style = 'width: 100%'
    toastEl.addEventListener('hidden.bs.toast', () => toastEl.remove())
  }
  document.body.addEventListener('toast', onToast)

})()
