FROM alpine:3.11.3
CMD [ "kubetrack" ]
COPY release/kubetrack-* /bin/kubetrack
