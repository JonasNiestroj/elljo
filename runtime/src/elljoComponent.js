import { setComponent } from './lifecycleHooks'

export default class EllJoComponent {

  constructor(options, props, events) {
    this.$ = {};
    this.$.afterRender = [];
    this.$props = {};
    this.$events = {};
    this.oldState = {};
    this.upading = false
    setComponent(this);

    if (events) {
      Object.keys(events).forEach(event => {
        if (!this.$events[event]) {
          this.$events[event] = [events[event]]
        } else {
          this.$events[event].push(events[event])
        }
      })
    }
    
    this.$event = (name) => {
      var callbacks = this.$events[name]
      if (callbacks) {
        const args = []
        for (let i = 1; i < arguments.length; i++) {
          args.push(arguments[i])
        }
        callbacks.forEach(callback => callback(...args))
      }
    }
    
    this.utils = {
      diffArray: function diffArray(one, two) {
        if (!Array.isArray(two)) {
          return one.slice();
        }

        var tlen = two.length
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
      }
    }
  }

  queueUpdate() {
		if(!this.updating) {
			this.updating = true;
			Promise.resolve().then(this.update)
		}
  }

  teardown() {
    this.mainFragment.teardown();
    mainFragment = null;
  }
}