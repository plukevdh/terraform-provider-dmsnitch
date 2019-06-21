build-linux:
	docker build -t dms-builder .
	docker run --rm -v "${PWD}":/go/src/terraform-provider-dmsnitch dms-builder go build

build:
	go get -d -v ./...
	go install -v ./...
	go build

install:
	cp terraform-provider-dmsnitch  ~/.terraform.d/plugins/
