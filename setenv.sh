GOPATH=~/go
PWD=$(pwd)
export GOPATH=$GOPATH:$PWD
export PATH=$PATH:${GOPATH//://bin:}/bin
