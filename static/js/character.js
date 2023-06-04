import { CharacterComponent } from "./character-component.js"
import { EventBus } from "./framework.js";

function upload(event) {
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
}

function handleopen(event) {
    EventBus.dispatch("open-click", {id: this.id});
}

function handleEdit(event) {
    EventBus.dispatch("edit-click", {id: this.id});
}

function save(event) {
    event.preventDefault();
    EventBus.dispatch("save-click", {id: this.id});
}

function cancel(event) {
    event.preventDefault();
    EventBus.dispatch("cancel-click", {id: this.id});
}

function createCharacter(character, node) {
    let id = "char-" + character.id;
    let openId = "open-" + character.id;
    let editId = "edit-" + character.id;
    let saveId = "save-" + character.id;
    let cancelId = "cancel-" + character.id;
    let callbacks = {
        [openId]:
        {"click":handleopen},
        [editId]:
        {"click":handleEdit},
        [saveId]:
        {"click":save},
        [cancelId]:
        {"click":cancel}
    }
    let characterComponent = new CharacterComponent(id, character, callbacks);
    characterComponent.prependToDom(node);
}

function fetchLast() {
    fetch("/api" + window.location.pathname + "0")
    .then(response => response.json())
    .then(character => {
        createCharacter(character, "container");
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
            createCharacter(character, "container");
        });
    });
}

function addCharacter() {
    let character = {
        "id": 0,
        "name": "",
        "description": "",
        "personality": "",
        "scenario": "",
        "mes_example": "",
        "first_mes": ""
    }
    createCharacter(character, "add-new");
}

document.getElementById("upload-button").addEventListener("click", upload);
EventBus.register("refresh", fetchLast);
addCharacter();
fetchAll();
