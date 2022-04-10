IMAGES = setup elasticsearch kibana nginx
TAG ?= latest

# --- Docker images -----------------------------------------------------------
define build-and-tag
	docker build --file docker/$(1).Dockerfile --tag ghcr.io/scorestack/$(1):$(2) .
endef

define build-tag-push
	$(call build-and-tag,$(1),$(2))
	docker push ghcr.io/scorestack/$(1):$(2)
endef

.PHONY: build-image
build-image:
	$(foreach image,$(IMAGES),$(call build-and-tag,$(image),$(TAG));)

.PHONY: build-push-image
build-push-image:
	$(foreach image,$(IMAGES),$(call build-tag-push,$(image),$(TAG)))