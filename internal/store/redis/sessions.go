package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"authservice/internal/model"
)

func (rdb *Redis) CreateSession(ctx context.Context, session *model.Session) error {
	return rdb.db.Set(
		ctx,
		sessionKey(session.Username),
		session,
		0,
	).Err()
}

func (rdb *Redis) GetSession(ctx context.Context, userID string) (*model.Session, error) {
	var session model.Session

	bytes, err := rdb.db.Get(ctx, sessionKey(userID)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if err := json.Unmarshal(bytes, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, err
}

func (rdb *Redis) DeleteSession(ctx context.Context, userID string) error {
	return rdb.db.Del(
		ctx,
		sessionKey(userID),
	).Err()
}

func (rdb *Redis) ListSessions(ctx context.Context, userID string) ([]*model.Session, error) {
	ks, err := rdb.db.Keys(
		ctx,
		sessionKeyList(),
	).Result()
	if err != nil {
		return nil, err
	}

	sessions := make([]*model.Session, 0, len(ks))
	for _, key := range ks {
		userID, _ := strings.CutPrefix(key, SessionKey)

		session, err := rdb.GetSession(ctx, userID)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)

	}

	return sessions, nil
}

func (rdb *Redis) ListSessionsUsingScan(ctx context.Context, userID string) ([]*model.Session, error) {
	var cursor uint64
	var keys []string

	for {
		var err error
		ks, cursor, err := rdb.db.Scan(ctx, cursor, sessionKeyList(), 10).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, ks...)
		if cursor == 0 {
			break
		}

	}

	sessions := make([]*model.Session, 0, len(keys))
	for _, key := range keys {
		userID, _ := strings.CutPrefix(key, SessionKey)

		session, err := rdb.GetSession(ctx, userID)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)

	}

	return sessions, nil
}

func sessionKey(userID string) string {
	return SessionKey + userID
}

func sessionKeyList() string {
	return SessionKey + "*"
}
