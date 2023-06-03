class Component {
    // parameters format : {"param":"value", ...}
    // bindings format : {"id" : {"event":function}, ...}
    constructor(id, parameters, callbacks) {
        this._id = id;
        this._callbacks = callbacks;
        for (var key in parameters) {
            this[key] = parameters[key];
        }
    }

    // overload this method to handle the subclasse's template
    _template() {
        return `
            <div>
                <p>Overload the _template method in subclass.</p>
            </div>
            `
    }

    _init() {
        this._element = this._createElementFromHTML(this._template());
        this._element.id = this._id;
        this._inDom = false;
    }

    _createElementFromHTML(htmlString) {
        var div = document.createElement('div');
        div.innerHTML = htmlString.trim();
        return div.firstChild;
    }

    _addCallbacks() {
        if (!this._inDom) {
            throw "Component not in DOM";
        }
        for (var key in this._callbacks) {
            var element = document.getElementById(key);
            if (!element) {
                continue;
            }
            for (var event in this._callbacks[key]) {
                element.addEventListener(event, this._callbacks[key][event]);
            }
        }
    }

    render() {
        var newElement = this._createElementFromHTML(this._template());
        this._element.parentNode.replaceChild(newElement, this._element);
        this._element = newElement;
        this._element.id = this._id;
        this._addCallbacks();
    }

    appendToDom(parent) {
        var parentElement = document.getElementById(parent);
        if (this._inDom) {
            this.removeFromDom();
        }
        this._inDom = true;
        parentElement.appendChild(this._element);
        this._addCallbacks();
    }

    prependToDom(parent) {
        var parentElement = document.getElementById(parent);
        if (this._inDom) {
          this.removeFromDom();
        }
        this._inDom = true;
        parentElement.insertBefore(this._element, parentElement.firstChild);
        this._addCallbacks();
      }

    removeFromDom() {
        this._inDom = false;
        this._element.parentNode.removeChild(this._element);
    }

}

const EventBus = {
    register(event, callback) {
        document.addEventListener(event, (e) => callback(e.detail));
    },
    dispatch(event, data) {
        document.dispatchEvent(new CustomEvent(event, { detail: data }));
    },
    remove(event, callback) {
        document.removeEventListener(event, callback);
    },
};

export { Component, EventBus };
