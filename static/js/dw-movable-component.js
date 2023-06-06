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
        console.log("up");
    }

    _handleDown(event) {
        console.log("down");
    }

}

export { DWMovableComponent }