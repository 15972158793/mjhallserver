export GOROOT=/usr/local/go
export GOBIN=$GOROOT/bin
export GOPATH=~/rjserver/baseserver/server:$GOROOT:~/jyserver/gopath
export PATH=$PATH:$GOROOT/bin:$GOBIN
export GOOS=darwin GOARCH=amd64

echo $(go version)
go clean

cd bin/server
go build main/login
echo build login ok !

go build main/game
echo build game ok !

go build main/center
echo build center ok !

cd ../backstage
go build main/backstage
echo build backstage ok !
