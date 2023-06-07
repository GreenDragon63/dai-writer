import { Component, EventBus } from "./framework.js" 

class DWComponent extends Component {
    constructor(id, parameters) {
        super(id, parameters);
        let openId = "open-" + parameters.id;
        let editId = "edit-" + parameters.id;
        let saveId = "save-" + parameters.id;
        let cancelId = "cancel-" + parameters.id;
        this._addCallbacks({
            [openId]: {"click":this._handleOpen.bind(this)},
            [editId]: {"click":this._handleEdit.bind(this)},
            [saveId]: {"click":this._handleSaved.bind(this)},
            [cancelId]: {"click":this._handleCanceled.bind(this)}
        });
        this._displayed = false;
        this._edition = false;
        this._edited = false;
        this._uri = "";
        if (this.displayed) {
            this._displayed = this.displayed;
        }
        this._init();
    }

    _handleOpen(event) {
        this._displayed = !this._displayed;
        if (this._displayed === false) {
            this._edition = false;
        }
        if (this.id === 0) {
            this._edition = this._displayed;
        }
        this.render();
    }

    _handleEdit(event) {
        this._edition = !this._edition;
        if (this._edition === true) {
            this._displayed = true;
        }
        if (this.id === 0) {
            this._displayed = this._edition;
        }
        this.render();
        if (this._edition === true) {
            const formElement = document.getElementById("form-"+this.id);
            const formInputs = formElement.querySelectorAll('input, select, textarea');

            formInputs.forEach(input => {
                input.addEventListener('input', this._handleInput.bind(this));
            });
        }
    }

    _handleSaved(event) {
        event.preventDefault();
        var form = document.getElementById("form-"+this.id);
        var formData = new FormData(form);

        var jsonData = {};
        formData.forEach(function(value, key) {
            if (key.endsWith("id")) {
                jsonData[key] = parseInt(value);
            } else if (key === "scenes") {
                jsonData[key] = value.split(",").map(function(num) {
                    return parseInt(num);
                });
            } else if (key === "lines") {
                jsonData[key] = value.split(",").map(function(num) {
                    return parseInt(num);
                });
            } else if (key === "characters") {
                jsonData[key] = value.split(",").map(function(num) {
                    return parseInt(num);
                });
            } else if (key === "displayed") {
                jsonData[key] = !!value;
            } else {
                jsonData[key] = value;
            }
        });
        self = this;
        fetch(this._uri + this.id, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(jsonData)
        })
        .then(function(response) {
            if (response.ok) {
                if (self.id !== 0) {
                    self._refresh();
                } else {
                    EventBus.dispatch("refresh");
                }
                self._edition = false;
                self._edited = false;
                self._displayed = false;
                self.render();
                return response.json();
            } else {
                alert("Save failed.");
            }
        })
        .catch(error => {
            alert("An error occurred. Please try again."+error);
        });
    }

    _handleCanceled(event) {
        event.preventDefault();
        this._edition = false;
        this._edited = false;
        this._displayed = false;
        this.render();
    }

    _handleInput(event) {
        if (this._edited === false) {
            this._edited = true;
            const openButton = document.getElementById("open-"+this.id);
            openButton.disabled = true;
            const editButton = document.getElementById("edit-"+this.id);
            editButton.disabled = true;
        }
    }

    _refresh() {
        var form = document.getElementById("form-"+this.id);
        var formData = new FormData(form);
        var self = this;
        formData.forEach(function(value, key) {
            if (key.endsWith("id")) {
                self[key] = parseInt(value);
            } else if (key === "scenes") {
                self[key] = value.split(",").map(function(num) {
                    return parseInt(num);
                });
            } else if (key === "lines") {
                self[key] = value.split(",").map(function(num) {
                    return parseInt(num);
                });
            } else if (key === "characters") {
                self[key] = value.split(",").map(function(num) {
                    return parseInt(num);
                });
            } else if (key === "displayed") {
                self[key] = !!value;
            } else {
                self[key] = value;
            }
        });
    }
}

export { DWComponent }