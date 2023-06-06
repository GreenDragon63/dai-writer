import { BookComponent } from "./book-component.js"
import { EventBus } from "./framework.js";

function createBook(book, node) {
    let id = "book-" + book.id;
    let bookComponent = new BookComponent(id, book);
    bookComponent.prependToDom(node);
}

function fetchLast() {
    fetch("/api" + window.location.pathname + "0")
    .then(response => response.json())
    .then(book => {
        createBook(book, "container");
    });
}

function fetchAll() {
    fetch("/api" + window.location.pathname)
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

EventBus.register("refresh", fetchLast);
addBook();
fetchAll();


