import {Page} from "./page.es6";

export class Dashboard extends Page {

    constructor() {
        super("pages/dashboard.html");
    }

    init() {
        console.log("Dashboard init")
    }

    destroy() {
        console.log("Dashboard destroy")
    }

}