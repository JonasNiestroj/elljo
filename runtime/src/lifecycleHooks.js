export let currentComponent = null;

export const setComponent = (component) => {
  currentComponent = component
}

export const afterRender = (callback) => {
  if (!currentComponent) {
    return
  }
  currentComponent.$.afterRender.push(callback)
}