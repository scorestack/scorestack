# --- Docker images -----------------------------------------------------------
.PHONY: build-image-setup
build-image-setup:
	docker build --file docker/setup.Dockerfile --tag scorestack-setup:latest docker/

.PHONY: build-image-elasticsearch
build-image-elasticsearch:
	docker build --file docker/elasticsearch.Dockerfile --tag scorestack-elasticsearch:latest docker/

.PHONY: build-image-kibana
build-image-kibana:
	docker build --file docker/kibana.Dockerfile --tag scorestack-kibana:latest .

.PHONY: build-image-nginx
build-image-nginx:
	docker build --file docker/nginx.Dockerfile --tag scorestack-nginx:latest docker/

.PHONY: build-image
build-image: build-image-setup build-image-elasticsearch build-image-kibana build-image-nginx

.PHONY: build-image-release
build-image-release: build-image
	docker tag scorestack-setup:latest scorestack-setup:0.8.2
	docker tag scorestack-elasticsearch:latest scorestack-elasticsearch:0.8.2
	docker tag scorestack-kibana:latest scorestack-kibana:0.8.2
	docker tag scorestack-nginx:latest scorestack-nginx:0.8.2