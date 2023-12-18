default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: generate
generate:
	go generate ./...

all: build api-server api-client

.PHONY: build
build: tmp
	go build -o tmp/terraform-provider-redirect-store ./main.go

.PHONY: api-server
api-server: tmp
	go build -o tmp/api-server ./api/cmd/server

.PHONY: api-client
api-client: tmp
	go build -o tmp/api-client ./api/cmd/client

tmp:
	mkdir -p $@
