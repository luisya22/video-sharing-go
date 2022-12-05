package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	tc "github.com/testcontainers/testcontainers-go"
)

func NewTestDB(t *testing.T) *sql.DB {
	postgresPort := nat.Port("5432/tcp")
	postgres, err := tc.GenericContainer(context.Background(),
		tc.GenericContainerRequest{
			ContainerRequest: tc.ContainerRequest{
				Image:        "postgres",
				ExposedPorts: []string{postgresPort.Port()},
				Env: map[string]string{
					"POSTGRES_PASSWORD": "pass",
					"POSTGRES_USER":     "user",
				},
				WaitingFor: wait.ForAll(
					wait.ForLog("database system is ready to accept connections"),
					wait.ForListeningPort(postgresPort),
				),
			},
			Started: true,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	hostPort, err := postgres.MappedPort(context.Background(), postgresPort)
	if err != nil {
		t.Fatal(err)
	}

	postgresURLTemplate := "postgres://user:pass@localhost:%s?sslmode=disable"
	postgresURL := fmt.Sprintf(postgresURLTemplate, hostPort.Port())

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("../tests/testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("../tests/testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	return db
}
