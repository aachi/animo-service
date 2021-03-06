package http

import (
	"context"
	"encoding/json"
	"net/http"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/revas/animo-service/pkg"
)

func DecodeResolveProfilesAliasesRequest(_ context.Context, request *http.Request) (interface{}, error) {
	var decoded animo.ResolveProfilesAliasesRequest
	if err := json.NewDecoder(request.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func DecodeInternalGetProfilesRequest(_ context.Context, request *http.Request) (interface{}, error) {
	var decoded animo.InternalGetProfilesRequest
	if err := json.NewDecoder(request.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func DecodeGetProfilesRequest(_ context.Context, request *http.Request) (interface{}, error) {
	var decoded animo.GetProfilesRequest
	if err := json.NewDecoder(request.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func DecodeSearchProfilesRequest(_ context.Context, request *http.Request) (interface{}, error) {
	var decoded animo.SearchProfilesRequest
	if err := json.NewDecoder(request.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func DecodeUpdateProfilesRequest(_ context.Context, request *http.Request) (interface{}, error) {
	var decoded animo.UpdateProfilesRequest
	if err := json.NewDecoder(request.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func MakeHandler(logger log.Logger, endpoints animo.Endpoints) http.Handler {
	options := []kithttp.ServerOption{
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerErrorLogger(logger),
	}

	resolveProfilesAliasesHandler := kithttp.NewServer(
		endpoints.ResolveAliasesEndpoint,
		DecodeResolveProfilesAliasesRequest,
		EncodeResponse,
		options...,
	)

	internalGetProfilesHandler := kithttp.NewServer(
		endpoints.GetProfilesEndpoint,
		DecodeInternalGetProfilesRequest,
		EncodeResponse,
		options...,
	)

	getProfilesHandler := kithttp.NewServer(
		endpoints.GetProfilesEndpoint,
		DecodeGetProfilesRequest,
		EncodeResponse,
		options...,
	)

	searchProfilesHandler := kithttp.NewServer(
		endpoints.SearchProfilesEndpoint,
		DecodeSearchProfilesRequest,
		EncodeResponse,
		options...,
	)

	updateProfilesHandler := kithttp.NewServer(
		endpoints.UpdateProfilesEndpoint,
		DecodeUpdateProfilesRequest,
		EncodeResponse,
		options...,
	)

	r := mux.NewRouter()

	r.Handle("/internal/animo.ResolveProfilesAliases/", resolveProfilesAliasesHandler).Methods("POST")
	r.Handle("/internal/animo.GetProfiles/", internalGetProfilesHandler).Methods("POST")
	r.Handle("/animo.GetProfiles/", getProfilesHandler).Methods("POST")
	r.Handle("/animo.SearchProfiles/", searchProfilesHandler).Methods("POST")
	r.Handle("/animo.UpdateProfiles/", updateProfilesHandler).Methods("POST")

	return r
}
