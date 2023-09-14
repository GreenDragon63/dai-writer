import prefix from "./config.js"
import { SceneComponent } from "./scene-component.js"
import { EventBus } from "./framework.js";

function createScene(scene, node) {
    let id = "scene-" + scene.id;
    let sceneComponent = new SceneComponent(id, scene);
    sceneComponent.appendToDom(node);
}

function fetchLast() {
    fetch(prefix + "/api" + window.location.pathname.replace(prefix,"") + "/0")
    .then(response => response.json())
    .then(scene => {
        createScene(scene, "container");
        EventBus.dispatch("refresh-order",{"cause":"fetchLast"});
    });
}

function fetchAll() {
    let decoded = window.location.pathname.replace(prefix,"").split('/');
    if (decoded.length < 3) {
        return
    }
    let book = decoded[2];
    let order;
    fetch(prefix + "/api/book/"+book)
    .then(response => response.json())
    .then(book => {
        order = book.scenes;
        fetch(prefix + "/api" + window.location.pathname.replace(prefix,""))
        .then(response => response.json())
        .then(data => {
            if (data === null) {
                return
            }
            order.forEach(id => {
                createScene(data[id-1], "container");
            });
        });
    });
}

function addScene() {
    let decoded = window.location.pathname.replace(prefix,"").split('/');
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

function addBreadcrumb() {
    const breadcrumb = document.getElementById("breadcrumb");
    breadcrumb.innerHTML = '<a href="'+prefix+'/">Home</a>/<a href="'+prefix+'/book">Edit books</a>';
}

EventBus.register("refresh", fetchLast);
addBreadcrumb();
addScene();
fetchAll();


