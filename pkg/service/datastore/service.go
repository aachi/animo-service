package datastore

import (
	"context"
	"errors"

	"cloud.google.com/go/datastore"
	"github.com/go-kit/kit/log"
	"google.golang.org/api/iterator"

	"github.com/revas/animo-service/pkg"
)

type GoogleDatastoreAnimoService struct {
	Logger log.Logger
	Client *datastore.Client
}

// Ensure InMemoryAnimoService implements the animo.AnimoService interface.
var _ animo.AnimoService = &GoogleDatastoreAnimoService{}

func (svc *GoogleDatastoreAnimoService) GetOrCreateProfile(ctx context.Context, identity string) (*animo.Profile, error) {
	profile, err := findProfileByIdentity(ctx, svc, identity)
	if err != nil {
		return nil, err
	}
	if profile.ID == "" {
		profile, err = createProfileFromIdentity(ctx, svc.Client, identity)
		if err != nil {
			return nil, err
		}
	}
	return profile, nil
}

func (svc *GoogleDatastoreAnimoService) ResolveProfilesAliases(context context.Context, profilesAliases []string) ([]string, error) {
	var profilesIds []string
	for _, alias := range profilesAliases {
		var profile *animo.Profile
		var err error
		if alias == "me" {
			profile, err = svc.GetOrCreateProfile(context, context.Value("Identity").(string))
		} else {
			profile, err = findProfileByAlias(context, svc, alias)
		}
		if err != nil {
			return nil, err
		}
		profilesIds = append(profilesIds, profile.ID)
	}
	return profilesIds, nil
}

func (svc *GoogleDatastoreAnimoService) GetProfiles(context context.Context, profilesIds []string) ([]*animo.Profile, error) {
	profilesKeys := makeDatastoreKeysFromIds(profilesIds)

	profiles := make([]*animo.Profile, len(profilesKeys))
	err := svc.Client.GetMulti(context, profilesKeys, profiles)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (svc *GoogleDatastoreAnimoService) SearchProfiles(context context.Context, filter string) ([]*animo.Profile, error) {
	query := datastore.NewQuery("Profiles").
		Filter("Alias >=", filter).
		Order("Alias").
		Limit(5)

	it := svc.Client.Run(context, query)
	var profiles []*animo.Profile
	for {
		var profile animo.Profile
		_, err := it.Next(&profile)
		if err == iterator.Done {
			break
		}
		if err != nil {
			svc.Logger.Log("error", err.Error())
		}
		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

func (svc *GoogleDatastoreAnimoService) UpdateProfiles(ctx context.Context, profilesIds []string, profiles []*animo.Profile) ([]*animo.Profile, error) {
	profilesKeys := makeDatastoreKeysFromIds(profilesIds)

	persistedProfiles := make([]*animo.Profile, len(profilesKeys))
	err := svc.Client.GetMulti(ctx, profilesKeys, persistedProfiles)
	if err != nil {
		return nil, err
	}

	for index, persistedProfile := range persistedProfiles {
		updatedProfile := profiles[index]
		if persistedProfile.Alias != updatedProfile.Alias && updatedProfile.Alias != "me" {
			profile, err := findProfileByAlias(ctx, svc, updatedProfile.Alias)
			if err != nil {
				return nil, err
			}
			if profile.ID != "" {
				return nil, errors.New("profile alias is not available")
			}
			persistedProfile.Alias = updatedProfile.Alias
		}
		if updatedProfile.Name == "" || updatedProfile.Email == "" {
			return nil, errors.New("profile values are empty")
		}
		persistedProfile.Name = updatedProfile.Name
		persistedProfile.Email = updatedProfile.Email
		persistedProfile.Picture = updatedProfile.Picture
	}

	_, err = svc.Client.PutMulti(ctx, profilesKeys, persistedProfiles)
	if err != nil {
		return nil, err
	}

	return persistedProfiles, nil
}
