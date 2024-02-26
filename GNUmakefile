TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
HOSTNAME=codefresh.io
PKG_NAME=codefresh
NAMESPACE=app
BINARY=terraform-provider-${PKG_NAME}
OS_ARCH=$$(go env GOOS)_$$(go env GOARCH)
TFPLUGINDOCS_VERSION=v0.14.1

default: build

tools:
	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint
	GO111MODULE=on go install github.com/bflad/tfproviderlint/cmd/tfproviderlint

build: fmtcheck
	go install
	go build -o ${BINARY}

install: build
	mv ${BINARY} $(HOME)/go/bin/

equivalence: build
	@echo "==> Preparing equivalence tests"
	mkdir -p testing/equivalence/.plugins/registry.terraform.io/codefresh-io/codefresh/0.6.0/${OS_ARCH}/
	cp terraform-provider-codefresh testing/equivalence/.plugins/registry.terraform.io/codefresh-io/codefresh/0.6.0/${OS_ARCH}/

	cd testing/equivalence;\
	./update-test-cases.sh;\

	@echo "==> Running equivalence tests for terraform"
	cd testing/equivalence;\
	equivalence-testing update --binary=$$(which terraform) --goldens=results/terraform --tests=test_cases --rewrites=rewrites.jsonc

	@echo "==> Running equivalence tests for opentofu"
	cd testing/equivalence;\
	equivalence-testing update --binary=$$(which tofu) --goldens=results/opentofu --tests=test_cases --rewrites=rewrites.jsonc

	@echo "==> Comparing results"
	cd testing/equivalence;\
	./compare-results.sh;\

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w $(GOFMT_FILES)

fmtcheck: SHELL:=/bin/bash
fmtcheck:
	@echo "==> Checking that code complies with gofmt requirements..."

	@gofmt_files=$$(find . -name '*.go' | grep -v vendor | xargs gofmt -l -s); \
	if [[ -n $${gofmt_files} ]]; then\
	    echo 'gofmt needs running on the following files:';\
	    echo "$${gofmt_files}";\
	    echo "You can use the command: \`make fmt\` to reformat code.";\
	    exit 1;\
	fi;
lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

docs-prepare:
	@echo "==> Setting up tfplugindocs..."
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@${TFPLUGINDOCS_VERSION}

docs: docs-prepare
	@echo "==> Generating Provider Documentation..."
	tfplugindocs generate

.PHONY: build test testacc vet fmt fmtcheck lint tools test-compile docs docs-prepare

