package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"

	_ "github.com/lib/pq"
	tc "github.com/testcontainers/testcontainers-go"
)

func NewTestDB(t *testing.T) *sql.DB {
	containerReq := tc.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_USER":     "user",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		),
	}

	dbContainer, err := tc.GenericContainer(context.Background(),
		tc.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	hostPort, err := dbContainer.MappedPort(context.Background(), "5432")
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

func Called(fnCalls map[string]int, key string) {

	if fnCalls == nil {
		fnCalls = make(map[string]int)
	}

	value, exists := fnCalls[key]

	if !exists {
		fnCalls[key] = 1
		return
	}

	fnCalls[key] = value + 1
}
