package db

import (
	"rprj/be/models"
)

/*
CREATE TABLE IF NOT EXISTS `rra_users_groups` (

	`user_id` varchar(16) NOT NULL DEFAULT '',
	`group_id` varchar(16) NOT NULL DEFAULT '',
	PRIMARY KEY (`user_id`,`group_id`),
	KEY `rra_users_groups_idx1` (`user_id`),
	KEY `rra_users_groups_idx2` (`group_id`),
	KEY `rra_users_groups_idx3` (`user_id`,`group_id`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
// CREATE
func CreateUserGroup(g models.DBUserGroup) error {
	_, err := DB.Exec(
		"INSERT INTO "+tablePrefix+"users_groups (user_id, group_id) VALUES (?, ?)",
		g.UserID, g.GroupID,
	)
	return err
}

// READ
func GetUserGroupsByUserID(userID string) ([]models.DBUserGroup, error) {
	rows, err := DB.Query(
		"SELECT user_id, group_id FROM "+tablePrefix+"users_groups WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.DBUserGroup
	for rows.Next() {
		var g models.DBUserGroup
		if err := rows.Scan(&g.UserID, &g.GroupID); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// DELETE
func DeleteUserGroup(userID string, groupID string) error {
	_, err := DB.Exec(
		"DELETE FROM "+tablePrefix+"users_groups WHERE user_id = ? AND group_id = ?",
		userID, groupID,
	)
	return err
}

// DELETE all groups for a user
func DeleteUserGroupsByUserID(userID string) error {
	_, err := DB.Exec(
		"DELETE FROM "+tablePrefix+"users_groups WHERE user_id = ?",
		userID,
	)
	return err
}

// DELETE all users for a group
func DeleteUserGroupsByGroupID(groupID string) error {
	_, err := DB.Exec(
		"DELETE FROM "+tablePrefix+"users_groups WHERE group_id = ?",
		groupID,
	)
	return err
}
