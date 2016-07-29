/*
 * CONFIG
 */

var include = [
    'src/**/*.html',
    'src/**/*.css'
];
var dist = 'dist';
var es6Main = 'src/MoselUI.es6';

/*
 * CONFIG END
 */

var fs = require("fs");
var browserify = require("browserify");
var bablify = require("babelify");
var cp = require("copyfiles");
var rm = require("del");
var mkdir = require("mkdirp");

//cleanup
rm.sync([dist]);
mkdir.sync(dist);

include.push(dist);
cp(include, 1,
    function (err, files) {
        if (typeof err != 'undefined') console.error(err);
        //console.log('Copied ' + files);
    });

browserify(es6Main)
    .transform(bablify, {presets: ["es2015"]})
    .bundle()
    .pipe(fs.createWriteStream('dist/moselui.min.js'));