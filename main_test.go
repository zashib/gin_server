package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
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

// Add handlers to router
func init() {
	router = getTasks(testDb)
	router = toggleTask(testDb)
	router = updateTask(testDb)
}

func TestMain(m *testing.M) {
	var err error

	initDb, err = sql.Open("postgres", testConnectionString())
	if err != nil {
		fmt.Println("can not open test DB")
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(initDb),                 // You database connection
		testfixtures.Dialect("postgres"),              // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory("fixtures"),            // the directory containing the YAML files
		testfixtures.DangerousSkipTestDatabaseCheck(), // will refuse to load fixtures if the database name (or database filename for SQLite) doesn't contains "test"
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

func TestToggleTask(t *testing.T) {
	prepareTestDatabase()
	id := "1"
	//tasks := make(map[int64]TasksCategory)
	initResponse := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/task", nil)
	router.ServeHTTP(initResponse, req)
	assert.Equal(t, 200, initResponse.Code)
	fmt.Println(initResponse.Body)

	req, _ = http.NewRequest("PUT", "/toggle_task/"+id, nil)
	router.ServeHTTP(initResponse, req)
	assert.Equal(t, 200, initResponse.Code)

	endResponse := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/task", nil)
	router.ServeHTTP(endResponse, req)
	assert.Equal(t, 200, endResponse.Code)
	fmt.Println(endResponse.Body)
	//assert.Equal(t, true,  endResponse.Body)
}

func TestUpdateTask(t *testing.T) {
	prepareTestDatabase()
	id := "1"
	requestBody, _ := json.Marshal(map[string]string{
		"title": "song",
	})

	initResponse := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/task", nil)
	router.ServeHTTP(initResponse, req)
	assert.Equal(t, 200, initResponse.Code)
	fmt.Println(initResponse.Body)

	req, _ = http.NewRequest("PUT", "/task/"+id, bytes.NewBuffer(requestBody))
	router.ServeHTTP(initResponse, req)
	assert.Equal(t, 200, initResponse.Code)

	endResponse := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/task", nil)
	router.ServeHTTP(endResponse, req)
	assert.Equal(t, 200, endResponse.Code)
	fmt.Println(endResponse.Body)
}

func testConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", "localhost", "5435", "user", "test_db", "pwd")
}
