PLUGIN_NAME = dpcsoftware/oci-secrets
VERSION = 0.1
PLATFORM ?= amd64
TAG ?= $(VERSION)-$(PLATFORM)

all: create

clean:
	rm -rf ./plugin/rootfs

rootfs: clean
	docker buildx build --platform $(PLATFORM) -t oci-secrets-rootfs .
	docker create --name tmp oci-secrets-rootfs true
	mkdir -p ./plugin/rootfs
	docker export tmp | tar -x -C ./plugin/rootfs
	docker rm -vf tmp
	docker rmi oci-secrets-rootfs

create: rootfs
	docker plugin rm -f $(PLUGIN_NAME):$(TAG) || true
	docker plugin create $(PLUGIN_NAME):$(TAG) ./plugin

push: create
	docker plugin push $(PLUGIN_NAME):$(TAG)
