import { SceneComponent } from "./scene-component.js"
import { EventBus } from "./framework.js";

function createScene(scene, node) {
    let id = "scene-" + scene.id;
    let sceneComponent = new SceneComponent(id, scene);
    sceneComponent.prependToDom(node);
}

function fetchLast() {
    fetch("/api" + window.location.pathname + "/0")
    .then(response => response.json())
    .then(scene => {
        createScene(scene, "container");
    });
}

function fetchAll() {
    let decoded = window.location.pathname.split('/');
    if (decoded.length < 3) {
        return
    }
    let book = decoded[2];
    let order;
    fetch("/api/book/"+book)
    .then(response => response.json())
    .then(book => {
        order = book.scenes;
    });
    fetch("/api" + window.location.pathname)
    .then(response => response.json())
    .then(data => {
        if (data === null) {
            return
        }
        order.slice().reverse().forEach(id => {
            createScene(data[id-1], "container");
        });
    });
}

function addScene() {
    let decoded = window.location.pathname.split('/');
    if (decoded.length < 3) {
        return
    }
    let book = decoded[2];
    let scene = {
        "id": 0,
        "book_id": book,
        "name": "",
        "description": "",
    }
    createScene(scene, "add-new");
}

EventBus.register("refresh", fetchLast);
addScene();
fetchAll();


