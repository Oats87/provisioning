ARG UBI_IMAGE=registry.access.redhat.com/ubi7/ubi-minimal:latest
ARG GO_IMAGE=briandowns/rancher-build-base:v0.1.1

FROM ${UBI_IMAGE} as ubi

FROM ${GO_IMAGE} as builder
ARG TAG=""
RUN apt update     && \
    apt upgrade -y && \
    apt install -y ca-certificates git
RUN git clone --depth=1 http://github.com/rancher/provisioning
RUN cd provisioning && \
    git fetch --all --tags --prune && \
    go build
RUN echo $(pwd) && ls

FROM ubi
RUN microdnf update -y && \
    rm -rf /var/cache/yum
ENV KUBECONFIG /root/.kube/config
COPY --from=builder /go/provisioning/provisioning /usr/local/bin

ENTRYPOINT ["provisioning"]