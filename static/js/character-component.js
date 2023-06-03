import { Component } from "./framework.js" 
import { EventBus } from "./framework.js";

class CharacterComponent extends Component {
    constructor(id, parameters, callbacks) {
        super(id, parameters, callbacks);
        this._displayed = false;
        this._init();
        EventBus.register("eye-click", this._handleEye.bind(this));
    }

    _handleEye(event) {
        if (event.id === "eye-"+this.id) {
            this._displayed = !this._displayed;
            this.render();
        }
    }

    _template() {
        if (this._displayed) {
            return `
                <div class="element">
                    <div class="image-container">
                        <img src="/static/img/placeholder.svg" alt="Image">
                    </div>
                    <div class="content">
                        <div>
                            <p>Name : ${this.name}</p>
                        </div>
                        <div>
                            <p>Permanent tokens (long term memory)</p>
                            <p>Description : ${this.description}</p>
                            <p>Personality : ${this.personality}</p>
                            <p>Scenario : ${this.scenario}</p>
                            <p>Non permanent tokens (short term memory)</p>
                            <p>Message examples : ${this.mes_example}</p>
                            <p>First message : ${this.first_mes}</p>
                        </div>
                    </div>
                    <div class="buttons buttons-right">
                        <button><i class="fa-solid fa-pen-to-square"></i></button>
                        <button id="eye-${this.id}"><i class="fa-regular fa-eye"></i></button>
                    </div>
                </div>
                `
        } else {
            return `
                <div class="element">
                    <div class="content">
                        <div>
                            <p>Name : ${this.name}</p>
                        </div>
                    </div>
                    <div class="buttons buttons-right">
                        <button><i class="fa-solid fa-pen-to-square"></i></button>
                        <button id="eye-${this.id}"><i class="fa-regular fa-eye-slash"></i></button>
                    </div>
                </div>
                `   
        }
    }
}

export { CharacterComponent }