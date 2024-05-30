import:
    cdk8s import

build:
    go build -o . ./...

generate: build
    ./k8s-home
    ls -l dist/

check-versions *ARGS: build
    ./k8s-home check-versions {{ARGS}}
