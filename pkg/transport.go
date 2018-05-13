package animo

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type ResolveAliasesRequest struct {
	ProfilesAliases []string `json:"profilesAliases"`
}

type ResolveAliasesResponse struct {
	ProfilesIds []string `json:"profilesIds"`
	Err         string   `json:"err,omitempty"`
}

type GetProfilesRequest struct {
	ProfilesAliases []string `json:"profilesAliases"`
}

type GetProfilesResponse struct {
	Profiles []*Profile `json:"profiles"`
	Err      string     `json:"err,omitempty"`
}

type SearchProfilesRequest struct {
	Filter string `json:"filter"`
}

type SearchProfilesResponse struct {
	Profiles []*Profile `json:"profiles"`
	Err      string     `json:"err,omitempty"`
}

type UpdateProfilesRequest struct {
	ProfilesAliases []string   `json:"profilesAliases"`
	Profiles        []*Profile `json:"profiles"`
}

type UpdateProfilesResponse struct {
	Profiles []*Profile `json:"profiles"`
	Err      string     `json:"err,omitempty"`
}

type Endpoints struct {
	ResolveAliasesEndpoint endpoint.Endpoint
	GetProfilesEndpoint    endpoint.Endpoint
	SearchProfilesEndpoint endpoint.Endpoint
	UpdateProfilesEndpoint endpoint.Endpoint
}

func MakeResolveAliasesEndpoint(svc AnimoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ResolveAliasesRequest)
		ids, err := svc.ResolveAliases(ctx, req.ProfilesAliases)
		if err != nil {
			return ResolveAliasesResponse{[]string{}, err.Error()}, nil
		}
		return ResolveAliasesResponse{ids, ""}, nil
	}
}

func MakeGetProfilesEndpoint(svc AnimoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetProfilesRequest)
		profilesIds, err := svc.ResolveAliases(ctx, req.ProfilesAliases)
		profiles, err := svc.GetProfiles(ctx, profilesIds)
		if err != nil {
			return GetProfilesResponse{[]*Profile{}, err.Error()}, nil
		}
		return GetProfilesResponse{profiles, ""}, nil
	}
}

func MakeSearchProfilesEndpoint(svc AnimoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SearchProfilesRequest)
		profiles, err := svc.SearchProfiles(ctx, req.Filter)
		if err != nil {
			return SearchProfilesResponse{[]*Profile{}, err.Error()}, nil
		}
		return SearchProfilesResponse{profiles, ""}, nil
	}
}

func MakeUpdateProfilesEndpoint(svc AnimoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateProfilesRequest)

		if len(req.ProfilesAliases) != 1 || req.ProfilesAliases[0] != "me" {
			return UpdateProfilesResponse{[]*Profile{}, "access is allowed only to self resource"}, nil
		}
		if len(req.ProfilesAliases) != len(req.Profiles) {
			return nil, errors.New("aliases and profiles have different lengths")
		}

		profilesIds, err := svc.ResolveAliases(ctx, req.ProfilesAliases)
		profiles, err := svc.UpdateProfiles(ctx, profilesIds, req.Profiles)
		if err != nil {
			return UpdateProfilesResponse{[]*Profile{}, err.Error()}, nil
		}
		return UpdateProfilesResponse{profiles, ""}, nil
	}
}
