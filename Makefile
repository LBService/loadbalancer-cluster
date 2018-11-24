all: build

VERSION?=0.0.1


fmt:
	find . -type f -name "*.go" | grep -v "./vendor*" | xargs gofmt -s -w

build: clean fmt
	go build -o lbaas -ldflags "-X github.com/LBService/loadbalancer-cluster/pkg/version.gitCommit=`git rev-parse HEAD` -X github.com/LBService/loadbalancer-cluster/pkg/version.buildDate=`date '+%Y-%m-%dT%H:%M:%S'`" -a github.com/LBService/loadbalancer-cluster/cmd/lbaas

local-run:
	MY_POD_NAMESPACE=default MY_POD_NAME=testpod  ./lbaas
    
clean:
	rm -f lbaas cmd/lbaas/lbaas
