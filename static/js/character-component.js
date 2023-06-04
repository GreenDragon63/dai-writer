import { DaiComponent } from "./dai-component.js";

class CharacterComponent extends DaiComponent {
    constructor(id, parameters, callbacks) {
        super(id, parameters, callbacks);
        this._uri = '/api/character/';
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
                        <form id="form-${this.id}" method="POST" action="/api/character/${this.id}">
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
                        <button id="edit-${this.id}"><i class="fa-solid fa-arrow-up-right-from-square"></i></button>
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
                            <button id="edit-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
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
                        <button id="edit-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
                        <button id="open-${this.id}"><i class="fa-solid fa-folder-closed"></i></button>
                    </div>
                </div>
                `   
        }
    }
}

export { CharacterComponent }