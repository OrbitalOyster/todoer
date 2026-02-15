'use strict'

function onFormSubmit(formId, callback) {
  const element = document.getElementById(formId)
  if (!element)
    throw new Error("Oh noes!")
  element.addEventListener("submit", event => {
    event.preventDefault()
    const form = event.target,
      formData = new FormData(form),
      formObject = Object.fromEntries(formData.entries())
    callback(formObject)
  })
}
