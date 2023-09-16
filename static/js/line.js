import prefix from "../../config.js";;
import { LineComponent } from "./line-component.js"
import { EventBus } from "./framework.js";

function createLine(line, node) {
    let id = "line-" + line.id;
    let lineComponent = new LineComponent(id, line);
    lineComponent.appendToDom(node);
}

function fetchLast() {
    fetch(prefix + "/api" + window.location.pathname.replace(prefix,"") + "/0")
    .then(response => response.json())
    .then(line => {
        createLine(line, "container");
        EventBus.dispatch("refresh-order",{"cause":"fetchLast"});
    });
}

function fetchAll() {
    let decoded = window.location.pathname.replace(prefix,"").split('/');
    if (decoded.length < 4) {
        return
    }
    let book = decoded[2];
    let scene = decoded[3];
    let order;
    fetch(prefix + "/api/scene/"+book+"/"+scene)
    .then(response => response.json())
    .then(scene => {
        order = scene.lines;
        fetch(prefix + "/api" + window.location.pathname.replace(prefix,""))
        .then(response => response.json())
        .then(data => {
            if (data === null) {
                return
            }
            order.forEach(id => {
                createLine(data[id-1], "container");
            });
        });
    });
}

function addLine() {
    let decoded = window.location.pathname.replace(prefix,"").split('/');
    if (decoded.length < 4) {
        return
    }
    let book = decoded[2];
    let scene = decoded[3];
    let line = {
        "id": 0,
        "book_id": book,
        "scene_id": scene,
        "character_id": 0,
        "content": [""],
        "current": 0,
        "token": 0,
    }
    createLine(line, "add-new");
}

function addBreadcrumb() {
    const breadcrumb = document.getElementById("breadcrumb");
    let decoded = window.location.pathname.replace(prefix,"").split('/');
    if (decoded.length < 4) {
        return
    }
    let book = decoded[2];
    breadcrumb.innerHTML = '<a href="'+prefix+'/">Home</a>/<a href="'+prefix+'/book">Edit books</a>/<a href="'+prefix+'/scene/'+book+'">Edit scenes</a>';
}

EventBus.register("refresh", fetchLast);
addBreadcrumb();
addLine();
fetchAll();


