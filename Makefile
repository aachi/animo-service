install:
	dep ensure

lint:
	go fmt ./...
	goimports -w -d $(find . -type f -name '*.go' -not -path "./vendor/*")
	golint ./pkg/...

run:
	GOOGLE_APPLICATION_CREDENTIALS=~/path-to-certificate.json go run ./cmd/service/main.go -gcp-project=project-id -signature-secret=secret

docker-build:
	docker build . --tag animo-service

docker-run:
	docker run -e "SIGNATURE_SECRET=secret" -e "GCP_PROJECT=project-id" -e "GOOGLE_APPLICATION_CREDENTIALS=~/path-to-certificate.json" -p 8080:8080 animo-service