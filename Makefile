# Bump these on release, and for now update the deployment files
VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
BUILD_NUMBER  ?= 12

IMAGE_TAG ?= $(VERSION_MAJOR).$(VERSION_MINOR).$(BUILD_NUMBER)
REGISTRY_USER ?= rajchaudhuri

.PHONY: all
all: localprovisioner localprovisioner-image

out/kutti-localprovisioner: cmd/kutti-localprovisioner/main.go
	CGO_ENABLED=0 go build -o out/kutti-localprovisioner cmd/kutti-localprovisioner/main.go


.PHONY: localprovisioner
localprovisioner: out/kutti-localprovisioner

.PHONY: localprovisioner-image
localprovisioner-image: out/kutti-localprovisioner build/package/kutti-localprovisioner/local.Dockerfile
	docker image build -t $(REGISTRY_USER)/kutti-localprovisioner:$(IMAGE_TAG) -f build/package/kutti-localprovisioner/local.Dockerfile .

out/kutti: cmd/kutti/main.go cmd/kutti/cmd/*.go
	CGO_ENABLED=0 go build -o $@ $<


.PHONY: kutticmd
kutticmd: out/kutti

.PHONY: clean
clean:
	rm -rf out/*
	docker image rm $(REGISTRY_USER)/kutti-localprovisioner:$(IMAGE_TAG)