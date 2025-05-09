// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/absmach/supermq"
	"github.com/absmach/supermq-contrib/opcua"
)

var _ supermq.Response = (*browseRes)(nil)

type browseRes struct {
	Nodes []opcua.BrowsedNode `json:"nodes"`
}

func (res browseRes) Code() int {
	return http.StatusOK
}

func (res browseRes) Headers() map[string]string {
	return map[string]string{}
}

func (res browseRes) Empty() bool {
	return false
}
