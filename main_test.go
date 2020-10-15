package main

import (
	"database/sql"
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	initDb   *sql.DB
	fixtures *testfixtures.Loader
	testDb   = createDB(testConnectionString())
	router   = gin.New()
)

// Add handlers to test_docker_compose router
func init() {
	router = getTasks(testDb)
	router = toggleTask(testDb)
	router = updateTask(testDb)
}

func TestMain(m *testing.M) {
	var err error

	initDb, err = sql.Open("postgres", testConnectionString())
	if err != nil {
		fmt.Println("can not open test_docker_compose DB")
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(initDb),                 // You database connection
		testfixtures.Dialect("postgres"),              // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory("fixtures"),            // the directory containing the YAML files
		testfixtures.DangerousSkipTestDatabaseCheck(), // will refuse to load fixtures if the database name (or database filename for SQLite) doesn't contains "test_docker_compose"
	)
	if err != nil {
		fmt.Println("can not create test_docker_compose Loader")
	}

	os.Exit(m.Run())

}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		fmt.Println("can not prepare test_docker_compose Database")
	}
}

func TestToggleTask(t *testing.T) {
	prepareTestDatabase()
	// run server using httptest
	server := httptest.NewServer(router)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	e.GET("/task").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("tasks").
		Object().Value("1").Object().ValueEqual("Status", false)

	e.PUT("/toggle_task/1").
		Expect().
		Status(http.StatusOK)

	e.GET("/task").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("tasks").
		Object().Value("1").Object().ValueEqual("Status", true)

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

	e.GET("/task").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("tasks").
		Object().Value("1").Object().ValueEqual("Title", "test_1")

	e.PUT("/task/1").WithJSON(testTitle).
		Expect().
		Status(http.StatusOK)

	e.GET("/task").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("tasks").
		Object().Value("1").Object().ValueEqual("Title", "song")
}

func testConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", "localhost", "5435", "user", "test_db", "pwd")
}
