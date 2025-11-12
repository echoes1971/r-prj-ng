package db

import (
	"database/sql"
	"log"

	"rprj/be/models"
)

// CREATE
func CreateGroup(g models.DBGroup) error {
	if g.ID == "" {
		newID, _ := uuid16HexGo()
		log.Print("newID=", newID)
		g.ID = newID
	}
	_, err := DB.Exec(
		"INSERT INTO "+tablePrefix+"groups (id, name, description) VALUES (?, ?, ?)",
		g.ID, g.Name, g.Description,
	)
	return err
}

// READ
func GetGroupByID(id string) (*models.DBGroup, error) {
	row := DB.QueryRow(
		"SELECT id, name, description FROM "+tablePrefix+"groups WHERE id = ?",
		id,
	)

	var g models.DBGroup
	err := row.Scan(&g.ID, &g.Name, &g.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Nessun gruppo trovato
		}
		return nil, err
	}
	return &g, nil
}

// UPDATE
func UpdateGroup(g models.DBGroup) error {
	_, err := DB.Exec(
		"UPDATE "+tablePrefix+"groups SET name = ?, description = ? WHERE id = ?",
		g.Name, g.Description, g.ID,
	)
	return err
}

// DELETE
func DeleteGroup(id string) error {
	_, err := DB.Exec(
		"DELETE FROM "+tablePrefix+"groups WHERE id = ?",
		id,
	)
	return err
}

// SEARCH
func SearchGroupsBy(search string, orderBy string) ([]models.DBGroup, error) {
	query := "SELECT id, name, description FROM " + tablePrefix + "groups"
	if search != "" {
		query += " WHERE name LIKE ? OR description LIKE ?"
	}
	if orderBy != "" {
		query += " ORDER BY " + orderBy
	}

	likePattern := "%" + search + "%"

	var rows *sql.Rows
	var err error
	if search != "" {
		rows, err = DB.Query(query, likePattern, likePattern)
	} else {
		rows, err = DB.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.DBGroup
	for rows.Next() {
		var g models.DBGroup
		err := rows.Scan(&g.ID, &g.Name, &g.Description)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
