import $ from "jquery";
import {Controller} from "./pages/page.es6";
import {Dashboard} from "./pages/dashboard.es6";

class MoselUI extends Controller {

    constructor() {
        super($('#moselui'), 'pages/moselui.html');
    }

    init() {
        console.log('Init MoselUI');
        var content = $('#content');

        this.pages = {
            dashboard: new Dashboard(content),
            page2: new Controller(content, "pages/page2.html")
        };

        this.loadContent('dashboard');
    }

    loadContent(pageName) {
        if (pageName in this.pages) {
            this.pages[pageName].load();
        }
    }

    logIn() {
        console.log("login");
    }
}

$(document).ready(() => {
    window.app = new MoselUI();
    window.app.load();
});