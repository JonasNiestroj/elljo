import { setComponent, currentComponent } from './lifecycleHooks';

export default class EllJoComponent {
  constructor(options, props, events) {
    this.$ = {};
    this.$.afterRender = [];
    this.$.beforeDestroy = [];
    this.$.mounted = [];
    this.$.update = [];
    this.$props = {};
    this.$propsBindings = {};
    this.$events = {};
    this.$slots = {};
    this.oldState = {};
    this.updating = false;

    if (options.slots) {
      this.$slots = options.slots
    }

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
  }

  updateValue(name, func) {
    currentComponent[name + 'IsDirty'] = true;
    if (this.$propsBindings[name]) {
      for (let i = 0; i < this.$propsBindings[name].length; i++) {
        this[this.$propsBindings[name][i]].$props[name] = func;
      }
    }

    this.queueUpdate();
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
