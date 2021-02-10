GOBIN ?= $(shell go env GOPATH)/bin
VERSION := $$(make -s show-version)

.PHONY: show-version
show-version: $(GOBIN)/gobump
	@gobump show -r .

PHONY: deploy
deploy:
	gcloud app deploy -q

.PHONY: tag
tag:
	git tag -a "v$(VERSION)" -m "Release $(VERSION)"
	git push --tags
