FROM golang:1.20.5-alpine as build

RUN mkdir /build
WORKDIR /build
COPY src/ ./

RUN go get
RUN go build

FROM scratch

COPY --from=build /etc/ssl/cert.pem /etc/ssl/cert.pem
COPY --from=build /build/oci-secrets-docker /bin/oci-secrets
