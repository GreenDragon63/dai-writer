import { DWMovableComponent } from "./dw-movable-component.js";
import selectCharacter from "./select-character.js";

class SceneComponent extends DWMovableComponent {
    constructor(id, parameters, callbacks) {
        super(id, parameters, callbacks);
        if (this.characters === undefined) {
            this.characters = [];
        }
        let addId = "add-" + parameters.id;
        this._addCallbacks({
            [addId]: {"click":this._handleAdd.bind(this)},
        });
        this._uri = '/api/scene/' + this.book_id + '/';
        if ((this.lines === undefined) || (this.lines === null)) {
            this.lines = [];
        }
        this.numLines = this.lines.length;
    }

    _handleAdd(event) {
        event.preventDefault();
        let element = document.getElementById("character-"+this.id);
        let charId = parseInt(element.value)
        if (!this.characters.includes(charId)) {
            this.characters.push(charId);
        }
        this.render();
        console.log(this.characters)
    }

    _handleRemove(event) {
        event.preventDefault();
    }

    _template() {
        if (this.id === 0) {
            var linkLines = "";
        } else {
            var linkLines = `<a href="/line/${this.book_id}/${this.id}/" class="custom-button button-link ml2 mt2">Edit scenes</a>`
        }
        var charaList = "";
        if (this.characters !== undefined) {
            this.characters.forEach(character => {
                charaList += `<li>${selectCharacter.name(character)}</li>`;
            });
            charaList = `<ul>${charaList}</ul>`;
        }
        if (this._edition) {
            return `
            <div class="element">
                <div class="content">
                    <form id="form-${this.id}" method="POST" action="/api/scene/${this.book_id}/${this.id}">
                        <input type="hidden" name="id" value="${this.id}">
                        <input type="hidden" name="book_id" value="${this.book_id}">
                        <div>
                            <p>Name : </p><input type="text" value="${this.name}" name="name" class="custom-input w100">
                        </div>
                        <div class="mt2">
                            <label>Description: </label><textarea name="description" class="custom-textarea">${this.description}</textarea>
                            ${charaList}
                            <label>Add a character to the scene: </label>${selectCharacter.code("character-"+this.id)}
                            <button id="add-${this.id}" class="custom-button">Add</button>
                        </div>
                        <input type="hidden" name="characters" value="${this.characters}">
                        <input type="hidden" name="lines" value="${this.lines}">
                        <button id="save-${this.id}" type="submit" class="custom-button ml2 mt2">Save</button>
                        <button id="cancel-${this.id}" type="button" class="custom-button mt2">Cancel</button>
                    </form>
                    <div class="buttons buttons-center">
                        <button id="up-${this.id}"><i class="fa-solid fa-chevron-up"></i></button>
                        <button id="down-${this.id}"><i class="fa-solid fa-chevron-down"></i></button>
                    </div>
                    ${linkLines}
                </div>
                <div class="buttons buttons-right">
                    <button id="edit-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
                    <button  id="open-${this.id}"><i class="fa-regular fa-eye"></i></button>
                </div>
            </div>
            `
        }
        if (this.id === 0) {
            return `
                <div class="element">
                    <div class="content">
                        <div>
                            <p>Create a new scene</p>
                        </div>
                    </div>
                    <div class="buttons buttons-right">
                        <button id="edit-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
                        <button  id="open-${this.id}"><i class="fa-regular fa-eye"></i></button>
                    </div>
                </div>
                `
        }
        if (this._displayed) {
                return `
                    <div class="element">
                        <div class="content">
                            <div>
                                <p>Name : ${this.name}</p>
                            </div>
                            <div>
                                <p>Description : ${this.description}</p>
                                <p>Lines : ${this.numLines}</p>
                            </div>
                        </div>
                        <div class="buttons buttons-right">
                            <button id="edit-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
                            <button  id="open-${this.id}"><i class="fa-regular fa-eye"></i></button>
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
                    <button  id="open-${this.id}"><i class="fa-regular fa-eye"></i></button>
                </div>
            </div>
            `
    }
}

export { SceneComponent }