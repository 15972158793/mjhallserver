export GOROOT=D:/Go
export GOBIN=D:/Develop/GoWork/bin
export GOPATH=D:/Develop/GoWork
export PATH=$PATH:$GOROOT/bin:$GOBIN
export GOOS=$1 GOARCH=amd64

echo $(go version)
go clean

go install $2
