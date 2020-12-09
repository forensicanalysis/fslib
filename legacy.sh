go mod download
go mod vendor

rm -rf path
mkdir -p path/src/github.com/forensicanalysis/fslib go1.9.7 go1.2.2

cp -r vendor/* path/src
zip -r path.zip *
unzip path.zip -d path/src/github.com/forensicanalysis/fslib
rm path.zip

curl -Lso go1.9.7.darwin-amd64.tar.gz https://golang.org/dl/go1.9.7.darwin-amd64.tar.gz
tar xfvz go1.9.7.darwin-amd64.tar.gz
mv go go1.9.7/root
./go1.9.7/root/bin/go version
GOPATH=$PWD/path GOROOT=$PWD/go1.9.7/root ./go1.9.7/root/bin/go build -v github.com/forensicanalysis/fslib/...

curl -Lso go1.2.2.darwin-amd64.tar.gz https://golang.org/dl/go1.2.2.darwin-amd64-osx10.8.tar.gz
tar xfvz go1.2.2.darwin-amd64.tar.gz
mv go go1.2.2/root
./go1.2.2/root/bin/go version
GOPATH=$PWD/path GOROOT=$PWD/go1.2.2/root ./go1.2.2/root/bin/go build -v github.com/forensicanalysis/fslib/...
