package data

import (
	"context"
	"os"

	"auth/ent"
	"auth/ent/schema"
	"auth/internal/pkg/secrets"
)

func makeTheCreator(client *ent.Client) *ent.UserCreate {
	email := os.Getenv(`ADMIN_EMAIL`)
	if email == `` {
		panic(`env variable ADMIN_EMAIL was not set`)
	}
	phone := os.Getenv(`ADMIN_PHONE`)
	if phone == `` {
		panic(`env variable ADMIN_PHONE was not set`)
	}
	password := os.Getenv(`ADMIN_PASSWORD`)
	if password == `` {
		panic(`env variable ADMIN_PASSWORD was not set`)
	}
	return client.User.Create().
		SetDisplayName("Admin, the Creator").
		SetType(schema.TypeAdmin).
		SetEmail(email).
		SetPhone(phone).
		SetPasswordHash(secrets.MustMakeHash(password))
}

func SeedMainEntities(ctx context.Context, client *ent.Client) error {
	seeders := []func(context.Context, *ent.Client) error{
		seedTheCreator,
	}

	for _, seeder := range seeders {
		if err := seeder(ctx, client); err != nil {
			return err
		}
	}
	return nil
}

func seedTheCreator(ctx context.Context, client *ent.Client) error {
	_, err := client.User.Get(ctx, 1)
	if ent.IsNotFound(err) {
		_, err := makeTheCreator(client).Save(ctx)
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
