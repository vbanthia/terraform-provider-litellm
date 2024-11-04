HOSTNAME=registry.terraform.io
NAMESPACE=local
NAME=litellm
VERSION=1.0.0
OS_ARCH=darwin_amd64

default: install

build:
	go build -o terraform-provider-${NAME}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv terraform-provider-${NAME} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/terraform-provider-${NAME}_v${VERSION}

.PHONY: build install
