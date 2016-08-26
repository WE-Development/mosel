#!/usr/bin/env bash

#install moseld
pushd $GOPATH

go get github.com/WE-Development/mosel/moseld
sudo cp bin/moseld /usr/local/bin/

popd
#create config
pushd $(dirname $0)

sudo mkdir /etc/mosel
sudo cp moseld.conf /etc/mosel/

popd