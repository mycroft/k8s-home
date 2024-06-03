import:
    cdk8s import

build:
    go build -o . ./...

diff: generate
    sh contrib/diff.sh

generate *ARGS: build
    ./k8s-home {{ARGS}}
    ls -l dist/

lint:
    golangci-lint run -c .golangci.yaml

check-versions *ARGS: build
    ./k8s-home check-versions {{ARGS}}
