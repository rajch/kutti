# Bump these on release, and for now update the deployment files
VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
BUILD_NUMBER  ?= 12

IMAGE_TAG ?= $(VERSION_MAJOR).$(VERSION_MINOR).$(BUILD_NUMBER)
REGISTRY_USER ?= rajchaudhuri

.PHONY: all
all: kutticmd

out/kutti-localprovisioner: cmd/kutti-localprovisioner/main.go
	CGO_ENABLED=0 go build -o out/kutti-localprovisioner cmd/kutti-localprovisioner/main.go


.PHONY: localprovisioner
localprovisioner: out/kutti-localprovisioner

.PHONY: localprovisioner-image
localprovisioner-image: out/kutti-localprovisioner build/package/kutti-localprovisioner/local.Dockerfile
	docker image build -t $(REGISTRY_USER)/kutti-localprovisioner:$(IMAGE_TAG) -f build/package/kutti-localprovisioner/local.Dockerfile .

KUTTICMDFILES = cmd/kutti/main.go \
				cmd/kutti/cmd/*.go \
				pkg/clustermanager/*.go \
				pkg/vboxdriver/*.go \
				pkg/core/*.go \
				internal/pkg/configfilemanager/*.go \
				internal/pkg/fileutils/*.go \
				internal/pkg/kuttilog/*.go

out/kutti: $(KUTTICMDFILES)
	CGO_ENABLED=0 go build -o $@ $<

.PHONY: kutticmd
kutticmd: out/kutti

out/kutti.exe: $(KUTTICMDFILES)
	CGO_ENABLED=0 GOOS=windows go build -o $@ $<

.PHONY: kutticmd-windows
kutticmd-windows: out/kutti.exe

.PHONY: clean
clean:
	rm -rf out/*
	docker image rm $(REGISTRY_USER)/kutti-localprovisioner:$(IMAGE_TAG)