// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/absmach/supermq-contrib/consumers/notifiers"
	"github.com/absmach/supermq-contrib/consumers/notifiers/api"
	"github.com/absmach/supermq-contrib/consumers/notifiers/mocks"
	apiutil "github.com/absmach/supermq/api/http/util"
	smqlog "github.com/absmach/supermq/logger"
	smqauthn "github.com/absmach/supermq/pkg/authn"
	authnmocks "github.com/absmach/supermq/pkg/authn/mocks"
	"github.com/absmach/supermq/pkg/errors"
	svcerr "github.com/absmach/supermq/pkg/errors/service"
	"github.com/absmach/supermq/pkg/uuid"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	contentType  = "application/json"
	email        = "user@example.com"
	contact1     = "email1@example.com"
	contact2     = "email2@example.com"
	token        = "token"
	invalidToken = "invalid"
	topic        = "topic"
	instanceID   = "5de9b29a-feb9-11ed-be56-0242ac120002"
	validID      = "d4ebb847-5d0e-4e46-bdd9-b6aceaaa3a22"
)

var (
	notFoundRes   = toJSON(apiutil.ErrorRes{Msg: svcerr.ErrNotFound.Error()})
	unauthRes     = toJSON(apiutil.ErrorRes{Msg: svcerr.ErrAuthentication.Error()})
	invalidRes    = toJSON(apiutil.ErrorRes{Err: apiutil.ErrInvalidQueryParams.Error(), Msg: apiutil.ErrValidation.Error()})
	missingTokRes = toJSON(apiutil.ErrorRes{Err: apiutil.ErrBearerToken.Error(), Msg: apiutil.ErrValidation.Error()})
)

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	contentType string
	token       string
	body        io.Reader
}

func (tr testRequest) make() (*http.Response, error) {
	req, err := http.NewRequest(tr.method, tr.url, tr.body)
	if err != nil {
		return nil, err
	}

	if tr.token != "" {
		req.Header.Set("Authorization", apiutil.BearerPrefix+tr.token)
	}

	if tr.contentType != "" {
		req.Header.Set("Content-Type", tr.contentType)
	}

	req.Header.Set("Referer", "http://localhost")

	return tr.client.Do(req)
}

func newServer() (*httptest.Server, *mocks.Service, *authnmocks.Authentication) {
	logger := smqlog.NewMock()
	svc := new(mocks.Service)
	mux := chi.NewRouter()
	authn := new(authnmocks.Authentication)
	api.MakeHandler(svc, mux, logger, instanceID, authn)
	return httptest.NewServer(mux), svc, authn
}

func toJSON(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func TestCreate(t *testing.T) {
	ss, svc, authn := newServer()
	defer ss.Close()

	sub := notifiers.Subscription{
		Topic:   topic,
		Contact: contact1,
	}

	emptyTopic := notifiers.Subscription{Contact: contact1}
	emptyContact := notifiers.Subscription{Topic: "topic123"}

	cases := []struct {
		desc        string
		req         notifiers.Subscription
		contentType string
		token       string
		authnRes    smqauthn.Session
		authnErr    error
		status      int
		err         error
	}{
		{
			desc:        "add successfully",
			req:         sub,
			contentType: contentType,
			token:       token,
			authnRes:    smqauthn.Session{UserID: validID},
			status:      http.StatusCreated,
			err:         nil,
		},
		{
			desc:        "add an existing subscription",
			req:         sub,
			contentType: contentType,
			token:       token,
			authnRes:    smqauthn.Session{UserID: validID},
			status:      http.StatusConflict,
			err:         svcerr.ErrConflict,
		},
		{
			desc:        "add with empty topic",
			req:         emptyTopic,
			contentType: contentType,
			token:       token,
			authnRes:    smqauthn.Session{UserID: validID},
			status:      http.StatusBadRequest,
			err:         apiutil.ErrInvalidTopic,
		},
		{
			desc:        "add with empty contact",
			req:         emptyContact,
			contentType: contentType,
			token:       token,
			authnRes:    smqauthn.Session{UserID: validID},
			status:      http.StatusBadRequest,
			err:         svcerr.ErrMalformedEntity,
		},
		{
			desc:        "add with invalid auth token",
			req:         sub,
			contentType: contentType,
			token:       invalidToken,
			status:      http.StatusUnauthorized,
			authnErr:    svcerr.ErrAuthentication,
			err:         svcerr.ErrAuthentication,
		},
		{
			desc:        "add with empty auth token",
			req:         sub,
			contentType: contentType,
			token:       "",
			status:      http.StatusUnauthorized,
			err:         svcerr.ErrAuthentication,
		},

		{
			desc:        "add with invalid content type",
			req:         sub,
			contentType: "application/xml",
			token:       token,
			authnRes:    smqauthn.Session{UserID: validID},
			status:      http.StatusUnsupportedMediaType,
			err:         apiutil.ErrUnsupportedContentType,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			data := toJSON(tc.req)
			req := testRequest{
				client:      ss.Client(),
				method:      http.MethodPost,
				url:         fmt.Sprintf("%s/subscriptions", ss.URL),
				contentType: tc.contentType,
				token:       tc.token,
				body:        strings.NewReader(data),
			}

			authCall := authn.On("Authenticate", mock.Anything, tc.token).Return(tc.authnRes, tc.authnErr)
			svcCall := svc.On("CreateSubscription", mock.Anything, tc.authnRes, tc.req).Return(sub, tc.err)
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))

			var errRes respBody
			err = json.NewDecoder(res.Body).Decode(&errRes)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
			if errRes.Err != "" || errRes.Message != "" {
				err = errors.Wrap(errors.New(errRes.Err), errors.New(errRes.Message))
			}

			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
			authCall.Unset()
			svcCall.Unset()
		})
	}
}

func TestView(t *testing.T) {
	ss, svc, authn := newServer()
	defer ss.Close()

	sub := notifiers.Subscription{
		Topic:   topic,
		Contact: contact1,
		ID:      generateUUID(t),
		OwnerID: validID,
	}

	cases := []struct {
		desc     string
		id       string
		token    string
		authnRes smqauthn.Session
		authnErr error
		status   int
		err      error
		Sub      notifiers.Subscription
	}{
		{
			desc:     "view successfully",
			id:       sub.ID,
			token:    token,
			authnRes: smqauthn.Session{UserID: validID},
			status:   http.StatusOK,
			err:      nil,
			Sub:      sub,
		},
		{
			desc:     "view not existing",
			id:       "not existing",
			token:    token,
			authnRes: smqauthn.Session{UserID: validID},
			status:   http.StatusNotFound,
			err:      svcerr.ErrNotFound,
		},
		{
			desc:     "view with invalid auth token",
			id:       sub.ID,
			token:    invalidToken,
			status:   http.StatusUnauthorized,
			authnErr: svcerr.ErrAuthentication,
			err:      svcerr.ErrAuthentication,
		},
		{
			desc:   "view with empty auth token",
			id:     sub.ID,
			token:  "",
			status: http.StatusUnauthorized,
			err:    apiutil.ErrBearerToken,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			req := testRequest{
				client: ss.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/subscriptions/%s", ss.URL, tc.id),
				token:  tc.token,
			}

			authCall := authn.On("Authenticate", mock.Anything, tc.token).Return(tc.authnRes, tc.authnErr)
			svcCall := svc.On("ViewSubscription", mock.Anything, tc.authnRes, tc.id).Return(tc.Sub, tc.err)

			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected request error %s", tc.desc, err))
			var errRes respBody
			err = json.NewDecoder(res.Body).Decode(&errRes)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
			if errRes.Err != "" || errRes.Message != "" {
				err = errors.Wrap(errors.New(errRes.Err), errors.New(errRes.Message))
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
			authCall.Unset()
			svcCall.Unset()
		})
	}
}

func TestList(t *testing.T) {
	ss, svc, authn := newServer()
	defer ss.Close()

	const numSubs = 100
	var subs []subRes
	var sub notifiers.Subscription

	for i := 0; i < numSubs; i++ {
		sub = notifiers.Subscription{
			Topic:   fmt.Sprintf("topic.subtopic.%d", i),
			Contact: contact1,
			ID:      generateUUID(t),
		}
		if i%2 == 0 {
			sub.Contact = contact2
		}
		sr := subRes{
			ID:      sub.ID,
			OwnerID: validID,
			Contact: sub.Contact,
			Topic:   sub.Topic,
		}
		subs = append(subs, sr)
	}

	var contact2Subs []subRes
	for i := 20; i < 40; i += 2 {
		contact2Subs = append(contact2Subs, subs[i])
	}

	cases := []struct {
		desc     string
		query    map[string]string
		token    string
		authnRes smqauthn.Session
		authnErr error
		status   int

		err  error
		page notifiers.Page
	}{
		{
			desc: "list default limit",
			query: map[string]string{
				"offset": "5",
			},
			token:    token,
			authnRes: smqauthn.Session{UserID: validID},
			status:   http.StatusOK,

			err: nil,
			page: notifiers.Page{
				PageMetadata: notifiers.PageMetadata{
					Offset: 5,
					Limit:  20,
				},
				Total:         numSubs,
				Subscriptions: subscriptionsSlice(subs, 5, 25),
			},
		},
		{
			desc: "list not existing",
			query: map[string]string{
				"topic": "not-found-topic",
			},
			token:    token,
			authnRes: smqauthn.Session{UserID: validID},
			status:   http.StatusNotFound,
			err:      svcerr.ErrNotFound,
		},
		{
			desc: "list one with topic",
			query: map[string]string{
				"topic": "topic.subtopic.10",
			},
			token:    token,
			authnRes: smqauthn.Session{UserID: validID},
			status:   http.StatusOK,

			err: nil,
			page: notifiers.Page{
				PageMetadata: notifiers.PageMetadata{
					Offset: 0,
					Limit:  20,
				},
				Total:         1,
				Subscriptions: subscriptionsSlice(subs, 10, 11),
			},
		},
		{
			desc: "list with contact",
			query: map[string]string{
				"contact": contact2,
				"offset":  "10",
				"limit":   "10",
			},
			token:    token,
			authnRes: smqauthn.Session{UserID: validID},
			status:   http.StatusOK,
			err:      nil,
			page: notifiers.Page{
				PageMetadata: notifiers.PageMetadata{
					Offset: 10,
					Limit:  10,
				},
				Total:         50,
				Subscriptions: subscriptionsSlice(contact2Subs, 0, 10),
			},
		},
		{
			desc: "list with invalid query",
			query: map[string]string{
				"offset": "two",
			},
			token:    token,
			authnRes: smqauthn.Session{UserID: validID},
			status:   http.StatusBadRequest,
			err:      apiutil.ErrValidation,
		},
		{
			desc:     "list with invalid auth token",
			token:    invalidToken,
			status:   http.StatusUnauthorized,
			authnErr: svcerr.ErrAuthentication,
			err:      svcerr.ErrAuthentication,
		},
		{
			desc:   "list with empty auth token",
			token:  "",
			status: http.StatusUnauthorized,
			err:    apiutil.ErrBearerToken,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			req := testRequest{
				client: ss.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/subscriptions%s", ss.URL, makeQuery(tc.query)),
				token:  tc.token,
			}

			authCall := authn.On("Authenticate", mock.Anything, tc.token).Return(tc.authnRes, tc.authnErr)
			svcCall := svc.On("ListSubscriptions", mock.Anything, tc.authnRes, mock.Anything).Return(tc.page, tc.err)

			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))

			var bodyRes respBody
			err = json.NewDecoder(res.Body).Decode(&bodyRes)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
			if bodyRes.Err != "" || bodyRes.Message != "" {
				err = errors.Wrap(errors.New(bodyRes.Err), errors.New(bodyRes.Message))
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))

			svcCall.Unset()
			authCall.Unset()
		})
	}
}

func TestRemove(t *testing.T) {
	ss, svc, authn := newServer()
	defer ss.Close()
	id := generateUUID(t)

	cases := []struct {
		desc     string
		id       string
		token    string
		authnRes smqauthn.Session
		authnErr error
		status   int
		res      string
		err      error
	}{
		{
			desc:     "remove successfully",
			id:       id,
			token:    token,
			authnRes: smqauthn.Session{UserID: validID},
			status:   http.StatusNoContent,
			err:      nil,
		},
		{
			desc:     "remove not existing",
			id:       "not existing",
			token:    token,
			authnRes: smqauthn.Session{UserID: validID},
			status:   http.StatusNotFound,
			err:      svcerr.ErrNotFound,
		},
		{
			desc:     "remove empty id",
			id:       "",
			token:    token,
			authnRes: smqauthn.Session{UserID: validID},
			status:   http.StatusBadRequest,
			err:      svcerr.ErrMalformedEntity,
		},
		{
			desc:     "view with invalid auth token",
			id:       id,
			token:    invalidToken,
			status:   http.StatusUnauthorized,
			res:      unauthRes,
			authnErr: svcerr.ErrAuthentication,
			err:      svcerr.ErrAuthentication,
		},
		{
			desc:   "view with empty auth token",
			id:     id,
			token:  "",
			status: http.StatusUnauthorized,
			res:    missingTokRes,
			err:    svcerr.ErrAuthentication,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			req := testRequest{
				client: ss.Client(),
				method: http.MethodDelete,
				url:    fmt.Sprintf("%s/subscriptions/%s", ss.URL, tc.id),
				token:  tc.token,
			}

			authCall := authn.On("Authenticate", mock.Anything, tc.token).Return(tc.authnRes, tc.authnErr)
			svcCall := svc.On("RemoveSubscription", mock.Anything, tc.authnRes, tc.id).Return(tc.err)

			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))

			authCall.Unset()
			svcCall.Unset()
		})
	}
}

func makeQuery(m map[string]string) string {
	var ret string
	for k, v := range m {
		ret += fmt.Sprintf("&%s=%s", k, v)
	}
	if ret != "" {
		return fmt.Sprintf("?%s", ret[1:])
	}
	return ""
}

type subRes struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Contact string `json:"contact"`
	Topic   string `json:"topic"`
}
type page struct {
	Offset        uint     `json:"offset"`
	Limit         int      `json:"limit"`
	Total         uint     `json:"total,omitempty"`
	Subscriptions []subRes `json:"subscriptions,omitempty"`
}

func subscriptionsSlice(subs []subRes, start, end int) []notifiers.Subscription {
	var res []notifiers.Subscription
	for i := start; i < end; i++ {
		sub := subs[i]
		res = append(res, notifiers.Subscription{
			ID:      sub.ID,
			OwnerID: sub.OwnerID,
			Contact: sub.Contact,
			Topic:   sub.Topic,
		})
	}
	return res
}

func generateUUID(t *testing.T) string {
	idProvider := uuid.New()
	ulid, err := idProvider.ID()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	return ulid
}

type respBody struct {
	Err     string `json:"error"`
	Message string `json:"message"`
}
