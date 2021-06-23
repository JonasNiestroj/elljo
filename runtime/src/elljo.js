export class EllJo {

  plugins = []
  components = {}

  constructor() {
    if (!window.__elljo__) {
      window.__elljo__ = this
    }
  }

  addPlugin(plugin, options) {
    if (!plugin.init) {
      // TODO: Log error
      return
    }
    plugin.init(this, options)
  }

  addComponent(name, component) {
    this.components[name] = component
  }

  mount(component, to) {
    const toElement = document.querySelector(to)
    if (!toElement) {
      // TODO: Log error
      return
    }
    new component({ target: toElement })
  }
}