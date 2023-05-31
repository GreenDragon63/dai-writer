import { SceneComponent } from "./scene-component.js"




fetch("/api" + window.location.pathname)
    .then(response => response.json())
    .then(data => {
        var id = 0;
        data.forEach(scene => {
            var sceneComponent = new SceneComponent("scene" + id, scene);
            sceneComponent.appendToDom("container");
            id++;
        });
    });

