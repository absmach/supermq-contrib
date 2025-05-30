// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package mongodb

import (
	"context"

	"github.com/absmach/supermq-contrib/twins"
	"github.com/absmach/supermq/pkg/errors"
	repoerr "github.com/absmach/supermq/pkg/errors/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	maxNameSize            = 1024
	twinsCollection string = "twins"
)

type twinRepository struct {
	db *mongo.Database
}

var _ twins.TwinRepository = (*twinRepository)(nil)

// NewTwinRepository instantiates a MongoDB implementation of twin repository.
func NewTwinRepository(db *mongo.Database) twins.TwinRepository {
	return &twinRepository{
		db: db,
	}
}

func (tr *twinRepository) Save(ctx context.Context, tw twins.Twin) (string, error) {
	if len(tw.Name) > maxNameSize {
		return "", errors.ErrMalformedEntity
	}

	coll := tr.db.Collection(twinsCollection)

	if _, err := coll.InsertOne(ctx, tw); err != nil {
		return "", errors.Wrap(repoerr.ErrCreateEntity, err)
	}

	return tw.ID, nil
}

func (tr *twinRepository) Update(ctx context.Context, tw twins.Twin) error {
	if len(tw.Name) > maxNameSize {
		return errors.ErrMalformedEntity
	}

	coll := tr.db.Collection(twinsCollection)

	filter := bson.M{"id": tw.ID}
	update := bson.M{"$set": tw}
	res, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount < 1 {
		return repoerr.ErrNotFound
	}

	return nil
}

func (tr *twinRepository) RetrieveByID(ctx context.Context, twinID string) (twins.Twin, error) {
	coll := tr.db.Collection(twinsCollection)
	var tw twins.Twin

	filter := bson.M{"id": twinID}
	if err := coll.FindOne(ctx, filter).Decode(&tw); err != nil {
		return tw, repoerr.ErrNotFound
	}

	return tw, nil
}

func (tr *twinRepository) RetrieveByAttribute(ctx context.Context, channel, subtopic string) ([]string, error) {
	coll := tr.db.Collection(twinsCollection)

	findOptions := options.Aggregate()
	prj1 := bson.M{
		"$project": bson.M{
			"definition": bson.M{
				"$arrayElemAt": []interface{}{"$definitions.attributes", -1},
			},
			"id":  true,
			"_id": 0,
		},
	}
	match := bson.M{
		"$match": bson.M{
			"definition.channel": channel,
			"$or": []interface{}{
				bson.M{"definition.subtopic": subtopic},
				bson.M{"definition.subtopic": twins.SubtopicWildcard},
			},
		},
	}
	prj2 := bson.M{
		"$project": bson.M{
			"id": true,
		},
	}

	cur, err := coll.Aggregate(ctx, []bson.M{prj1, match, prj2}, findOptions)
	if err != nil {
		return []string{}, errors.Wrap(repoerr.ErrViewEntity, err)
	}
	defer cur.Close(ctx)

	if err := cur.Err(); err != nil {
		return []string{}, nil
	}

	var ids []string
	for cur.Next(ctx) {
		var elem struct {
			ID string `json:"id"`
		}
		err := cur.Decode(&elem)
		if err != nil {
			return ids, nil
		}
		ids = append(ids, elem.ID)
	}

	return ids, nil
}

func (tr *twinRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata twins.Metadata) (twins.Page, error) {
	coll := tr.db.Collection(twinsCollection)

	findOptions := options.Find()
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	filter := bson.M{}

	if owner != "" {
		filter["owner"] = owner
	}
	if name != "" {
		filter["name"] = name
	}
	if len(metadata) > 0 {
		filter["metadata"] = metadata
	}
	cur, err := coll.Find(ctx, filter, findOptions)
	if err != nil {
		return twins.Page{}, errors.Wrap(repoerr.ErrViewEntity, err)
	}

	results, err := decodeTwins(ctx, cur)
	if err != nil {
		return twins.Page{}, errors.Wrap(repoerr.ErrViewEntity, err)
	}

	total, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return twins.Page{}, errors.Wrap(repoerr.ErrViewEntity, err)
	}

	return twins.Page{
		Twins: results,
		PageMetadata: twins.PageMetadata{
			Total:  uint64(total),
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (tr *twinRepository) Remove(ctx context.Context, twinID string) error {
	coll := tr.db.Collection(twinsCollection)

	filter := bson.M{"id": twinID}
	res, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(repoerr.ErrRemoveEntity, err)
	}

	if res.DeletedCount < 1 {
		return repoerr.ErrNotFound
	}

	return nil
}

func decodeTwins(ctx context.Context, cur *mongo.Cursor) ([]twins.Twin, error) {
	defer cur.Close(ctx)
	var results []twins.Twin
	for cur.Next(ctx) {
		var elem twins.Twin
		err := cur.Decode(&elem)
		if err != nil {
			return []twins.Twin{}, nil
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return []twins.Twin{}, nil
	}
	return results, nil
}
