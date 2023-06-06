import { BookComponent } from "./book-component.js"
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

function createBook(book, node) {
    let id = "book-" + book.id;
    let openId = "open-" + book.id;
    let editId = "edit-" + book.id;
    let saveId = "save-" + book.id;
    let cancelId = "cancel-" + book.id;
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
    let bookComponent = new BookComponent(id, book, callbacks);
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


