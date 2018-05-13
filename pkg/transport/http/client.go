package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	khttp "github.com/go-kit/kit/transport/http"

	"github.com/revas/animo-service/pkg"
)

func EncodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func DecodeResolveAliasesResponse(_ context.Context, response *http.Response) (interface{}, error) {
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	var resp animo.ResolveAliasesResponse
	err := json.NewDecoder(response.Body).Decode(&resp)
	return resp, err
}

func CreateResolveProfilesAliasesClientEndpoint(url *url.URL) endpoint.Endpoint {
	return khttp.NewClient(
		"GET",
		url,
		EncodeRequest,
		DecodeResolveAliasesResponse,
	).Endpoint()
}
