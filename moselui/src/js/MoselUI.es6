import $ from "jquery";

window.app = new class {

    constructor() {
        this.pages = {
            dashboard: "pages/dashboard.html",
            page2: "pages/page2.html"
        };

        $(document).ready(() => {
            this.loadContent('dashboard');
        });
    }

    loadContent(page) {
        $('#content').load(this.pages[page]);
    }

    logIn() {
        console.log("login");
    }
}();