import { EventBus } from "./framework.js";

class SelectCharacter {
    constructor() {
        this._characters = [];
        this._sceneCharacters = [];
        self = this;
        fetch("/api/character/infos/")
        .then(response => response.json())
        .then(data => {
            if (data === null) {
                return
            }
            data.forEach(character => {
                self._characters.push(character);
                self._options += `<option value="${character.id}">${character.name}</option>`
            });
        });
    }

    all(id) {
        return `<select id="${id}" class="custom-select">${this._options}</select>`
    }

    scene(id, book_id, scene_id) {
        if (this._sceneCharacters.length === 0) {
            self = this;
            fetch("/api/scene/"+book_id+"/"+scene_id)
            .then(response => response.json())
            .then(scene => {
                if (scene === null) {
                    return
                }
                self._sceneCharacters = scene.characters;
                self._sceneOptions = "";
                self._sceneCharacters.forEach(character => {
                    self._sceneOptions += `<option value="${character}">${self.name(character)}</option>`
                });
                EventBus.dispatch("chara-list");
                return `<select id="${id}" class="custom-select" name="character_id">${this._sceneOptions}</select>`
            });
        }
        return `<select id="${id}" class="custom-select" name="character_id">${this._sceneOptions}</select>`
    }

    name(id) {
        var result = "";
        this._characters.forEach(character => {
            if (character.id === id) {
                result = character.name;
            }
        });
        return result
    }

    infos(id) {
        var result = "";
        this._characters.forEach(character => {
            if (character.id === id) {
                result = "Description: " + self.escapeHtml(character.description) + "\nScenario: " + character.scenario;
            }
        });
        return result
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML.replaceAll('"', '”');
    }

}

const selectCharacter = new SelectCharacter();

export default selectCharacter;