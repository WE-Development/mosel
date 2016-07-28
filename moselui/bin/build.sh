#!/bin/bash

npm install

rm -r distribution
mkdir distribution

node node_modules/babel-cli/bin/babel.js src --out-file distribution/moselui.min.js

node node_modules/copyfiles/copyfiles -u 1 'src/**/*.html' distribution