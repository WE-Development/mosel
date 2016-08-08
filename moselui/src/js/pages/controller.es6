import $ from "jquery";

export class Controller {

    constructor(container, view) {
        this.container = container;
        this.view = view;
    }

    init() {
        console.debug("Init", this)
    }

    load() {
        this.container.ready(() => {
            console.debug(this.container, this.view);

            if (this.container.controller instanceof Controller) {
                this.container.controller.destroy();
            }

            this.container.controller = this;
            this.container.load(this.view);
            $(this.container).ready(() => this.init());
        });
    }

    destroy() {
        console.debug("Destroy", this)
    }

    getChild(selector) {
        return this.container.find(selector);
    }
}
