package db

import (
	"database/sql"
	"fmt"

	"rprj/be/models"
)

// CREATE
// func CreateUser(u models.DBUser) (string, error) {
// 	if u.ID == "" {
// 		newID, _ := uuid16HexGo()
// 		log.Print("newID=", newID)
// 		u.ID = newID
// 	}
// 	_, err := DB.Exec(
// 		"INSERT INTO "+tablePrefix+"users (id, login, pwd, pwd_salt, fullname, group_id) VALUES (?, ?, ?, ?, ?, ?)",
// 		u.ID, u.Login, u.Pwd, u.PwdSalt, u.Fullname, u.GroupID,
// 	)
// 	return u.ID, err
// }

// READ (by Login)
func GetUserByLogin(login string) (*models.DBUser, error) {
	row := DB.QueryRow("SELECT id, login, pwd, pwd_salt, fullname, group_id FROM "+tablePrefix+"users WHERE login = ?", login)
	var u models.DBUser
	err := row.Scan(&u.ID, &u.Login, &u.Pwd, &u.PwdSalt, &u.Fullname, &u.GroupID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

// READ (by ID)
func GetUserByID(id string) (*models.DBUser, error) {
	row := DB.QueryRow("SELECT id, login, pwd, pwd_salt, fullname, group_id FROM "+tablePrefix+"users WHERE id = ?", id)
	var u models.DBUser
	err := row.Scan(&u.ID, &u.Login, &u.Pwd, &u.PwdSalt, &u.Fullname, &u.GroupID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

// UPDATE
// func UpdateUser(u models.DBUser, updatePwd bool) error {
// 	if updatePwd {
// 		_, err := DB.Exec(
// 			"UPDATE "+tablePrefix+"users SET login=?, pwd=?, pwd_salt=?, fullname=?, group_id=? WHERE id=?",
// 			u.Login, u.Pwd, u.PwdSalt, u.Fullname, u.GroupID, u.ID,
// 		)
// 		return err
// 	}
// 	// Update without password
// 	_, err := DB.Exec(
// 		"UPDATE "+tablePrefix+"users SET login=?, fullname=?, group_id=? WHERE id=?",
// 		u.Login, u.Fullname, u.GroupID, u.ID,
// 	)
// 	return err
// }

// DELETE
func DeleteUser(id string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get user to find personal group_id
	var groupID string
	err = tx.QueryRow("SELECT group_id FROM "+tablePrefix+"users WHERE id=?", id).Scan(&groupID)
	if err != nil {
		return err
	}

	// Delete user-group associations
	_, err = tx.Exec("DELETE FROM "+tablePrefix+"users_groups WHERE user_id=?", id)
	if err != nil {
		return err
	}

	// Delete user
	_, err = tx.Exec("DELETE FROM "+tablePrefix+"users WHERE id=?", id)
	if err != nil {
		return err
	}

	// Delete personal group (if it exists)
	if groupID != "" {
		_, err = tx.Exec("DELETE FROM "+tablePrefix+"groups WHERE id=?", groupID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Get all users
func GetAllUsers(searchBy string, orderBy string) ([]models.DBUser, error) {
	if orderBy == "" {
		orderBy = "id"
	}
	query := "SELECT id, login, pwd, pwd_salt, fullname, group_id FROM " + tablePrefix + "users"
	if searchBy != "" {
		query += " WHERE login LIKE ? OR fullname LIKE ?"
		searchPattern := "%" + searchBy + "%"
		query += fmt.Sprintf(" ORDER BY %s", orderBy)
		rows, err := DB.Query(query, searchPattern, searchPattern)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var users []models.DBUser
		for rows.Next() {
			var u models.DBUser
			if err := rows.Scan(&u.ID, &u.Login, &u.Pwd, &u.PwdSalt, &u.Fullname, &u.GroupID); err != nil {
				return nil, err
			}
			users = append(users, u)
		}
		return users, nil
	}

	// Senza filtro di ricerca
	rows, err := DB.Query("SELECT id, login, pwd, pwd_salt, fullname, group_id FROM " + tablePrefix + "users ORDER BY " + orderBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.DBUser
	for rows.Next() {
		var u models.DBUser
		if err := rows.Scan(&u.ID, &u.Login, &u.Pwd, &u.PwdSalt, &u.Fullname, &u.GroupID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// EXTRA: Count
func CountUsers() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM " + tablePrefix + "users").Scan(&count)
	if err != nil {
		return 0, err
	}
	fmt.Println("Numero utenti:", count)
	return count, nil
}

// CreateUser creates a user, personal group, and associations in a single transaction
func CreateUser(u models.DBUser, login string, additionalGroupIDs []string) (*models.DBUser, string, error) {
	tx, err := DB.Begin()
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback() // Rollback if not committed

	// Check that user with same login does not already exist
	var existingUserID string
	err = tx.QueryRow("SELECT id FROM "+tablePrefix+"users WHERE login = ?", login).Scan(&existingUserID)
	if err != sql.ErrNoRows {
		if err == nil {
			return nil, "", fmt.Errorf("user with login '%s' already exists", login)
		}
		return nil, "", err
	}

	// Generate IDs
	userID, _ := uuid16HexGo()
	groupID, _ := uuid16HexGo()

	// Create personal group
	_, err = tx.Exec(
		"INSERT INTO "+tablePrefix+"groups (id, name, description) VALUES (?, ?, ?)",
		groupID, login+"'s group", "Personal group for "+login,
	)
	if err != nil {
		return nil, "", err
	}

	// Create user with personal group as primary
	u.ID = userID
	u.GroupID = groupID
	_, err = tx.Exec(
		"INSERT INTO "+tablePrefix+"users (id, login, pwd, pwd_salt, fullname, group_id) VALUES (?, ?, ?, ?, ?, ?)",
		u.ID, u.Login, u.Pwd, u.PwdSalt, u.Fullname, u.GroupID,
	)
	if err != nil {
		return nil, "", err
	}

	// Add user to personal group
	_, err = tx.Exec(
		"INSERT INTO "+tablePrefix+"users_groups (user_id, group_id) VALUES (?, ?)",
		userID, groupID,
	)
	if err != nil {
		return nil, "", err
	}

	// Add user to additional groups
	for _, gID := range additionalGroupIDs {
		if gID == groupID {
			continue // Skip personal group (already added)
		}
		_, err = tx.Exec(
			"INSERT INTO "+tablePrefix+"users_groups (user_id, group_id) VALUES (?, ?)",
			userID, gID,
		)
		if err != nil {
			return nil, "", err
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, "", err
	}

	return &u, groupID, nil
}

// UpdateUser updates user and group associations in a single transaction
func UpdateUser(u models.DBUser, updatePwd bool, groupIDs []string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update user
	if updatePwd {
		_, err = tx.Exec(
			"UPDATE "+tablePrefix+"users SET login=?, pwd=?, pwd_salt=?, fullname=?, group_id=? WHERE id=?",
			u.Login, u.Pwd, u.PwdSalt, u.Fullname, u.GroupID, u.ID,
		)
	} else {
		_, err = tx.Exec(
			"UPDATE "+tablePrefix+"users SET login=?, fullname=?, group_id=? WHERE id=?",
			u.Login, u.Fullname, u.GroupID, u.ID,
		)
	}
	if err != nil {
		return err
	}

	// Delete all existing group associations
	_, err = tx.Exec("DELETE FROM "+tablePrefix+"users_groups WHERE user_id=?", u.ID)
	if err != nil {
		return err
	}

	// Recreate group associations
	for _, groupID := range groupIDs {
		_, err = tx.Exec(
			"INSERT INTO "+tablePrefix+"users_groups (user_id, group_id) VALUES (?, ?)",
			u.ID, groupID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
