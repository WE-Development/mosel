import $ from "jquery";
import {Page} from "./pages/page.es6";
import {Dashboard} from "./pages/dashboard.es6";

window.app = new class {

    constructor() {
        this.pages = {
            dashboard: new Dashboard(),
            page2: new Page("pages/page2.html")
        };

        $(document).ready(() => {
            this.loadContent('dashboard');
        });
    }

    loadContent(pageName) {

        if (typeof this.currentPage != 'undefined') {
            this.currentPage.destroy();
        }

        if (!(pageName in this.pages)) {
            return;
        }

        this.currentPage = this.pages[pageName];
        var content = $('#content');

        content.load(this.currentPage.html);
        content.ready(() => this.currentPage.init());
    }

    logIn() {
        console.log("login");
    }
}();