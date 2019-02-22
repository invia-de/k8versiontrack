FROM golang

# Add project directory to Docker image.
ADD . /go/src/github.com/invia-de/k8versiontrack
COPY ./scripts/kubetoken /usr/bin/kubetoken

ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET Az00P54fhK2SMggW
ENV KUBECTL_VERSION=v1.11.6
# Install kubectl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
    chmod +x kubectl && \
    mv kubectl /usr/bin/kubectl && \
    echo "source <(kubectl completion bash)" >> ~/.bashrc


WORKDIR /go/src/github.com/invia-de/k8versiontrack

RUN go get
RUN go build
EXPOSE 8888
CMD ./k8versiontrack
