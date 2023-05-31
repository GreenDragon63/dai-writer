import { BookComponent } from "./book-component.js"

fetch("/api" + window.location.pathname)
    .then(response => response.json())
    .then(data => {
        var id = 0;
        data.forEach(book => {
            var bookComponent = new BookComponent("book" + id, book);
            bookComponent.appendToDom("container");
            id++;
        });
    });

