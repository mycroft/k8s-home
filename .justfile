build:
    go build -o . ./...

generate: build
    ./k8s-home
    ls -l dist/

check-versions: build
    ./k8s-home -check-version
