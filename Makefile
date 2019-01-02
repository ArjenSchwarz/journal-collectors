.PHONY: all
all: build package deploy

.PHONY: github-actions
github-actions: prep-actions build

.PHONY: build
build: deps test clean compile

.PHONY: deps
deps:
	go get -v ./...

.PHONY: test
test:
	go get -u golang.org/x/lint/golint
	golint ./...
	go test ./...

.PHONY: clean
clean:
	rm -rf ./output

.PHONY: compile
compile:
	mkdir -p output/github output/instapaper output/pinboard
	cp github/config.yml output/github/
	cp instapaper/config.yml output/instapaper/
	cp pinboard/config.yml output/pinboard/
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o output/github/github ./github
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o output/instapaper/instapaper ./instapaper
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o output/pinboard/pinboard ./pinboard

.PHONY: aws
aws: package deploy

.PHONY: package
package:
	aws cloudformation package --template-file ./template.yaml --s3-bucket public.ig.nore.me --output-template-file packaged-template.yaml

.PHONY: deploy
deploy:
	aws cloudformation deploy --template-file ./packaged-template.yaml --stack-name journal-collectors
