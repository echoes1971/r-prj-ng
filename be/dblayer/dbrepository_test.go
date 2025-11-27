package dblayer

import (
	"fmt"
	"log"
	"sync"
	"time"

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

	// Step 2: Create a new DBEntity of type "users"
	userEntity := Factory.GetInstanceByTableName("users")
	if userEntity == nil {
		t.Fatal("Failed to get DBEntity for 'users'")
	}

	// Step 3: Set the "login" column to ...
	userEntity.SetValue("login", "a")

	// Step 4: Create DBRepository with the connection
	repo := NewDBRepository(dbContext, Factory, DbConnection)

	repo.Verbose = false

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

func TestInsertUser(t *testing.T) {
	// Setup
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   "rprj",
	}

	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false

	// Create a new user entity
	newUser := Factory.GetInstanceByTableName("users")
	if newUser == nil {
		t.Fatal("Failed to get DBEntity for 'users'")
	}

	// Set values (id must be set manually since it's not auto-increment)
	login := "testuser_" + Random4digits()
	newUser.SetValue("login", login)
	newUser.SetValue("pwd", "testpassword")
	newUser.SetValue("fullname", "Test User")
	// newUser.SetValue("group_id", "-3") // -3 is the default "users" group

	// Insert the user (transaction is handled internally)
	newUser, newerr := repo.Insert(newUser)
	if newerr != nil {
		t.Fatal("Failed to insert user:", newerr)
	}

	// Verify insertion by searching
	searchEntity := Factory.GetInstanceByTableName("users")
	searchEntity.SetValue("login", login)
	results, err := repo.Search(searchEntity, false, true, "")
	if err != nil {
		t.Fatal("Failed to search for inserted user:", err)
	}
	if len(results) == 0 {
		t.Fatal("Inserted user not found")
	}
	// Verify the group_id is set correctly
	if results[0].GetValue("group_id") == "" {
		t.Fatal("Inserted user's group_id is not set")
	}
	// Verify the group exists
	groupEntity := Factory.GetInstanceByTableName("groups")
	groupEntity.SetValue("id", results[0].GetValue("group_id"))
	groupResults, err := repo.Search(groupEntity, false, true, "")
	if err != nil {
		t.Fatal("Failed to search for user's group:", err)
	}
	if len(groupResults) == 0 {
		t.Fatal("User's group not found")
	}

	// Print success message

	log.Printf("Successfully inserted and found user: %s", results[0].GetValue("login"))

	// Cleanup: delete the test user
	_, err = repo.Delete(newUser)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test user: %v", err)
	}
}

func TestConcurrentMayhem(t *testing.T) {
	// Setup
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   "rprj",
	}

	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false

	var wgCreate sync.WaitGroup
	var wgDelete sync.WaitGroup

	var wgMayhem sync.WaitGroup

	// Rate limiter to control the number of concurrent goroutines
	// Buffer size should be much smaller than MaxOpenConns to leave room for nested operations
	// Each Insert creates: user + group + usergroup (3 operations in 1 transaction)
	// Then launches a Delete goroutine (which does: delete usergroup + delete group + delete user)
	max_concurrent := 10 // Small buffer - ticker will refill gradually
	burstyLimiter := make(chan time.Time, max_concurrent)
	for range max_concurrent {
		burstyLimiter <- time.Now()
	}
	go func() {
		for range time.Tick(25 * time.Millisecond) {
			burstyLimiter <- time.Now()
		}
	}()

	concurrent_routines := 500
	userPrefix := "mayhem_" + Random4digits()

	for i := range concurrent_routines {
		<-burstyLimiter // Acquire a token

		wgCreate.Add(1)
		go func(index int) {
			defer wgCreate.Done()
			// Create a new user entity
			newUser := Factory.GetInstanceByTableName("users")
			if newUser == nil {
				t.Log("Failed to get DBEntity for 'users'")
				return
			}

			// Set values (id must be set manually since it's not auto-increment)
			login := fmt.Sprintf("%s_%04d", userPrefix, index)
			newUser.SetValue("login", login)
			newUser.SetValue("pwd", "testpassword")
			newUser.SetValue("fullname", fmt.Sprintf("Concurrent User %04d", index))
			// Insert the user (transaction is handled internally)
			_, newerr := repo.Insert(newUser)
			if newerr != nil {
				t.Logf("Failed to insert user %s: %v", login, newerr)
				return
			}
			t.Logf("Successfully inserted user: %s", login)

			// Cleanup: delete the test user
			wgDelete.Add(1)
			go func(delUser DBEntityInterface) {
				defer wgDelete.Done()
				_, delErr := repo.Delete(delUser)
				if delErr != nil {
					t.Logf("Warning: Failed to cleanup test user %s: %v", delUser.GetValue("login"), delErr)
				} else {
					t.Logf("Successfully deleted user: %s", delUser.GetValue("login"))
				}
			}(newUser)
		}(i)
	}

	// for i := concurrent_routines - 1; i > 0; i-- {
	// 	<-burstyLimiter
	// 	wgMayhem.Add(1)
	// 	go func(index int) {
	// 		defer wgMayhem.Done()
	// 		// Create a new user entity
	// 		searchUser := factory.GetInstanceByTableName("users")
	// 		if searchUser == nil {
	// 			t.Log("Failed to get DBEntity for 'users'")
	// 			return
	// 		}
	// 		searchUser.SetValue("login", fmt.Sprintf("%s_%04d", userPrefix, index))
	// 		results, err := repo.Search(searchUser, false, true, "")
	// 		if err != nil {
	// 			t.Logf("Failed to search for user %s: %v", searchUser.GetValue("login"), err)
	// 			return
	// 		}
	// 		if len(results) == 0 {
	// 			t.Logf("User %s not found during mayhem search", searchUser.GetValue("login"))
	// 			return
	// 		}
	// 		t.Logf("Mayhem search found user: %s", results[0].ToString())
	// 		_, delErr := repo.Delete(results[0])
	// 		if delErr != nil {
	// 			t.Logf("Warning: Failed to cleanup test user %s: %v", results[0].GetValue("login"), delErr)
	// 		} else {
	// 			t.Logf("Successfully deleted user: %s", results[0].GetValue("login"))
	// 		}
	// 	}(i)
	// }

	wgCreate.Wait()
	wgDelete.Wait()
	wgMayhem.Wait()

	// Search for any remaining users with the mayhem prefix
	searchEntity := Factory.GetInstanceByTableName("users")
	searchEntity.SetValue("login", userPrefix+"_%")
	results, err := repo.Search(searchEntity, false, true, "")
	if err != nil {
		t.Fatal("Failed to search for mayhem users:", err)
	}
	if len(results) != 0 {
		t.Fatalf("Some mayhem users were not deleted, count: %d", len(results))
	} else {
		t.Log("All mayhem users successfully deleted.")
	}

	// Search for the personal groups created
	groupSearchEntity := Factory.GetInstanceByTableName("groups")
	groupSearchEntity.SetValue("name", userPrefix+"_%'s group")
	groupResults, err := repo.Search(groupSearchEntity, false, true, "")
	if err != nil {
		t.Fatal("Failed to search for mayhem groups:", err)
	}
	if len(groupResults) != 0 {
		t.Fatalf("Some mayhem groups were not deleted, count: %d", len(groupResults))
	} else {
		t.Log("All mayhem groups successfully deleted.")
	}
}

func TestCRUDUser(t *testing.T) {
	// Setup
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   "rprj",
	}

	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false

	// Create a new user entity
	newUser := Factory.GetInstanceByTableName("users")
	if newUser == nil {
		t.Fatal("Failed to get DBEntity for 'users'")
	}

	// Set values (id must be set manually since it's not auto-increment)
	login := "testuser_" + Random4digits()
	newUser.SetValue("login", login)
	newUser.SetValue("pwd", "testpassword")
	newUser.SetValue("fullname", "Test User")

	// newUser.SetValue("group_id", "-3") // -3 is the default "users" group

	// Insert the user (transaction is handled internally)
	newUser, newerr := repo.Insert(newUser)
	if newerr != nil {
		t.Fatal("Failed to insert user:", newerr)
	}

	// Verify insertion by searching
	searchEntity := Factory.GetInstanceByTableName("users")
	searchEntity.SetValue("login", login)
	results, err := repo.Search(searchEntity, false, true, "")
	if err != nil {
		t.Fatal("Failed to search for inserted user:", err)
	}
	if len(results) == 0 {
		t.Fatal("Inserted user not found")
	}
	// Verify the group_id is set correctly
	if results[0].GetValue("group_id") == "" {
		t.Fatal("Inserted user's group_id is not set")
	}
	// Verify the group exists
	groupEntity := Factory.GetInstanceByTableName("groups")
	groupEntity.SetValue("id", results[0].GetValue("group_id"))
	groupResults, err := repo.Search(groupEntity, false, true, "")
	if err != nil {
		t.Fatal("Failed to search for user's group:", err)
	}
	if len(groupResults) == 0 {
		t.Fatal("User's group not found")
	}

	// Modify the user's fullname
	results[0].SetValue("fullname", "Updated Test User")
	_, err = repo.Update(results[0])
	if err != nil {
		t.Fatal("Failed to update user:", err)
	}

	// Verify the update
	verifyEntity := Factory.GetInstanceByTableName("users")
	verifyEntity.SetValue("login", login)
	verifyResults, err := repo.Search(verifyEntity, false, true, "")
	if err != nil {
		t.Fatal("Failed to search for updated user:", err)
	}
	if len(verifyResults) == 0 {
		t.Fatal("Updated user not found")
	}
	if verifyResults[0].GetValue("fullname") != "Updated Test User" {
		t.Fatal("User's fullname was not updated correctly")
	}
	log.Printf("Successfully updated user: %s", verifyResults[0].GetValue("login"))

	// Cleanup: delete the test user
	_, err = repo.Delete(newUser)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test user: %v", err)
	}

	// Verify deletion
	finalResults, err := repo.Search(searchEntity, false, true, "")
	if err != nil {
		t.Fatal("Failed to search for deleted user:", err)
	}
	if len(finalResults) != 0 {
		t.Fatal("User was not deleted successfully")
	}
	log.Printf("Successfully deleted user: %s", login)
}

func TestCRUDMayhem(t *testing.T) {
	// Setup
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   "rprj",
	}

	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false

	var wgCreate sync.WaitGroup
	var wgDelete sync.WaitGroup

	// Rate limiter to control the number of concurrent goroutines
	// Buffer size should be much smaller than MaxOpenConns to leave room for nested operations
	// Each Insert creates: user + group + usergroup (3 operations in 1 transaction)
	// Then launches a Delete goroutine (which does: delete usergroup + delete group + delete user)
	max_concurrent := 10 // Small buffer - ticker will refill gradually
	burstyLimiter := make(chan time.Time, max_concurrent)
	for range max_concurrent {
		burstyLimiter <- time.Now()
	}
	go func() {
		for range time.Tick(25 * time.Millisecond) {
			burstyLimiter <- time.Now()
		}
	}()

	concurrent_routines := 300
	userPrefix := "mayhem_" + Random4digits()

	for i := range concurrent_routines {
		<-burstyLimiter // Acquire a token
		wgCreate.Add(1)
		go func(index int) {
			defer wgCreate.Done()
			// Create a new user entity
			newUser := Factory.GetInstanceByTableName("users")
			if newUser == nil {
				t.Log("Failed to get DBEntity for 'users'")
				return
			}

			// Set values (id must be set manually since it's not auto-increment)
			login := fmt.Sprintf("%s_%04d", userPrefix, index)
			newUser.SetValue("login", login)
			newUser.SetValue("pwd", "testpassword")
			newUser.SetValue("fullname", fmt.Sprintf("Concurrent User %04d", index))
			// Insert the user (transaction is handled internally)
			_, newerr := repo.Insert(newUser)
			if newerr != nil {
				t.Logf("Failed to insert user %s: %v", login, newerr)
				return
			}
			t.Logf("Successfully inserted user: %s", login)

			// Verify insertion by searching
			searchEntity := Factory.GetInstanceByTableName("users")
			searchEntity.SetValue("login", login)
			results, err := repo.Search(searchEntity, false, true, "")
			if err != nil {
				t.Logf("Failed to search for inserted user %s: %v", login, err)
				return
			}
			if len(results) == 0 {
				t.Logf("Inserted user %s not found", login)
				return
			}
			t.Logf("Verified inserted user: %s", login)

			// Update the user's fullname
			results[0].SetValue("fullname", "Updated Concurrent User")
			_, err = repo.Update(results[0])
			if err != nil {
				t.Logf("Failed to update user %s: %v", login, err)
				return
			}
			t.Logf("Successfully updated user: %s", login)

			// Verify the update
			verifyEntity := Factory.GetInstanceByTableName("users")
			verifyEntity.SetValue("login", login)
			verifyResults, err := repo.Search(verifyEntity, false, true, "")
			if err != nil {
				t.Logf("Failed to search for updated user %s: %v", login, err)
				return
			}
			if len(verifyResults) == 0 {
				t.Logf("Updated user %s not found", login)
				return
			}
			if verifyResults[0].GetValue("fullname") != "Updated Concurrent User" {
				t.Logf("User %s's fullname was not updated correctly", login)
				return
			}
			t.Logf("Verified updated user: %s", login)

			// Cleanup: delete the test user
			_, err = repo.Delete(verifyResults[0])
			if err != nil {
				t.Logf("Failed to delete user %s: %v", login, err)
				return
			}
			t.Logf("Successfully deleted user: %s", login)
		}(i)
	}

	wgCreate.Wait()
	wgDelete.Wait()

	// Search for any remaining users with the mayhem prefix
	searchEntity := Factory.GetInstanceByTableName("users")
	searchEntity.SetValue("login", userPrefix+"_%")
	results, err := repo.Search(searchEntity, false, true, "")
	if err != nil {
		t.Fatal("Failed to search for mayhem users:", err)
	}
	if len(results) != 0 {
		t.Fatalf("Some mayhem users were not deleted, count: %d", len(results))
	} else {
		t.Log("All mayhem users successfully deleted.")
	}

	// Search for the personal groups created
	groupSearchEntity := Factory.GetInstanceByTableName("groups")
	groupSearchEntity.SetValue("name", userPrefix+"_%'s group")
	groupResults, err := repo.Search(groupSearchEntity, false, true, "")
	if err != nil {
		t.Fatal("Failed to search for mayhem groups:", err)
	}
	if len(groupResults) != 0 {
		t.Fatalf("Some mayhem groups were not deleted, count: %d", len(groupResults))
	} else {
		t.Log("All mayhem groups successfully deleted.")
	}
}
