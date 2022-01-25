default: testacc

VERSION=1.1.1
REPO=susan.to/terraform/contentful

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build:
	go build

install: build
	mkdir -p ~/.terraform.d/plugins/${REPO}/${VERSION}/darwin_amd64
	mv terraform-provider-contentful ~/.terraform.d/plugins/${REPO}/${VERSION}/darwin_amd64
