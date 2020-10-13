package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
	"os"
	"testing"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Loader
	testDb   = createDB(testConnectionString())
	router   = gin.New()
)

func TestMain(m *testing.M) {
	var err error

	db, err = sql.Open("postgres", testConnectionString())
	if err != nil {
		fmt.Println("can not open test DB")
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(db),                     // You database connection
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

func TestX(t *testing.T) {
	prepareTestDatabase()
	router.GET("/task", func(c *gin.Context) {
		type TasksCategory struct {
			Title   string
			Content string
			Status  bool
			Name    string
		}
		var tasks []TasksCategory
		findErr := testDb.Table("tasks").Join("LEFT", "category", "category.id = tasks.category_id").
			Find(&tasks)
		if findErr != nil {
			fmt.Println(findErr)
		}
		c.JSON(200, gin.H{
			"tasks": tasks,
		})
	})

}

//// Rid of debug output
//func init() {
//	gin.SetMode(gin.TestMode)
//}

func testConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", "localhost", "5435", "user", "test_db", "pwd")
}
