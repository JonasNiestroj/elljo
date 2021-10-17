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

export const beforeDestroy = (callback) => {
  if (!currentComponent) {
    return
  }
  currentComponent.$.beforeDestroy.push(callback)
}

export const mounted = (callback) => {
  if (!currentComponent) {
    return
  }
  currentComponent.$.mounted.push(callback)
}

export const update = (callback) => {
  if (!currentComponent) {
    return
  }
  currentComponent.$.update.push(callback)
}