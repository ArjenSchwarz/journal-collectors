.PHONY: deps clean build package

deps:
	go get -u ./...

clean:
	rm -rf ./output

build:
	mkdir -p output/{github,instapaper,pinboard}
	cp github/config.yml output/github/
	cp instapaper/config.yml output/instapaper/
	cp pinboard/config.yml output/pinboard/
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o output/github/github ./github
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o output/instapaper/instapaper ./instapaper
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o output/pinboard/pinboard ./pinboard

package:
	sam package --template-file ./template.yaml --s3-bucket public.ig.nore.me --output-template-file packaged-template.yaml

deploy:
	sam deploy --template-file ./packaged-template.yaml --stack-name journal-collectors