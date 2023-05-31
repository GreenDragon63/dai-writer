import { LineComponent } from "./line-component.js"

fetch("/api" + window.location.pathname)
    .then(response => response.json())
    .then(data => {
        var id = 0;
        data.forEach(line => {
            var lineComponent = new LineComponent("line" + id, line);
            lineComponent.appendToDom("container");
            id++;
        });
    });
