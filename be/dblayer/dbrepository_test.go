package dblayer

import (
	"database/sql"
	"log"

	"testing"

	_ "github.com/go-sql-driver/mysql"
)

/*
1. Connect to "root:mysecret@tcp(localhost:3306)/rproject"
2. create a new DBEntity of type "users"
3. set the "login" column to "u"
4. create dbrepository with the connection
5. search for the user with login "u"
6. print the user
*/
func TestSearchUserByLogin(t *testing.T) {
	// Step 1: Connect to the database
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   "rprj",
	}
	factory := NewDBEFactory(true)
	user := NewDBUser()
	factory.Register(&user.DBEntity)
	dbConnection, err := sql.Open("mysql", "root:mysecret@tcp(localhost:3306)/rproject")
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}
	defer dbConnection.Close()

	// Step 2: Create a new DBEntity of type "users"
	userEntity := factory.GetInstanceByTableName("users")
	if userEntity == nil {
		t.Fatal("Failed to get DBEntity for 'users'")
	}

	// Step 3: Set the "login" column to ...
	userEntity.SetValue("login", "a")

	// Step 4: Create DBRepository with the connection
	repo := NewDBRepository(dbContext, factory, dbConnection)

	repo.Verbose = true

	// Step 5: Search for the user with specified login
	results, err := repo.Search(userEntity, true, true, "login")
	if err != nil {
		t.Fatal("Failed to search for user:", err)
	}
	if len(results) == 0 {
		t.Fatal("No user found with login 'u'")
	}

	// Step 6: Print the results
	for _, user := range results {
		log.Printf("- %s\t%s\t%s\n", user.GetValue("id"), user.GetValue("login"), user.GetValue("fullname"))
	}
}
