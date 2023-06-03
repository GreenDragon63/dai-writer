import { CharacterComponent } from "./character-component.js"
import { EventBus } from "./framework.js";

document.getElementById("upload-button").addEventListener("click", function(event) {
    event.preventDefault();

    const fileInput = document.getElementById("upload-card");
    if (fileInput.files.length === 0) {
        alert("Please select a file.");
        return;
    }

    const formData = new FormData();
    formData.append("file", fileInput.files[0]);

    fetch("/api/upload", {
        method: "POST",
        body: formData
    })
    .then(response => {
        if (response.ok) {
            document.getElementById("upload-card").value = "";
            document.getElementById("fileName").textContent = "Choose File";
            fetchLast();
        } else {
            alert("Upload failed. Please choose a character card.");
        }
    })
    .catch(error => {
        alert("An error occurred. Please try again.");
    });
});

function handleEye(event) {
    EventBus.dispatch("eye-click", {id: this.id});
}

function createCharacter(character) {
    let id = "char-" + character.id;
    let eyeId = "eye-" + character.id;
    let callbacks = {
        [eyeId]:
        {"click":handleEye}
    }
    let characterComponent = new CharacterComponent(id, character, callbacks);
    characterComponent.prependToDom("container");
}

function fetchLast() {
    fetch("/api" + window.location.pathname + "0")
    .then(response => response.json())
    .then(character => {
        createCharacter(character);
    });
}

function fetchAll() {
    fetch("/api" + window.location.pathname)
    .then(response => response.json())
    .then(data => {
        if (data === null) {
            return
        }
        data.forEach(character => {
            createCharacter(character);
        });
    });
}

fetchAll();
