deploy: build image push

build-app:
	go build -o effective-status

image:
	docker build -t ialexin/effective-status .

push:
	docker push ialexin/effective-status
