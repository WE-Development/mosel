var fs = require("fs");
var browserify = require("browserify");
var bablify = require("babelify");
var cp = require("copyfiles");
var rm = require("del");
var mkdir = require("mkdirp");

//cleanup
rm.sync(['dist']);
mkdir.sync('dist');

cp([
        './src/**/*.html',

        './dist'
    ], true,
    function (err, files) {
        if (typeof err != 'undefined') console.error(err);
        console.log('Copied ' + files);
    });

browserify('src/MoselUI.es6')
    .transform(bablify, {presets: ["es2015"]})
    .bundle()
    .pipe(fs.createWriteStream('dist/moselui.min.js'));