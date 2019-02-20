IMAGE := k8sops:latest

start_container:
	docker run \
		--rm \
		-it \
		-p 8888:8888 \
		--name k8sops \
	${IMAGE}
build_container:
	docker build -t k8sops .
run_container:
	docker run \
                --rm \
                -it \
                -v ${PWD}:/go/src/github.com/invia-de/K8VersionTrack \
                -p 8888:8888 \
                --name k8sops \
        ${IMAGE} \
	/bin/bash

