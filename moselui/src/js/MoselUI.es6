import $ from "jquery";

class Page {

    constructor(html) {
        this.html = html;
    }

    init() {
    }

    destroy() {
    }
}

class Dashboard extends Page {

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