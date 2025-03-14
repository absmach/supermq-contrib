// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package notifiers

import (
	"context"
)

// Subscription represents a user Subscription.
type Subscription struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Contact string `json:"contact"`
	Topic   string `json:"topic"`
}

// Page represents page metadata with content.
type Page struct {
	PageMetadata
	Total         uint           `json:"total"`
	Subscriptions []Subscription `json:"subscriptions"`
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Offset uint `json:"offset"`
	// Limit values less than 0 indicate no limit.
	Limit   int    `json:"limit"`
	Topic   string `json:"topic,omitempty"`
	Contact string `json:"contact,omitempty"`
}

// SubscriptionsRepository specifies a Subscription persistence API.
//
//go:generate mockery --name SubscriptionsRepository --output=./mocks --filename repository.go --quiet --note "Copyright (c) Abstract Machines"
type SubscriptionsRepository interface {
	// Save persists a subscription. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, sub Subscription) (Subscription, error)

	// Retrieve retrieves the subscription for the given id.
	Retrieve(ctx context.Context, id string) (Subscription, error)

	// RetrieveAll retrieves all the subscriptions for the given page metadata.
	RetrieveAll(ctx context.Context, pm PageMetadata) (Page, error)

	// Remove removes the subscription for the given ID.
	Remove(ctx context.Context, id string) error
}
