import { Component, EventBus } from "./framework.js" 

class DWComponent extends Component {
    constructor(id, parameters, callbacks) {
        super(id, parameters, callbacks);
        this._displayed = false;
        this._edition = false;
        this._edited = false;
        this._uri = "";
        this._init();
        EventBus.register("open-click", this._handleOpen.bind(this));
        EventBus.register("edit-click", this._handleEdit.bind(this));
        EventBus.register("save-click", this._handleSaved.bind(this));
        EventBus.register("cancel-click", this._handleCanceled.bind(this));
    }

    _handleOpen(event) {
        if (event.id === "open-"+this.id) {
            this._displayed = !this._displayed;
            if (this._displayed === false) {
                this._edition = false;
            }
            if (this.id === 0) {
                this._edition = this._displayed;
            }
            this.render();
        }
    }

    _handleEdit(event) {
        if (event.id === "edit-"+this.id) {
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
    }

    _handleSaved(event) {
        if (event.id === "save-"+this.id) {
            var form = document.getElementById("form-"+this.id);
            var formData = new FormData(form);

            var jsonData = {};
            formData.forEach(function(value, key) {
                if (key === "id") {
                    jsonData[key] = parseInt(value);
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
                    }
                    self._edition = false;
                    self._edited = false;
                    self._displayed = false;
                    self.render();
                    EventBus.dispatch("refresh");
                    return response.json();
                } else {
                    alert("Save failed.");
                }
            })
            .catch(error => {
                alert("An error occurred. Please try again."+error);
            });
        }
    }

    _handleCanceled(event) {
        if (event.id === "cancel-"+this.id) {
            this._edition = false;
            this._edited = false;
            this._displayed = false;
            this.render();
        }
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
            if (key === "id") {
                self[key] = parseInt(value);
            } else {
                self[key] = value;
            }
        });
    }
}

export { DWComponent }