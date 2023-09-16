import prefix from "../../config.js";
import { DWComponent } from "./dw-component.js";

class BookComponent extends DWComponent {
    constructor(id, parameters, callbacks) {
        super(id, parameters, callbacks);
        this._uri = prefix + '/api/book/';
        if ((this.scenes === undefined) || (this.scenes === null)) {
            this.scenes = [];
        }
        this.numScenes = this.scenes.length;
    }

    _template() {
        if (this._edition) {
            return `
            <div class="element">
                <div class="content">
                    <form id="form-${this.id}" method="POST" action="${prefix}/api/book/${this.id}">
                        <input type="hidden" name="id" value="${this.id}">
                        <div>
                            <p>Title : </p><input type="text" value="${this.name}" name="name" class="custom-input w100">
                        </div>
                        <div class="mt2">
                            <label>Description: </label><textarea name="description" class="custom-textarea">${this.description}</textarea>
                        </div>
                        <input type="hidden" name="scenes" value="${this.scenes}">
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
                            <p>Create a new book</p>
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
                        <div class="content">
                            <div>
                                <p>Title : ${this.name}</p>
                            </div>
                            <div>
                                <p>Description : ${this.description}</p>
                                <p>Scenes : ${this.numScenes}</p>
                            </div>
                            <a href="${prefix}/scene/${this.id}/" class="custom-button button-link ml2 mt2">Edit scenes</a>
                            <a href="${prefix}/api/export/${this.id}/" class="custom-button button-link ml2 mt2">Export book</a>
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
                        <p>Title : ${this.name}</p>
                    </div>
                    <a href="${prefix}/scene/${this.id}/" class="custom-button button-link ml2 mt2">Edit scenes</a>
                    <a href="${prefix}/api/export/${this.id}/" class="custom-button button-link ml2 mt2">Export book</a>
                </div>
                <div class="buttons buttons-right">
                    <button id="edit-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
                    <button id="open-${this.id}"><i class="fa-solid fa-folder-closed"></i></button>
                </div>
            </div>
            `
    }
}

export { BookComponent }