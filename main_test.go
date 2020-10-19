package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
)

var (
	initDB   *sql.DB
	fixtures *testfixtures.Loader
	testDB   = createDB(testConnectionString())
	router   = gin.New()
)

// Add handlers to test router
func init() {
	router = getTasks(testDB)
	router = toggleTask(testDB)
	router = updateTask(testDB)
}

func TestToggleTask(t *testing.T) {
	prepareTestDatabase()
	// run server using httptest
	server := httptest.NewServer(router)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	getData(e, "Status", false)

	e.PUT("/toggle_task/1").
		Expect().
		Status(http.StatusOK)

	getData(e, "Status", true)
}

func TestUpdateTask(t *testing.T) {
	prepareTestDatabase()

	testTitle := map[string]interface{}{
		"Title": "song",
	}
	// run server using httptest
	server := httptest.NewServer(router)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	getData(e, "Title", "test_1")

	e.PUT("/task/1").WithJSON(testTitle).
		Expect().
		Status(http.StatusOK)

	getData(e, "Title", "song")
}

func testConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		"localhost", "5435", "user", "test_db", "pwd")
}

func TestMain(m *testing.M) {
	var err error

	initDB, err = sql.Open("postgres", testConnectionString())
	if err != nil {
		fmt.Println("can not open test DB")
	}

	fixtures, err = testfixtures.New(
		// You database connection
		testfixtures.Database(initDB),
		// Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Dialect("postgres"),
		// the directory containing the YAML files
		testfixtures.Directory("fixtures"),
		// refuse load fixtures if the db name (or db filename for SQLite) doesn't contains "test"
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	if err != nil {
		fmt.Println("can not create test Loader")
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		fmt.Println("can not prepare test Database")
	}
}

func getData(e *httpexpect.Expect, key string, value interface{}) *httpexpect.Object {
	return e.GET("/task").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("tasks").
		Object().Value("1").Object().ValueEqual(key, value)
}
