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

function save(event) {
    event.preventDefault();
    var id = event.target.id.split("-")[1];
    var form = document.getElementById("edit-"+id);
    var formData = new FormData(form);

    var jsonData = {};
    formData.forEach(function(value, key) {
        if (key === "id") {
            jsonData[key] = parseInt(value);
        } else {
            jsonData[key] = value;
        }
    });

    fetch('/api/character/' + id, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(jsonData)
    })
    .then(function(response) {
        if (response.ok) {
            EventBus.dispatch("saved", {id: id});
            return response.json();
        } else {
            alert("Save failed.");
        }
    })
    .catch(error => {
        alert("An error occurred. Please try again.");
    });
}


function handleopen(event) {
    EventBus.dispatch("open-click", {id: this.id});
}

function handlePen(event) {
    EventBus.dispatch("pen-click", {id: this.id});
}

function cancel(event) {
    event.preventDefault();
    var id = event.target.id.split("-")[1];
    EventBus.dispatch("canceled", {id: id});
}

function createCharacter(character) {
    let id = "char-" + character.id;
    let openId = "open-" + character.id;
    let penId = "pen-" + character.id;
    let saveId = "save-" + character.id;
    let cancelId = "cancel-" + character.id;
    let callbacks = {
        [openId]:
        {"click":handleopen},
        [penId]:
        {"click":handlePen},
        [saveId]:
        {"click":save},
        [cancelId]:
        {"click":cancel}
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

document.getElementById("upload-button").addEventListener("click", upload);
fetchAll();
