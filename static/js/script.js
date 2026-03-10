'use strict'

let confirmMsg = null,
  htmxConfirmMsg = null,
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

  /* On toast */
  document.body.addEventListener('toast', function() {
    const toastEl = document.querySelector('.toast-container').lastElementChild
    new bootstrap.Toast(toastEl).show()
    const progressBar = toastEl.querySelector('.progress-bar')
    if (progressBar) {
      progressBar.style = 'width: 100%; transition: width linear 10s'
      /* Halt progress bar animation on mouse over */
      toastEl.addEventListener(
        'mouseover',
        () => progressBar.style = 'width: 0%; transition: none'
      )
      toastEl.addEventListener(
        'mouseleave',
        () => progressBar.style = 'width: 100%; transition: width linear 10s'
      )
    }
    /* Remove element after delay */
    toastEl.addEventListener('hidden.bs.toast', () => toastEl.remove())
  })

})()
