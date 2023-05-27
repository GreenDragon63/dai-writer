class Component {
    // parameters format : {"param":"value", ...}
    // callbacks format : {"id" : {"event":function}, ...}
    constructor(id, parameters, callbacks) {
        this._id = id;
        this._callbacks = callbacks;
        for (var key in parameters) {
            this[key] = parameters[key];
            console.log(key + " : " + this[key]);
        }
        this._element = this._createElementFromHTML(this._template());
        this._element.id = id;
        this._inDom = false;
    }

    // overload this method to handle the subclasse's template
    _template() {
        return `
            <div>
                <p>Overload the _template method in subclass. Param : ${this.param} </p>
                <button id="button-${this._id}">Click me</button>
            </div>
            `
    }

    _createElementFromHTML(htmlString) {
        var div = document.createElement('div');
        div.innerHTML = htmlString.trim();
        return div.firstChild;
    }

    _addCallbacks() {
        if (!this._inDom) {
            throw "Componant not in DOM";
        }
        for (var key in this._callbacks) {
            var element = document.getElementById(key+"-"+this._id);
            for (var event in this._callbacks[key]) {
                element.addEventListener(event, this._callbacks[key][event]);
            }
        }
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

    removeFromDom() {
        this._inDom = false;
        this._element.parentNode.removeChild(this._element);
    }

}

c = new Component("testid", {"param":"test"}, {"button":{"click":() => alert("Working")}});
c.appendToDom("container");
c.removeFromDom();
c.appendToDom("container");
