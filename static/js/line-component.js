import prefix from "../../config.js";;
import { DWMovableComponent } from "./dw-movable-component.js";
import { EventBus } from "./framework.js";
import selectCharacter from "./select-character.js";

class LineComponent extends DWMovableComponent {
    constructor(id, parameters, callbacks) {
        super(id, parameters, callbacks);
        this._uri = prefix + '/api/line/' + this.book_id + '/' + this.scene_id + '/';
        if ((this.lines === undefined) || (this.lines === null)) {
            this.lines = [];
        }
        let left = "left-" + parameters.id;
        let right = "right-" + parameters.id;
        let generate = "generate-" + parameters.id;
        this._addCallbacks({
            [left]: {"click":this._handleLeft.bind(this)},
            [right]: {"click":this._handleRight.bind(this)},
            [generate]: {"click":this._handleGenerate.bind(this)}
        });
        this._init();
        EventBus.register("refresh-order", this._saveOrder.bind(this));
        EventBus.register("chara-list", this.render.bind(this));
    }

    _save() {
        let line = {
            "id": this.id,
            "book_id": this.book_id,
            "scene_id": this.scene_id,
            "displayed": this._displayed,
            "character_id": this.character_id,
            "content": this.content,
            "current": this.current,
            "tokens": this.tokens,
        };
        fetch(this._uri + this.id, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(line)
        });
    }

    _handleGenerate(event) {
        event.preventDefault();
        var modal = document.getElementById("modal");
        modal.style.display = "flex";
        self = this;
        if (this.id === 0) {
            var charaSelect = document.getElementById("character-"+this.id);
            var character = charaSelect.value;
        } else {
            var character = this.character_id
        }
        fetch(prefix + "/api/generate/" + this.book_id + "/" + this.scene_id + "/" + + character + "/" + this.id, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({input: document.getElementById("system-input-"+this.id).value})
        })
        .then(response => response.json())
        .then(text => {
            if ((self.content.length == 1) && (self.content[0] == "")) {
                self.content[0] = text;
            } else {
                self.content.push(text);
                self.current = self.content.length - 1;
            }
            self.render();
            modal.style.display = "none";
        });
    }

    _handleLeft(event) {
        event.preventDefault();
        if (this.current > 0) {
            this.current--;
            this.current_content = this.content[this.current];
            this.render();
        }
    }

    _handleRight(event) {
        event.preventDefault();
        if (this.current < this.content.length - 1) {
            this.current++;
            this.current_content = this.content[this.current];
            this.render();
        }
    }

    _handleOpen(event) {
        super._handleOpen(event);
        if (this.id === 0) {
            return
        }
        this._save();
    }

    _saveOrder(event) {
        if (typeof event !== "undefined") {
            if (this.id !==0) {
                return;
            }
        }
        let lineList = [];
        let container = document.getElementById("container");
        let lines = container.children;
        for (var i = 0; i < lines.length; i++) {
            var line = lines[i];
            let decoded = line.id.split("-");
            if (decoded.length !== 2) {
                throw "Line id, wrong format"
            }
            lineList.push(parseInt(decoded[1]));
        }
        fetch(prefix + "/api/scene/"+ this.book_id + "/" + this.scene_id)
        .then(response => response.json())
        .then(scene => {
            scene.lines = lineList;
            fetch(prefix + "/api/scene/"+ this.book_id + "/" + this.scene_id, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(scene)
            });
        });
    }

    _markdownParser(text) {
        var toHTML = text
            .replace(/\n/gim, '<br>')
            .replace(/^### (.*$)/gim, '<h3>$1</h3>')
            .replace(/^## (.*$)/gim, '<h2>$1</h2>')
            .replace(/^# (.*$)/gim, '<h1>$1</h1>')
            .replace(/\*\*(.*)\*\*/gim, '<b>$1</b>')
            .replace(/\*(.*)\*/gim, '<i>$1</i>');
        return toHTML.trim();
    }

    _template() {
        if (this.id === 0) {
            var arrows = "";
            var character = selectCharacter.scene("character-"+this.id, this.book_id, this.scene_id)
        } else {
            var arrows = `
            <div class="buttons buttons-center">
                <button id="up-${this.id}"><i class="fa-solid fa-chevron-up"></i></button>
                <button id="down-${this.id}"><i class="fa-solid fa-chevron-down"></i></button>
            </div>`;
            var character = `<input type="hidden" name="character_id" value="${this.character_id}">`+selectCharacter.name(this.character_id)
        }
        if (this.id === 0) {
            return `
                <div class="element">
                    <div class="content">
                        <div>
                            <form id="form-${this.id}" method="POST" action="${prefix}/api/scene/${this.book_id}/${this.id}">
                                <input type="hidden" name="id" value="${this.id}">
                                <input type="hidden" name="book_id" value="${this.book_id}">
                                <input type="hidden" name="scene_id" value="${this.scene_id}">
                                <input type="hidden" name="displayed" value="${this._displayed}">
                                <input type="hidden" name="current" value="${this.current}">
                                <input type="hidden" name="content" value="${btoa(encodeURIComponent(JSON.stringify(this.content)))}">
                                <input type="hidden" name="current_content" value="">
                                <p>Character : ${character}
                                <button id="save-${this.id}" type="submit" class="custom-button ml2 mt2">Create a new line of dialogue</button>
                                </p>
                            </form>
                        </div>
                    </div>
                </div>
                `
        }
        if (this._edition) {
            return `
            <div class="element">
                <div class="image-container">
                    <img src="${prefix}/api/avatar/${this.character_id}" alt="Image">
                </div>
                <div class="content">
                    <form id="form-${this.id}" method="POST" action="${prefix}/api/scene/${this.book_id}/${this.id}">
                        <input type="hidden" name="id" value="${this.id}">
                        <input type="hidden" name="book_id" value="${this.book_id}">
                        <input type="hidden" name="scene_id" value="${this.scene_id}">
                        <input type="hidden" name="displayed" value="${this._displayed}">
                        <input type="hidden" name="current" value="${this.current}">
                        <input type="hidden" name="content" value="${btoa(encodeURIComponent(JSON.stringify(this.content)))}">
                        <div>
                            <p>Character : ${character}</p>
                        </div>
                        <div class="mt2">
                            <label>Content: </label>
                            <div class="buttons">
                                <button id="left-${this.id}"><i class="fa-regular fa-square-caret-left"></i></button>
                                &nbsp;&nbsp;${this.current+1}/${this.content.length}
                                <button id="right-${this.id}"><i class="fa-regular fa-square-caret-right"></i></button>
                            </div>
                            <textarea name="current_content" class="custom-textarea">${this.content[this.current]}</textarea>
                        </div>
                        <button id="save-${this.id}" type="submit" class="custom-button ml2 mt2">Save</button>
                        <button id="cancel-${this.id}" type="button" class="custom-button mt2">Cancel</button>
                    </form>
                    ${arrows}
                </div>
                <div class="buttons buttons-right">
                    <button id="edit-${this.id}"><i class="fa-regular fa-rectangle-xmark"></i></button>
                    <button  id="open-${this.id}"><i class="fa-regular fa-eye"></i></button>
                </div>
            </div>
            `
        }
        if (this._displayed) {
            if (document.getElementById("system-input-"+this.id) !== null) {
                var systemInput = document.getElementById("system-input-"+this.id).value;
            } else {
                var systemInput = "";
            }
            return `
                <div class="element">
                    <div class="image-container">
                        <img src="${prefix}/api/avatar/${this.character_id}" alt="Image">
                    </div>
                    <div class="content">
                        <div>
                            <p><b>${selectCharacter.name(this.character_id)}</b></p>
                            <p><input type="text" class="w100 ml0 custom-input" id="system-input-${this.id}" value="${systemInput}" placeholder="Optional system input for generation"></p>
                        </div>
                        <div>
                            <p>${this._markdownParser(this.content[this.current])}</p>
                        </div>
                        <div class="buttons buttons-center">
                            <button id="up-${this.id}"><i class="fa-solid fa-chevron-up"></i></button>
                            <button id="down-${this.id}"><i class="fa-solid fa-chevron-down"></i></button>
                        </div>
                    </div>
                    <div class="buttons buttons-right">
                        <button id="generate-${this.id}"><i class="fa-solid fa-rotate"></i></button>
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
                        <p><b>${selectCharacter.name(this.character_id)}: </b> ${this.content[this.current].substring(0,50)}...</p>
                    </div>
                </div>
                <div class="buttons buttons-right">
                    <button id="edit-${this.id}"><i class="fa-solid fa-pen-to-square"></i></button>
                    <button  id="open-${this.id}"><i class="fa-regular fa-eye-slash"></i></button>
                </div>
            </div>
            `
    }
}

export { LineComponent }