import { LineComponent } from "./line-component.js"
import { EventBus } from "./framework.js";

function createLine(line, node) {
    let id = "line-" + line.id;
    let lineComponent = new LineComponent(id, line);
    lineComponent.appendToDom(node);
}

function fetchLast() {
    fetch("/api" + window.location.pathname + "/0")
    .then(response => response.json())
    .then(line => {
        createLine(line, "container");
        EventBus.dispatch("refresh-order",{"cause":"fetchLast"});
    });
}

function fetchAll() {
    let decoded = window.location.pathname.split('/');
    if (decoded.length < 4) {
        return
    }
    let book = decoded[2];
    let scene = decoded[3];
    let order;
    fetch("/api/scene/"+book+"/"+scene)
    .then(response => response.json())
    .then(scene => {
        order = scene.lines;
    });
    fetch("/api" + window.location.pathname)
    .then(response => response.json())
    .then(data => {
        if (data === null) {
            return
        }
        order.forEach(id => {
            createLine(data[id-1], "container");
        });
    });
}

function addLine() {
    let decoded = window.location.pathname.split('/');
    if (decoded.length < 4) {
        return
    }
    let book = decoded[2];
    let scene = decoded[3];
    let line = {
        "id": 0,
        "book_id": book,
        "scene_id": scene,
        "character": null,
        "content": "",
        "token": 0,
    }
    createLine(line, "add-new");
}

EventBus.register("refresh", fetchLast);
addLine();
fetchAll();


