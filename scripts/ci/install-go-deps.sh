#!/bin/bash

set -v

sudo yum install -y https://www.rdoproject.org/repos/rdo-release.rpm
sudo yum install -y perl-Graph-Easy

export GOROOT=/opt/go
export GOPATH=$WORKSPACE
export PATH=$PATH:$GOROOT/bin

# Get the Go dependencies
go get -f -u github.com/axw/gocov/gocov
go get -f -u github.com/mattn/goveralls
go get -f -u golang.org/x/tools/cmd/cover
go get -f -u github.com/golang/lint/golint
go get -f -u github.com/tebeka/go2xunit

export PATH=$PATH:$GOPATH/bin

# speedup govendor sync command
curl -o /tmp/vendor.tgz http://46.231.132.68:8080/v1/AUTH_0ec9e4f4f3044236b4d18536ccfcb182/skydive/vendor/vendor.tar.gz

pushd ${GOPATH}/src/github.com/skydive-project/skydive
go get -f -u github.com/kardianos/govendor
govendor sync -n | perl -pe 's|fetch \"(.*)\"$|vendor/\1|g' | sort -u > vendor.fetch.list
cat vendor.fetch.list | xargs tar -xvzf /tmp/vendor.tgz --exclude "vendor/vendor.json"
# remove installed
find vendor/ -mindepth 2 -type f | xargs dirname | sort -u > vendor.installed.list
echo "package to be removed/cleanup"
diff -u vendor.fetch.list vendor.installed.list | grep '^\+v' | perl -pe 's|^\+(.*)|\1|' | tee /dev/stdout | xargs -n 1 rm -rf
rm -f vendor.fetch.list vendor.installed.list
popd
