import $ from "jquery";

export class Controller {

    constructor(container, html) {
        this.container = container;
        this.html = html;
    }

    init() {
    }

    load() {
        console.debug(this.container, this.html);

        if (typeof this.container.controller != 'undefined') {
            this.container.controller.destroy();
        }

        this.container.controller = this;
        this.container.load(this.html);
        $(this.container).ready(() => this.init());
    }

    destroy() {
    }
}
