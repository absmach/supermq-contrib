// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package notifiers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/absmach/supermq-contrib/consumers/notifiers"
	"github.com/absmach/supermq-contrib/consumers/notifiers/mocks"
	"github.com/absmach/supermq/consumers"
	smqmocks "github.com/absmach/supermq/consumers/mocks"
	smqauthn "github.com/absmach/supermq/pkg/authn"
	"github.com/absmach/supermq/pkg/errors"
	repoerr "github.com/absmach/supermq/pkg/errors/repository"
	svcerr "github.com/absmach/supermq/pkg/errors/service"
	"github.com/absmach/supermq/pkg/messaging"
	"github.com/absmach/supermq/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	total        = 100
	exampleUser1 = "token1"
	exampleUser2 = "token2"
	validID      = "d4ebb847-5d0e-4e46-bdd9-b6aceaaa3a22"
)

func newService() (notifiers.Service, *mocks.SubscriptionsRepository) {
	repo := new(mocks.SubscriptionsRepository)
	notifier := new(smqmocks.Notifier)
	idp := uuid.NewMock()
	from := "exampleFrom"
	return notifiers.New(repo, idp, notifier, from), repo
}

func TestCreateSubscription(t *testing.T) {
	svc, repo := newService()

	cases := []struct {
		desc    string
		token   string
		session smqauthn.Session
		sub     notifiers.Subscription
		id      string
		err     error
		userID  string
	}{
		{
			desc:   "test success",
			token:  exampleUser1,
			sub:    notifiers.Subscription{Contact: exampleUser1, Topic: "valid.topic"},
			id:     uuid.Prefix + fmt.Sprintf("%012d", 1),
			err:    nil,
			userID: validID,
		},
		{
			desc:   "test already existing",
			token:  exampleUser1,
			sub:    notifiers.Subscription{Contact: exampleUser1, Topic: "valid.topic"},
			id:     "",
			err:    repoerr.ErrConflict,
			userID: validID,
		},
		{
			desc:  "test with empty token",
			token: "",
			sub:   notifiers.Subscription{Contact: exampleUser1, Topic: "valid.topic"},
			id:    "",
			err:   svcerr.ErrAuthentication,
		},
	}

	for _, tc := range cases {
		if tc.token == exampleUser1 {
			tc.session = smqauthn.Session{UserID: tc.userID}
		}
		repoCall1 := repo.On("Save", context.Background(), mock.Anything).Return(tc.id, tc.err)
		id, err := svc.CreateSubscription(context.Background(), tc.session, tc.sub)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.id, id, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.id, id))
		repoCall1.Unset()
	}
}

func TestViewSubscription(t *testing.T) {
	svc, repo := newService()
	sub := notifiers.Subscription{
		Contact: exampleUser1,
		Topic:   "valid.topic",
		ID:      generateUUID(t),
		OwnerID: validID,
	}

	cases := []struct {
		desc    string
		token   string
		session smqauthn.Session
		id      string
		sub     notifiers.Subscription
		err     error
		userID  string
	}{
		{
			desc:   "test success",
			token:  exampleUser1,
			id:     validID,
			sub:    sub,
			err:    nil,
			userID: validID,
		},
		{
			desc:   "test not existing",
			token:  exampleUser1,
			id:     "not_exist",
			sub:    notifiers.Subscription{},
			err:    svcerr.ErrNotFound,
			userID: validID,
		},
		{
			desc:  "test with empty token",
			token: "",
			id:    validID,
			sub:   notifiers.Subscription{},
			err:   svcerr.ErrAuthentication,
		},
	}

	for _, tc := range cases {
		if tc.token == exampleUser1 {
			tc.session = smqauthn.Session{UserID: tc.userID}
		}
		repoCall1 := repo.On("Retrieve", context.Background(), tc.id).Return(tc.sub, tc.err)
		sub, err := svc.ViewSubscription(context.Background(), tc.session, tc.id)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.sub, sub, fmt.Sprintf("%s: expected %v got %v\n", tc.desc, tc.sub, sub))
		repoCall1.Unset()
	}
}

func TestListSubscriptions(t *testing.T) {
	svc, repo := newService()
	sub := notifiers.Subscription{Contact: exampleUser1, OwnerID: exampleUser1}
	topic := "topic.subtopic"
	var subs []notifiers.Subscription
	for i := 0; i < total; i++ {
		tmp := sub
		if i%2 == 0 {
			tmp.Contact = exampleUser2
			tmp.OwnerID = exampleUser2
		}
		tmp.Topic = fmt.Sprintf("%s.%d", topic, i)
		tmp.ID = generateUUID(t)
		tmp.OwnerID = validID
		subs = append(subs, tmp)
	}

	var offsetSubs []notifiers.Subscription
	for i := 20; i < 40; i += 2 {
		offsetSubs = append(offsetSubs, subs[i])
	}

	cases := []struct {
		desc            string
		token           string
		session         smqauthn.Session
		pageMeta        notifiers.PageMetadata
		page            notifiers.Page
		err             error
		authenticateErr error
		userID          string
	}{
		{
			desc:  "test success",
			token: exampleUser1,
			pageMeta: notifiers.PageMetadata{
				Offset: 0,
				Limit:  3,
			},
			err: nil,
			page: notifiers.Page{
				PageMetadata: notifiers.PageMetadata{
					Offset: 0,
					Limit:  3,
				},
				Subscriptions: subs[:3],
				Total:         total,
			},
			authenticateErr: nil,
			userID:          validID,
		},
		{
			desc:  "test not existing",
			token: exampleUser1,
			pageMeta: notifiers.PageMetadata{
				Limit:   10,
				Contact: "empty@example.com",
			},
			page:            notifiers.Page{},
			err:             svcerr.ErrNotFound,
			authenticateErr: nil,
			userID:          validID,
		},
		{
			desc:  "test with empty token",
			token: "",
			pageMeta: notifiers.PageMetadata{
				Offset: 2,
				Limit:  12,
				Topic:  "topic.subtopic.13",
			},
			page:            notifiers.Page{},
			err:             svcerr.ErrAuthentication,
			authenticateErr: svcerr.ErrAuthentication,
		},
		{
			desc:  "test with topic",
			token: exampleUser1,
			pageMeta: notifiers.PageMetadata{
				Limit: 10,
				Topic: fmt.Sprintf("%s.%d", topic, 4),
			},
			page: notifiers.Page{
				PageMetadata: notifiers.PageMetadata{
					Limit: 10,
					Topic: fmt.Sprintf("%s.%d", topic, 4),
				},
				Subscriptions: subs[4:5],
				Total:         1,
			},
			err:             nil,
			authenticateErr: nil,
			userID:          validID,
		},
		{
			desc:  "test with contact and offset",
			token: exampleUser1,
			pageMeta: notifiers.PageMetadata{
				Offset:  10,
				Limit:   10,
				Contact: exampleUser2,
			},
			page: notifiers.Page{
				PageMetadata: notifiers.PageMetadata{
					Offset:  10,
					Limit:   10,
					Contact: exampleUser2,
				},
				Subscriptions: offsetSubs,
				Total:         uint(total / 2),
			},
			err:             nil,
			authenticateErr: nil,
			userID:          validID,
		},
	}

	for _, tc := range cases {
		if tc.token == exampleUser1 {
			tc.session = smqauthn.Session{UserID: tc.userID}
		}
		repoCall1 := repo.On("RetrieveAll", context.Background(), tc.pageMeta).Return(tc.page, tc.err)
		page, err := svc.ListSubscriptions(context.Background(), tc.session, tc.pageMeta)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.page, page, fmt.Sprintf("%s: got unexpected page\n", tc.desc))
		repoCall1.Unset()
	}
}

func TestRemoveSubscription(t *testing.T) {
	svc, repo := newService()
	sub := notifiers.Subscription{
		ID: generateUUID(t),
	}

	cases := []struct {
		desc    string
		token   string
		session smqauthn.Session
		id      string
		err     error
		userID  string
	}{
		{
			desc:   "test success",
			token:  exampleUser1,
			id:     sub.ID,
			err:    nil,
			userID: validID,
		},
		{
			desc:   "test not existing",
			token:  exampleUser1,
			id:     "not_exist",
			err:    svcerr.ErrNotFound,
			userID: validID,
		},
		{
			desc:  "test with empty token",
			token: "",
			id:    sub.ID,
			err:   svcerr.ErrAuthentication,
		},
	}

	for _, tc := range cases {
		if tc.token == exampleUser1 {
			tc.session = smqauthn.Session{UserID: tc.userID}
		}
		repoCall1 := repo.On("Remove", context.Background(), tc.id).Return(tc.err)
		err := svc.RemoveSubscription(context.Background(), tc.session, tc.id)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall1.Unset()
	}
}

func TestConsume(t *testing.T) {
	svc, repo := newService()
	msg := messaging.Message{
		Channel:  "topic",
		Subtopic: "subtopic",
	}
	errMsg := messaging.Message{
		Channel:  "topic",
		Subtopic: "subtopic-2",
	}

	cases := []struct {
		desc string
		msg  *messaging.Message
		err  error
	}{
		{
			desc: "test success",
			msg:  &msg,
			err:  nil,
		},
		{
			desc: "test fail",
			msg:  &errMsg,
			err:  consumers.ErrNotify,
		},
	}

	for _, tc := range cases {
		repoCall := repo.On("RetrieveAll", context.TODO(), mock.Anything).Return(notifiers.Page{}, tc.err)
		err := svc.ConsumeBlocking(context.TODO(), tc.msg)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func generateUUID(t *testing.T) string {
	idProvider := uuid.New()
	ulid, err := idProvider.ID()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	return ulid
}
