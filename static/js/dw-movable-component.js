import { DWComponent } from "./dw-component.js";
import { Component, EventBus } from "./framework.js" 

class DWMovableComponent extends DWComponent {
    constructor(id, parameters, callbacks) {
        super(id, parameters, callbacks);
        EventBus.register("up-click", this._handleUp.bind(this));
        EventBus.register("down-click", this._handleDown.bind(this));
    }

    _handleUp(event) {

    }

    _handleDown(event) {
        
    }

}

export { DWMovableComponent }