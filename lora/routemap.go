// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package lora

import "context"

// RouteMapRepository store route map between Lora App Server and SupeMQ.
type RouteMapRepository interface {
	// Save stores/routes pair lora application topic & supermq channel.
	Save(context.Context, string, string) error

	// Channel returns supermq channel for given lora application.
	Get(context.Context, string) (string, error)

	// Removes mapping from cache.
	Remove(context.Context, string) error
}
