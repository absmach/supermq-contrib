// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package events

import (
	"context"
	"fmt"

	"github.com/absmach/supermq-contrib/twins"
	"github.com/absmach/supermq/pkg/errors"
	"github.com/go-redis/redis/v8"
)

const (
	prefix = "twin"
)

var (
	// errRedisTwinSave indicates error while saving Twin in redis cache.
	errRedisTwinSave = errors.New("failed to save twin in redis cache")

	// errRedisTwinUpdate indicates error while saving Twin in redis cache.
	errRedisTwinUpdate = errors.New("failed to update twin in redis cache")

	// errRedisTwinIDs indicates error while getting Twin IDs from redis cache.
	errRedisTwinIDs = errors.New("failed to get twin id from redis cache")

	// errRedisTwinRemove indicates error while removing Twin from redis cache.
	errRedisTwinRemove = errors.New("failed to remove twin from redis cache")
)

var _ twins.TwinCache = (*twinCache)(nil)

type twinCache struct {
	client *redis.Client
}

// NewTwinCache returns redis twin cache implementation.
func NewTwinCache(client *redis.Client) twins.TwinCache {
	return &twinCache{
		client: client,
	}
}

func (tc *twinCache) Save(ctx context.Context, twin twins.Twin) error {
	return tc.save(ctx, twin)
}

func (tc *twinCache) Update(ctx context.Context, twin twins.Twin) error {
	if err := tc.remove(ctx, twin.ID); err != nil {
		return errors.Wrap(errRedisTwinUpdate, err)
	}
	if err := tc.save(ctx, twin); err != nil {
		return errors.Wrap(errRedisTwinUpdate, err)
	}
	return nil
}

func (tc *twinCache) SaveIDs(ctx context.Context, channel, subtopic string, ids []string) error {
	for _, id := range ids {
		if err := tc.client.SAdd(ctx, attrKey(channel, subtopic), id).Err(); err != nil {
			return errors.Wrap(errRedisTwinSave, err)
		}
		if err := tc.client.SAdd(ctx, twinKey(id), attrKey(channel, subtopic)).Err(); err != nil {
			return errors.Wrap(errRedisTwinSave, err)
		}
	}
	return nil
}

func (tc *twinCache) IDs(ctx context.Context, channel, subtopic string) ([]string, error) {
	ids, err := tc.client.SMembers(ctx, attrKey(channel, subtopic)).Result()
	if err != nil {
		return nil, errors.Wrap(errRedisTwinIDs, err)
	}
	idsWildcard, err := tc.client.SMembers(ctx, attrKey(channel, twins.SubtopicWildcard)).Result()
	if err != nil {
		return nil, errors.Wrap(errRedisTwinIDs, err)
	}
	ids = append(ids, idsWildcard...)
	return ids, nil
}

func (tc *twinCache) Remove(ctx context.Context, twinID string) error {
	return tc.remove(ctx, twinID)
}

func (tc *twinCache) save(ctx context.Context, twin twins.Twin) error {
	if len(twin.Definitions) < 1 {
		return nil
	}
	attributes := twin.Definitions[len(twin.Definitions)-1].Attributes
	for _, attr := range attributes {
		if err := tc.client.SAdd(ctx, attrKey(attr.Channel, attr.Subtopic), twin.ID).Err(); err != nil {
			return errors.Wrap(errRedisTwinSave, err)
		}
		if err := tc.client.SAdd(ctx, twinKey(twin.ID), attrKey(attr.Channel, attr.Subtopic)).Err(); err != nil {
			return errors.Wrap(errRedisTwinSave, err)
		}
	}
	return nil
}

func (tc *twinCache) remove(ctx context.Context, twinID string) error {
	twinKey := twinKey(twinID)
	attrKeys, err := tc.client.SMembers(ctx, twinKey).Result()
	if err != nil {
		return errors.Wrap(errRedisTwinRemove, err)
	}
	if err := tc.client.Del(ctx, twinKey).Err(); err != nil {
		return errors.Wrap(errRedisTwinRemove, err)
	}
	for _, attrKey := range attrKeys {
		if err := tc.client.SRem(ctx, attrKey, twinID).Err(); err != nil {
			return errors.Wrap(errRedisTwinRemove, err)
		}
	}
	return nil
}

func twinKey(twinID string) string {
	return fmt.Sprintf("%s:%s", prefix, twinID)
}

func attrKey(channel, subtopic string) string {
	return fmt.Sprintf("%s:%s-%s", prefix, channel, subtopic)
}
