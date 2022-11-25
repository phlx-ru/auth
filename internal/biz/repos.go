package biz

import (
	"context"
	"time"

	"auth/ent"
)

//go:generate mockgen -source ./${GOFILE} -destination ./repos_mock_test.go -package ${GOPACKAGE}

type SessionRepo interface {
	Create(context.Context, *ent.Session) (*ent.Session, error)
	FindByToken(ctx context.Context, token string) (*ent.Session, error)
}

type CodeRepo interface {
	Create(context.Context, *ent.Code) (*ent.Code, error)
	FindForUser(ctx context.Context, userID int) (*ent.Code, error)
}

type HistoryRepo interface {
	Create(context.Context, *ent.History) (*ent.History, error)
	FindUserEvents(ctx context.Context, userID, limit, offset int) ([]*ent.History, error)
	FindLastUserEvents(
		ctx context.Context,
		userID int,
		types []string,
		interval time.Duration,
	) ([]*ent.History, error)
}

type UserRepo interface {
	Create(context.Context, *ent.User) (*ent.User, error)
	Update(context.Context, *ent.User) (*ent.User, error)
	FindByID(ctx context.Context, id int) (*ent.User, error)
	FindByEmail(ctx context.Context, email string) (*ent.User, error)
	FindByPhone(ctx context.Context, phone string) (*ent.User, error)
	Activate(ctx context.Context, userID int) (*ent.User, error)
	Deactivate(ctx context.Context, userID int) (*ent.User, error)
}
