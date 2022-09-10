package biz

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"auth/ent"
	"auth/internal/pkg/jwt"
	"auth/internal/pkg/template"
	"auth/internal/pkg/texts"
)

func (a *AuthUsecase) userByUsername(ctx context.Context, username string) (*ent.User, error) {
	var u *ent.User
	var err error
	if strings.Contains(username, `@`) {
		u, err = a.userRepo.FindByEmail(ctx, username)
	} else {
		u, err = a.userRepo.FindByPhone(ctx, username)
	}
	if ent.IsNotFound(err) {
		return nil, errors.New(
			template.MustInterpolate(
				texts.UsernameNotFound, map[string]any{
					"username": username,
				},
			),
		)
	}
	return u, err
}

func historyModel(userID int, event string, stats *Stats) *ent.History {
	history := &ent.History{
		UserID: userID,
		Event:  event,
	}
	if stats != nil {
		history.IP = &stats.IP
		history.UserAgent = &stats.UserAgent
	}
	return history
}

func sessionModelByLogin(userID int, passwordHash string, dto *LoginDTO) *ent.Session {
	session := &ent.Session{
		UserID:    userID,
		Token:     makeToken(dto.Username, passwordHash),
		IsActive:  true,
		ExpiredAt: expiredAt(dto.Remember),
	}
	extendSessionStats(session, dto.Stats)
	return session
}

func sessionModelByLoginByCode(userID int, code string, dto *LoginByCodeDTO) *ent.Session {
	session := &ent.Session{
		UserID:    userID,
		Token:     makeToken(dto.Username, code),
		IsActive:  true,
		ExpiredAt: expiredAt(dto.Remember),
	}
	extendSessionStats(session, dto.Stats)
	return session
}

func codeModel(userID int, code string) *ent.Code {
	return &ent.Code{
		UserID:  userID,
		Content: code,
	}
}

func makeToken(username, secret string) string {
	return jwt.Make(username + time.Now().String() + secret)
}

func mustMakeCode(length int) string {
	max := math.Pow(10, float64(length))
	res, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`%0`+fmt.Sprintf(`%d`, length)+`d`, res.Int64())
}

func makeHash(length int, username, secret string) string {
	h := sha1.New()
	h.Write([]byte(username + time.Now().String() + secret))
	start := 6
	return hex.EncodeToString(h.Sum(nil))[start : start+length]
}

func makeRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)[:length]
}

func expiredAt(remember bool) time.Time {
	expiredInterval := sessionExpiredIntervalShort
	if remember {
		expiredInterval = sessionExpiredIntervalLong
	}
	return time.Now().Add(expiredInterval)
}

func extendSessionStats(session *ent.Session, stats *Stats) {
	if session != nil && stats != nil {
		session.IP = stats.IP
		session.UserAgent = stats.UserAgent
		session.DeviceID = stats.DeviceID
	}
}
