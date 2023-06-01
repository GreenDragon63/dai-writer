import { CharacterComponent } from "./character-component.js"

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
            alert("Upload successful!");
            document.getElementById("upload-card").value = "";
            document.getElementById("fileName").textContent = "Choose File";
        } else {
            alert("Upload failed. Please choose a character card.");
        }
    })
    .catch(error => {
        alert("An error occurred. Please try again.");
    });
});


fetch("/api" + window.location.pathname)
    .then(response => response.json())
    .then(data => {
        var id = 0;
        data.forEach(character => {
            var characterComponent = new CharacterComponent("char" + id, character);
            characterComponent.prependToDom("container");
            id++;
        });
    });

