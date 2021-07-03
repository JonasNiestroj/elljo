import { currentComponent } from './lifecycleHooks'

export default function Observer(value, name) {
  if (!value || value.__observer__ || (!Array.isArray(value) && typeof value !== 'object')) {
    return
  }
  this.value = value
  value.__observer__ = this
  if (Array.isArray(value)) {
    for (var i = 0; i < value.length; i++) {
      new Observer(value[i], name)
    }
  } else if (typeof value === 'object' && value !== null) {
    const keys = Object.keys(value)
    for (var i = 0; i < keys.length; i++) {
      const key = keys[i]
      // Check if current property is configurable
      var prop = Object.getOwnPropertyDescriptor(value, key)
      if ((prop && !prop.configurable) || key === '__observer__') {
        continue
      }
      let keyValue = value[key];
      // TODO: Check for already existing getter/setter
      new Observer(keyValue, name)
      Object.defineProperty(value, key, {
        enumerable: true,
        configurable: true,
        get: function () {
          return keyValue
        },
        set: function (newValue) {
          currentComponent[name + 'IsDirty'] = true;
          currentComponent.oldState[name] = keyValue;
          currentComponent.queueUpdate();
          keyValue = newValue
          new Observer(newValue, name)
        }
      })
    }
  }
}