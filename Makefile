deploy: build image push

build-app:
	go build -o effective-status

image:
	docker build -t ialexin/effective-status:1.0 .

push:
	docker push ialexin/effective-status
