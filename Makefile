IMG ?= docker.io/imroc/openkruise-dummy-gameserver:latest

build:
	docker buildx build --push -t ${IMG} .
