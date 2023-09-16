import prefix from "../../config.js";
import { DWComponent } from "./dw-component.js";
import { EventBus } from "./framework.js";

class CharacterComponent extends DWComponent {
    constructor(id, parameters, callbacks = {}) {
        super(id, parameters, callbacks);
        this._uri = prefix + '/api/character/';
        let clone = "clone-" + parameters.id;
        this._addCallbacks({
            [clone]: {"click":this._handleClone.bind(this)},
        });
    }

    _handleClone(event) {
        event.preventDefault();
        fetch(prefix  + "/api/clone/" + this.id)
        .then(function(response) {
            if (response.ok) {
                EventBus.dispatch("refresh");
            } else {
                alert("Save failed.");
            }
        })
        .catch(error => {
            alert("An error occurred. Please try again."+error);
        });
    }

    _template() {
        if (this._edition) {
            return `
            <div class="element">
                <div class="image-container">
                    <img src="${prefix}/api/avatar/${this.id}" alt="Image">
                </div>
                <div class="content">
                    <form id="form-${this.id}" method="POST" action="${prefix}/api/character/${this.id}">
                        <input type="hidden" name="id" value="${this.id}">
                        <div>
                            <p>Name : </p><input type="text" value="${this.name}" name="name" class="custom-input">
                            <label>Is human: </label><input type="checkbox" id="is_human-${this.id}" name="is_human" ${this.is_human ? "checked" : ""}>
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
                    <button id="edit-${this.id}"><i class="fa-solid fa-arrow-up-right-from-square"></i></button>
                    <button id="open-${this.id}"><i class="fa-solid fa-folder-open"></i></button>
                </div>
            </div>
            `
        }
        if (this.id === 0) {
            return `
                <div class="element">
                    <div class="content">
                        <div>
                            <p>Create a new character</p>
                        </div>
                    </div>
                    <div class="buttons buttons-right">
                        <button id="edit-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
                        <button id="open-${this.id}"><i class="fa-solid fa-folder-closed"></i></button>
                    </div>
                </div>
                `
        }
        if (this._displayed) {
                return `
                    <div class="element">
                        <div class="image-container">
                            <img src="${prefix}/api/avatar/${this.id}" alt="Image">
                        </div>
                        <div class="content">
                            <div>
                                <p>Name : ${this.name}</p>
                                <p>Is human : ${this.is_human ? "Yes" : "No"}</p>
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
                            <button id="clone-${this.id}" class="custom-button ml2 mt2">Clone character</button>
                        </div>
                        <div class="buttons buttons-right">
                            <button id="edit-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
                            <button id="open-${this.id}"><i class="fa-solid fa-folder-open"></i></button>
                        </div>
                    </div>
                    `
        }
        return `
            <div class="element">
                <div class="content">
                    <div>
                        <p>Name : ${this.name}</p>
                    </div>
                </div>
                <div class="buttons buttons-right">
                    <button id="edit-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
                    <button id="open-${this.id}"><i class="fa-solid fa-folder-closed"></i></button>
                </div>
            </div>
            `
    }
}

export { CharacterComponent }