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
                this.container.removeAttr('data-container');
            }

            this.container.controller = this;
            this.container.load(this.view, () => this.init());
            this.container.attr('data-container', '');
        });
    }

    destroy() {
        console.debug("Destroy", this)
    }

    getChild(selector) {
        return this.getChildren(selector).first();
    }

    getChildren(selector) {
        var c = this.container;
        return c.find(selector)
            .not(c.find('*[data-container=""]').find(selector));
    }
}
