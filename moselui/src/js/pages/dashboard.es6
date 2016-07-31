import {Controller} from "./controller.es6";

export class Dashboard extends Controller {

    constructor(container) {
        super(container, "pages/dashboard.html");
    }

    init() {
        console.log("Dashboard init")
    }

    destroy() {
        console.log("Dashboard destroy")
    }

}