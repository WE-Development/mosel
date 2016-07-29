var fs = require("fs");
var browserify = require("browserify");
var bablify = require("babelify");

browserify('src/MoselUI.es6')
    .transform(bablify, {presets: ["es2015"]})
    .bundle()
    .pipe(fs.createWriteStream('dist/moselui.min.js'));