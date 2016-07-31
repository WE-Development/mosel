import $ from "jquery";
import {MoselUI} from "./pages/moselui.es6";

$(document).ready(function () {
    window.app = new MoselUI($('#moselui'));
    window.app.load();
});