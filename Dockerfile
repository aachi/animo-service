FROM golang:1.10 AS builder

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/revas/animo-service
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN bash -c 'ls'
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app ./cmd/service

FROM alpine:3.5
COPY --from=builder /app ./

ENTRYPOINT ./app -gcp-project=${GCP_PROJECT} -signature-secret=${SIGNATURE_SECRET}
