export GOROOT=/usr/local/go
export GOBIN=$GOROOT/bin
export GOPATH=~/rjserver/baseserver/server:$GOROOT:~/jyserver/gopath
export PATH=$PATH:$GOROOT/bin:$GOBIN
export GOOS=linux GOARCH=amd64

echo $(go version)
go clean

cd bin/server
go build -ldflags "-w -s" main/login
echo build login ok !

go build -ldflags "-w -s" main/game
echo build game ok !

go build -ldflags "-w -s" main/center
echo build center ok !

cd ../backstage
go build -ldflags "-w -s" main/backstage
echo build backstage ok !
