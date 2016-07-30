import $ from "jquery";

window.app = new class {

    constructor() {
        console.log("Init mosel ui");
        console.log($);
        //console.log(bootstrap);
    }

    logIn() {
        console.log("login");
    }
}();