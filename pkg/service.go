package animo

import "context"

type Profile struct {
	ID       string `json:"-"`
	Identity string `json:"-"`
	Alias    string `json:"alias"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
}

type AnimoService interface {
	GetOrCreateProfile(context context.Context, userIdentity string) (*Profile, error)
	ResolveAliases(context context.Context, profilesAliases []string) ([]string, error)
	GetProfiles(context context.Context, profilesIds []string) ([]*Profile, error)
	SearchProfiles(context context.Context, filter string) ([]*Profile, error)
	UpdateProfiles(context context.Context, profilesIds []string, profiles []*Profile) ([]*Profile, error)
}
