import { DWComponent } from "./dw-component.js";
import { Component, EventBus } from "./framework.js" 

class DWMovableComponent extends DWComponent {
    constructor(id, parameters, callbacks = {}) {
        super(id, parameters, callbacks);
        let upId = "up-" + parameters.id;
        let downId = "down-" + parameters.id;
        this._addCallbacks({
            [upId]: {"click":this._handleUp.bind(this)},
            [downId]: {"click":this._handleDown.bind(this)}
        });
    }

    _handleUp(event) {
        let previousElement = this._element.previousElementSibling;
        if (previousElement) {
            this._element.parentNode.insertBefore(this._element, previousElement);
            this._saveOrder();
        }
    }

    _handleDown(event) {
        let nextElement = this._element.nextElementSibling;
        if (nextElement) {
            this._element.parentNode.insertBefore(nextElement, this._element);
            this._saveOrder();
        }
    }

    _saveOrder() {
        throw "Need to overload _saveOrder method in subclass"
    }

}

export { DWMovableComponent }