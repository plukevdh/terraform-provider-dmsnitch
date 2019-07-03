build-linux: clean
	docker build -t dms-builder .
	docker run --rm -v "${PWD}":/go/src/terraform-provider-dmsnitch dms-builder go build

build: clean
	go get -d -v ./...
	go install -v ./...
	go build

install:
	cp terraform-provider-dmsnitch  ~/.terraform.d/plugins/

clean:
	rm -f terraform-provider-dmsnitch
