import prefix from "./config.js"
import { BookComponent } from "./book-component.js"
import { EventBus } from "./framework.js";

function createBook(book, node) {
    let id = "book-" + book.id;
    let bookComponent = new BookComponent(id, book);
    bookComponent.prependToDom(node);
}

function fetchLast() {
    fetch(prefix + "/api" + window.location.pathname.replace(prefix,"") + "0")
    .then(response => response.json())
    .then(book => {
        createBook(book, "container");
    });
}

function fetchAll() {
    fetch(prefix + "/api" + window.location.pathname.replace(prefix,""))
    .then(response => response.json())
    .then(data => {
        if (data === null) {
            return
        }
        data.forEach(book => {
            createBook(book, "container");
        });
    });
}

function addBook() {
    let book = {
        "id": 0,
        "name": "",
        "description": "",
    }
    createBook(book, "add-new");
}

function addBreadcrumb() {
    const breadcrumb = document.getElementById("breadcrumb");
    breadcrumb.innerHTML = '<a href="'+prefix+'/">Home</a>';
}

EventBus.register("refresh", fetchLast);
addBreadcrumb();
addBook();
fetchAll();


