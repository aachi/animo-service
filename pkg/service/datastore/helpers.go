package datastore

import (
	"context"
	"crypto/rand"

	"cloud.google.com/go/datastore"
	"github.com/oklog/ulid"
	"github.com/revas/animo-service/pkg"
	"google.golang.org/api/iterator"
)

func makeDatastoreKeysFromIds(profilesIds []string) []*datastore.Key {
	var keys []*datastore.Key
	for _, id := range profilesIds {
		key := datastore.NameKey("Profiles", id, nil)
		keys = append(keys, key)
	}
	return keys
}

func findProfileByIdentity(ctx context.Context, svc *GoogleDatastoreAnimoService, identity string) (*animo.Profile, error) {
	query := datastore.NewQuery("Profiles").
		Filter("Identity =", identity).
		Order("Identity").
		Limit(1)

	profile := &animo.Profile{}
	it := svc.Client.Run(ctx, query)
	_, err := it.Next(profile)
	if err != iterator.Done && err != nil {
		return nil, err
	}
	return profile, nil
}

func findProfileByAlias(ctx context.Context, svc *GoogleDatastoreAnimoService, alias string) (*animo.Profile, error) {
	query := datastore.NewQuery("Profiles").
		Filter("Alias =", alias).
		Order("Alias").
		Limit(1)

	profile := &animo.Profile{}
	it := svc.Client.Run(ctx, query)
	_, err := it.Next(profile)
	if err != iterator.Done && err != nil {
		return nil, err
	}
	return profile, nil
}

func createProfileFromIdentity(ctx context.Context, client *datastore.Client, identity string) (*animo.Profile, error) {
	profile := &animo.Profile{
		ID:       ulid.MustNew(ulid.Now(), rand.Reader).String(),
		Identity: identity,
		Name:     "",
		Picture:  "",
		Alias:    ulid.MustNew(ulid.Now(), rand.Reader).String(),
	}
	_, err := client.Put(ctx, datastore.NameKey("Profiles", profile.ID, nil), profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
