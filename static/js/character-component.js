import { Component } from "./framework.js" 
import { EventBus } from "./framework.js";

class CharacterComponent extends Component {
    constructor(id, parameters, callbacks) {
        super(id, parameters, callbacks);
        this._displayed = false;
        this._edition = false;
        this._edited = false;
        this._init();
        EventBus.register("open-click", this._handleopen.bind(this));
        EventBus.register("pen-click", this._handlePen.bind(this));
        EventBus.register("saved", this._handleSaved.bind(this));
        EventBus.register("canceled", this._handleCanceled.bind(this));
    }

    _handleopen(event) {
        if (event.id === "open-"+this.id) {
            this._displayed = !this._displayed;
            if (this._edition === true) {
                this._edition = false;
            }
            this.render();
        }
    }

    _handlePen(event) {
        if (event.id === "pen-"+this.id) {
            this._edition = !this._edition;
            if (this._displayed === false) {
                this._displayed = true;
            }
            this.render();
            if (this._edition === true) {
                const formElement = document.getElementById("edit-"+this.id);
                const formInputs = formElement.querySelectorAll('input, select, textarea');

                formInputs.forEach(input => {
                    input.addEventListener('input', this._handleInput.bind(this));
                });
            }
        }
    }

    _handleInput(event) {
        console.log("input");
        if (this._edited === false) {
            this._edited = true;
            const openButton = document.getElementById("open-"+this.id);
            openButton.disabled = true;
            const editButton = document.getElementById("pen-"+this.id);
            editButton.disabled = true;
        }
    }

    _handleSaved(event) {
        if (event.id == this.id) {
            this._refresh();
            this._edition = false;
            this._edited = false;
            this._displayed = false;    
            this.render();
        }
    }

    _handleCanceled(event) {
        if (event.id == this.id) {
            this._edition = false;
            this._edited = false;
            this._displayed = false;
            this.render();
        }
    }

    _refresh() {
        var form = document.getElementById("edit-"+this.id);
        var formData = new FormData(form);
        var self = this;
        formData.forEach(function(value, key) {
            if (key === "id") {
                self[key] = parseInt(value);
            } else {
                self[key] = value;
            }
        }, this);
    }

    _template() {
        if (this._displayed) {
            if (this._edition) {
                return `
                <div class="element">
                    <div class="image-container">
                        <img src="/api/avatar/${this.id}" alt="Image">
                    </div>
                    <div class="content">
                        <form id="edit-${this.id}" method="POST" action="/api/character/${this.id}">
                            <input type="hidden" name="id" value="${this.id}">
                            <div>
                                <p>Name : </p><input type="text" value="${this.name}" name="name" class="custom-input">
                            </div>
                            <div class="mt2">
                                <p>Permanent tokens (long term memory) :</p>
                                <label>Description: </label><textarea name="description" class="custom-textarea">${this.description}</textarea>
                                <label>Personality: </label><textarea name="personality" class="custom-textarea">${this.personality}</textarea>
                                <label>Scenario: </label><textarea name="scenario" class="custom-textarea">${this.scenario}</textarea>
                            </div>
                            <div class="mt2">
                                <p>Temporary  tokens (short term memory) :</p>
                                <label>Message examples: </label><textarea name="mes_example" class="custom-textarea">${this.mes_example}</textarea>
                                <label>First message: </label><textarea name="first_mes" class="custom-textarea">${this.first_mes}</textarea>
                            </div>
                            <button id="save-${this.id}" type="submit" class="custom-button ml2 mt2">Save</button>
                            <button id="cancel-${this.id}" type="button" class="custom-button mt2">Cancel</button>
                        </form>
                    </div>
                    <div class="buttons buttons-right">
                        <button id="pen-${this.id}"><i class="fa-solid fa-arrow-up-right-from-square"></i></button>
                        <button id="open-${this.id}"><i class="fa-solid fa-folder-open"></i></button>
                    </div>
                </div>
                `
            } else {
                return `
                    <div class="element">
                        <div class="image-container">
                            <img src="/api/avatar/${this.id}" alt="Image">
                        </div>
                        <div class="content">
                            <div>
                                <p>Name : ${this.name}</p>
                            </div>
                            <div>
                                <p>Permanent tokens (long term memory) :</p>
                                <p>Description : ${this.description}</p>
                                <p>Personality : ${this.personality}</p>
                                <p>Scenario : ${this.scenario}</p>
                                <p>Temporary tokens (short term memory) :</p>
                                <p>Message examples : ${this.mes_example}</p>
                                <p>First message : ${this.first_mes}</p>
                            </div>
                        </div>
                        <div class="buttons buttons-right">
                            <button id="pen-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
                            <button id="open-${this.id}"><i class="fa-solid fa-folder-open"></i></button>
                        </div>
                    </div>
                    `
            }
        } else {
            return `
                <div class="element">
                    <div class="content">
                        <div>
                            <p>Name : ${this.name}</p>
                        </div>
                    </div>
                    <div class="buttons buttons-right">
                        <button id="pen-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
                        <button id="open-${this.id}"><i class="fa-solid fa-folder-closed"></i></button>
                    </div>
                </div>
                `   
        }
    }
}

export { CharacterComponent }