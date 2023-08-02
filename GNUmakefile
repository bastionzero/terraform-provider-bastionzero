TEST?=$$(go list ./...)
GOFMT_FILES?=$$(find . -name '*.go')
PKG_NAME?=bastionzero
ACCTEST_TIMEOUT?=120m
ACCTEST_PARALLELISM?=2

default: build

build: fmtcheck
	go install .

generate:
	rm -rf docs/
	mkdir docs
	go generate ./...

test: fmtcheck
	go test -count=1 $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test -count=1 -timeout=30s -parallel=4

testacc: 
	TF_ACC=1 go test -v ./$(PKG_NAME)/... -count=1 -timeout $(ACCTEST_TIMEOUT) -parallel=$(ACCTEST_PARALLELISM)

vet:
	@echo "go vet ."
	@go vet $$(go list ./...) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

goimports:
	@echo "==> Fixing imports code with goimports..."
	@find . -name '*.go' | grep -v generator-resource-id | while read f; do goimports -w "$$f"; done

fmt:
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

website:
	@echo "Use this site to preview markdown rendering: https://registry.terraform.io/tools/doc-preview"

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test ./bastionzero/sweep/... -v -sweep=1

.PHONY: build generate test testacc vet goimports fmt fmtcheck website sweep
