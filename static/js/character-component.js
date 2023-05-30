import { Component } from "./framework.js" 

class CharacterComponent extends Component {
    _template() {
        return `
            <div class="element">
                <div class="image-container">
                    <img src="/static/img/placeholder.svg" alt="Image">
                </div>
                <div class="content">
                    <div>
                        <p>${this.name}</p>
                    </div>
                    <div>
                        <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer dictum ante nec magna congue, non mattis sem volutpat. Suspendisse potenti. Duis auctor consectetur mauris eget consequat. Vivamus faucibus consequat mi, vel ultrices neque efficitur id. Mauris nec feugiat sapien. Etiam non volutpat velit. Fusce nec turpis sit amet mauris varius mattis ut et nisi. Aenean maximus lorem ut nunc finibus, eget tristique libero tincidunt. Quisque sed faucibus orci, non rhoncus tortor.<p>
                        <p>Sed id neque nibh. In euismod bibendum dolor, a aliquet dui fermentum a. Proin aliquam neque vitae turpis semper, a consequat sapien scelerisque. Nulla facilisi. Mauris nec mi leo. Proin lobortis, justo nec tincidunt placerat, felis sem semper ipsum, a efficitur purus felis at ante. Maecenas at mauris eget tortor tincidunt suscipit vel nec ante. Etiam auctor, tortor et dictum fermentum, velit nulla scelerisque lacus, sed lacinia mauris purus id quam. Curabitur varius vestibulum dui id iaculis. Integer sed elementum neque.</p>
                    </div>
                    <div class="buttons buttons-center">
                        <button><i class="fa-solid fa-chevron-up"></i></button>
                        <button><i class="fa-solid fa-chevron-down"></i></button>
                    </div>
                </div>
                <div class="buttons buttons-right">
                    <button><i class="fa-solid fa-rotate"></i></button>
                    <button><i class="fa-solid fa-pen-to-square"></i></button>
                    <button><i class="fa-regular fa-eye"></i></button>
                </div>
            </div>
            `
    }
}

export { CharacterComponent }