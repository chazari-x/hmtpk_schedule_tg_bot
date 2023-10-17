IMAGE_NAME=chazari/hmtpk_schedule
TAG=latest

.PHONY: build push

build:
	docker build -t $(IMAGE_NAME):$(TAG) .

push: build
	docker push $(IMAGE_NAME):$(TAG)
