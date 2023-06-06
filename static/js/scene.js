import { SceneComponent } from "./scene-component.js"
import { EventBus } from "./framework.js";

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

function up(event) {
    EventBus.dispatch("up-click", {id: this.id});
}

function down(event) {
    EventBus.dispatch("down-click", {id: this.id});
}

function createScene(scene, node) {
    let id = "scene-" + scene.id;
    let openId = "open-" + scene.id;
    let editId = "edit-" + scene.id;
    let saveId = "save-" + scene.id;
    let cancelId = "cancel-" + scene.id;
    let upId = "up-" + scene.id;
    let downId = "down-" + scene.id;
    let callbacks = {
        [openId]:
        {"click":handleopen},
        [editId]:
        {"click":handleEdit},
        [saveId]:
        {"click":save},
        [cancelId]:
        {"click":cancel},
        [upId]:
        {"click":up},
        [downId]:
        {"click":down}
    }
    let sceneComponent = new SceneComponent(id, scene, callbacks);
    sceneComponent.prependToDom(node);
}

function fetchLast() {
    fetch("/api" + window.location.pathname + "0")
    .then(response => response.json())
    .then(scene => {
        createScene(scene, "container");
    });
}

function fetchAll() {
    fetch("/api" + window.location.pathname)
    .then(response => response.json())
    .then(data => {
        if (data === null) {
            return
        }
        data.forEach(scene => {
            createScene(scene, "container");
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
        "book": book,
        "name": "",
        "description": "",
    }
    createScene(scene, "add-new");
}

EventBus.register("refresh", fetchLast);
addScene();
fetchAll();


