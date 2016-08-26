#!/usr/bin/env bash

#install moseld
pushd $GOPATH

go get github.com/WE-Development/mosel/moselnoded
sudo cp bin/moselnoded /usr/local/bin/

popd
#create config
pushd $(dirname $0)

sudo mkdir /etc/mosel
sudo cp moselnoded.conf /etc/mosel/

popd