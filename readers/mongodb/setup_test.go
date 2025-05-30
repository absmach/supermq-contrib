// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package mongodb_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	smqlog "github.com/absmach/supermq/logger"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testLog, _ = smqlog.New(os.Stdout, "info")

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		testLog.Error(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "7.0.5",
		Env: []string{
			"MONGO_INITDB_DATABASE=test",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	port = container.GetPort("27017/tcp")
	addr = fmt.Sprintf("mongodb://localhost:%s", port)

	if err := pool.Retry(func() error {
		_, err := mongo.Connect(context.Background(), options.Client().ApplyURI(addr))
		return err
	}); err != nil {
		testLog.Error(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	code := m.Run()

	if err := pool.Purge(container); err != nil {
		testLog.Error(fmt.Sprintf("Could not purge container: %s", err))
	}

	os.Exit(code)
}
