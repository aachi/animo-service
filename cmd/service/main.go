package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/datastore"
	"github.com/go-kit/kit/log"

	"github.com/revas/animo-service/pkg"
	revasjwt "github.com/revas/animo-service/pkg/auth/jwt"
	sdatastore "github.com/revas/animo-service/pkg/service/datastore"
	thttp "github.com/revas/animo-service/pkg/transport/http"
)

func main() {
	var (
		GCPCertificatePath = flag.String("gcp-certificate", "", "Google Cloud Platform certificate .json file to connect with cloud resources. More info on https://cloud.google.com/docs/authentication/production.")
		GCPProjectID       = flag.String("gcp-project", "", "Google Cloud Platform project ID.")
		HS256SigningKey    = flag.String("signature-secret", "", "HS256 JWT token signing key.")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	if *GCPProjectID == "" || *HS256SigningKey == "" {
		logger.Log("error", "Please provide a GCP Project ID and a Token Signing Key.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		if *GCPCertificatePath != "" {
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", *GCPCertificatePath)
		} else {
			logger.Log("error", "Please provide a GCP Certificate.")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	ctx := context.Background()

	validateToken := revasjwt.MakeAuthenticatorMiddleware(*HS256SigningKey)

	var svc animo.AnimoService
	if *GCPProjectID != "" {
		client, err := datastore.NewClient(ctx, *GCPProjectID)
		if err != nil {
			panic(err)
		}
		svc = &sdatastore.GoogleDatastoreAnimoService{
			Logger: logger,
			Client: client,
		}
	}

	endpoints := animo.Endpoints{
		ResolveAliasesEndpoint: validateToken(animo.MakeResolveAliasesEndpoint(svc)),
		GetProfilesEndpoint:    validateToken(animo.MakeGetProfilesEndpoint(svc)),
		SearchProfilesEndpoint: validateToken(animo.MakeSearchProfilesEndpoint(svc)),
		UpdateProfilesEndpoint: validateToken(animo.MakeUpdateProfilesEndpoint(svc)),
	}

	handler := thttp.MakeHandler(logger, endpoints)

	errs := make(chan error, 2)

	go func() {
		errs <- http.ListenAndServe(":8080", handler)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}
