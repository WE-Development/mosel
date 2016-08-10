/*
 * CONFIG
 */

var dist = 'dist';
var es2015Target = 'dist/js/moselui.min.js';

var jsInclude = [
    'src/js/setup.js',
    'src/js/index.es6'
];

var resourceInclude = [
    {
        up: 1,
        paths: [
            'src/**/*.html',
            'src/**/*.css'
        ]
    },
    {
        up: 3,
        paths: [
            'node_modules/bootstrap/dist/**/*.min.css',
            'node_modules/bootstrap/dist/**/*.min.css.map',
            'node_modules/bootstrap/dist/**/*.eot',
            'node_modules/bootstrap/dist/**/*.svg',
            'node_modules/bootstrap/dist/**/*.ttf',
            'node_modules/bootstrap/dist/**/*.woff',
            'node_modules/bootstrap/dist/**/*.woff2'
        ]
    }
];

/*
 * CONFIG END
 */

var browserify = require("browserify");
var bablify = require("babelify");
var uglifyify = require("uglifyify");

var fs = require("fs");
var cp = require("copyfiles");
var rm = require("del");
var mkdir = require("mkdirp");
var filePath = require("filepath");

//cleanup
rm.sync([dist]);
mkdir.sync(dist);

//include resources
resourceInclude.forEach(function(inc) {
    var toCopy = inc.paths;
    toCopy.push(dist);
    cp(toCopy, inc.up,
        function (err, files) {
            if (typeof err != 'undefined') console.error(err);
            //console.log('Copied ' + files);
        });
});

//compile es2015
var targetPath = filePath.create(es2015Target);
mkdir.sync(targetPath.dir().toString());

var bundler = new browserify({debug: true});
jsInclude.forEach(function(inc) { 
	bundler.add(inc);
});
bundler
    .transform(bablify, {presets: ["es2015"]})
    .transform(uglifyify, {global: true})
    .bundle()
    .pipe(fs.createWriteStream(es2015Target));
