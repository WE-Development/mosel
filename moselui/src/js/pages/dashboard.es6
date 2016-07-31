import {Controller} from "./page.es6";

export class Dashboard extends Controller {

    constructor(container) {
        super(container, "pages/dashboard.view");
    }

    init() {
        console.log("Dashboard init")
    }

    destroy() {
        console.log("Dashboard destroy")
    }

}