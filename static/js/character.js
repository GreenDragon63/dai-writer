import { CharacterComponent } from "./character-component.js"

fetch("/api/character/")
    .then(response => response.json())
    .then(data => {
        var id = 0;
        data.forEach(character => {
            var characterComponent = new CharacterComponent("char" + id, character);
            characterComponent.appendToDom("container");
            id++;
        });
    });

