import { setComponent } from './lifecycleHooks';
import Observer from './observer';

export default class EllJoComponent {
  constructor(options, props, events) {
    this.$ = {};
    this.$.afterRender = [];
    this.$.beforeDestroy = [];
    this.$.mounted = [];
    this.$.update = [];
    this.$props = {};
    this.$events = {};
    this.oldState = {};
    this.updating = false;
    setComponent(this);

    if (events) {
      Object.keys(events).forEach((event) => {
        if (!this.$events[event]) {
          this.$events[event] = [events[event]];
        } else {
          this.$events[event].push(events[event]);
        }
      });
    }

    this.$event = (name) => {
      var callbacks = this.$events[name];
      if (callbacks) {
        const args = [];
        for (let i = 1; i < arguments.length; i++) {
          args.push(arguments[i]);
        }
        callbacks.forEach((callback) => callback(...args));
      }
    };

    this.utils = {
      diffArray: function diffArray(one, two) {
        if (!Array.isArray(two)) {
          return one.slice();
        }

        var tlen = two.length;
        var olen = one.length;
        var idx = -1;
        var arr = [];

        while (++idx < olen) {
          var ele = one[idx];

          var hasEle = false;
          for (var i = 0; i < tlen; i++) {
            var val = two[i];

            if (ele === val) {
              hasEle = true;
              break;
            }
          }

          if (hasEle === false) {
            arr.push({ element: ele, index: idx });
          }
        }
        return arr;
      },
    };

    const propertyNames = Object.getOwnPropertyNames(this);
    for (let i = 0; i < propertyNames.length; i++) {
      const property = propertyNames[i];
      if (Array.isArray(this[property])) {
        patchArray(this[property], property);
        new Observer(this[property], property);
      } else {
        new Observer(this[property], property);
      }
    }

    function patchArray(array, name) {
      const methodsToPatch = [
        'push',
        'pop',
        'splice',
        'sort',
        'reverse',
        'shift',
        'unshift',
        'fill',
      ];
      methodsToPatch.forEach((method) => {
        const currentMethod = array[method];
        Object.defineProperty(array, method, {
          enumerable: false,
          configurable: false,
          writable: false,
          value: function () {
            const result = currentMethod.apply(this, arguments);
            this[name + 'IsDirty'] = true;
            this.oldState[name] = array;
            this.queueUpdate();
            new Observer(result, name);
            return result;
          },
        });
      });
    }
  }

  update() {
    const callbacks = this.$.update;
    for (let i = 0; i < callbacks.length; i++) {
      callbacks[i]();
    }
    this.updating = false;
    this.$.mainFragment.update();
    this.oldState = {};
  }

  queueUpdate() {
    if (!this.updating) {
      this.updating = true;
      Promise.resolve().then(() => this.update());
    }
  }

  teardown() {
    const callbacks = this.$.beforeDestroy;
    for (let i = 0; i < callbacks.length; i++) {
      callbacks[i]();
    }
    this.$.mainFragment.teardown();
    this.$.mainFragment = null;
  }
}
