package dblayer

import (
	"database/sql"
	"fmt"
	"log"
)

/*
CREATE TABLE `rprj_dbversion` (

	`model_name` varchar(255) NOT NULL,
	`version` int(11) NOT NULL,
	PRIMARY KEY (`model_name`),
	KEY `rprj_dbversion_0` (`model_name`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/
type DBVersion struct {
	DBEntity
}

func NewDBVersion() *DBVersion {
	columns := []Column{
		{Name: "model_name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "version", Type: "int(11)", Constraints: []string{"NOT NULL"}},
	}
	keys := []string{"model_name"}
	return &DBVersion{
		DBEntity: *NewDBEntity(
			"DBVersion",
			"dbversion",
			columns,
			keys,
			[]ForeignKey{},
			make(map[string]any),
		),
	}
}

/*
CREATE TABLE `rprj_users` (

	`id` varchar(16) NOT NULL,
	`login` varchar(255) NOT NULL,
	`pwd` varchar(255) NOT NULL,
	`pwd_salt` varchar(4) DEFAULT '',
	`fullname` text DEFAULT NULL,
	`group_id` varchar(16) NOT NULL,
	PRIMARY KEY (`id`),
	KEY `rprj_users_0` (`id`),
	KEY `rprj_users_1` (`group_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/
type DBUser struct {
	DBEntity
}

func NewDBUser() *DBUser {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "login", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "pwd", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "pwd_salt", Type: "varchar(4)", Constraints: []string{}},
		{Name: "fullname", Type: "text", Constraints: []string{}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
	}
	keys := []string{"id"}
	foreignKeys := []ForeignKey{
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
	}
	return &DBUser{
		DBEntity: *NewDBEntity(
			"DBUser",
			"users",
			columns,
			keys,
			foreignKeys,
			make(map[string]any),
		),
	}
}
func (dbUser *DBUser) NewInstance() DBEntityInterface {
	return NewDBUser()
}
func (dbUser *DBUser) GetValue(columnName string) string {
	return dbUser.DBEntity.GetValue(columnName)
}
func (dbUser *DBUser) SetValue(columnName string, value string) {
	dbUser.DBEntity.SetValue(columnName, value)
}
func (dbUser *DBUser) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	// 1. Check that user with same login does not already exist
	existingUser := dbUser.NewInstance()
	existingUser.SetValue("login", dbUser.GetValue("login"))
	results, err := dbr.Search(existingUser, false, false, "login")
	if err != nil {
		return err
	}
	if len(results) > 0 {
		return fmt.Errorf("user with login '%s' already exists", dbUser.GetValue("login"))
	}

	// 2. Generate IDs
	userID, _ := uuid16HexGo()
	groupID, _ := uuid16HexGo()

	// Create personal group
	group := NewDBGroup()
	group.SetValue("id", groupID)
	group.SetValue("name", dbUser.GetValue("login")+"'s group")
	group.SetValue("description", "Personal group for "+dbUser.GetValue("login"))
	_, err = dbr.insertWithTx(group, tx)
	if err != nil {
		log.Print("DBUser::beforeInsert: error inserting group:", err)
		return err
	}

	// 3. Set user ID and group ID
	dbUser.SetValue("id", userID)
	dbUser.SetValue("group_id", groupID)

	return nil
}
func (dbUser *DBUser) afterInsert(dbr *DBRepository, tx *sql.Tx) error {

	userID := dbUser.GetValue("id")
	groupID := dbUser.GetValue("group_id")

	// 4. Add user to personal group
	userGroup := NewUserGroup()
	userGroup.SetValue("user_id", userID)
	userGroup.SetValue("group_id", groupID)
	_, err := dbr.insertWithTx(userGroup, tx)
	if err != nil {
		log.Print("DBUser::afterInsert: error inserting userGroup:", err)
		return err
	}

	if !dbUser.HasMetadata("group_ids") {
		return nil
	}

	// Assign the user to the specified groups
	groupIDs := dbUser.GetMetadata("group_ids").([]string)
	for _, gID := range groupIDs {
		if gID == groupID {
			continue // Skip personal group (already added)
		}
		userGroup := dbr.GetInstanceByTableName("users_groups").(*UserGroup)
		userGroup.SetValue("user_id", userID)
		userGroup.SetValue("group_id", gID)
		_, err = dbr.insertWithTx(userGroup, tx)
		if err != nil {
			log.Print("DBUser::afterInsert: error inserting userGroup for additional group:", err)
			return err
		}
	}
	return nil
}

func (dbUser *DBUser) afterUpdate(dbr *DBRepository, tx *sql.Tx) error {

	// Update user-group associations if group_ids metadata is set
	if !dbUser.HasMetadata("group_ids") {
		return nil
	}

	userID := dbUser.GetValue("id")
	groupIDs := dbUser.GetMetadata("group_ids").([]string)

	// Delete existing associations
	userGroups := dbr.GetInstanceByTableName("users_groups")
	userGroupFilter := userGroups.NewInstance()
	userGroupFilter.SetValue("user_id", userID)
	results, err := dbr.searchWithTx(userGroupFilter, false, false, "user_id", tx)
	if err != nil {
		return err
	}
	for _, res := range results {
		_, err := dbr.deleteWithTx(res, tx)
		if err != nil {
			log.Print("DBUser::afterUpdate: error deleting existing userGroup:", err)
			return err
		}
	}

	// Add new associations
	for _, gID := range groupIDs {
		userGroup := dbr.GetInstanceByTableName("users_groups")
		userGroup.SetValue("user_id", userID)
		userGroup.SetValue("group_id", gID)
		_, err = dbr.insertWithTx(userGroup, tx)
		if err != nil {
			log.Print("DBUser::afterUpdate: error inserting userGroup for updated groups:", err)
			return err
		}
	}

	return nil
}

func (dbUser *DBUser) beforeDelete(dbr *DBRepository, tx *sql.Tx) error {
	// Delete all user-group associations for this user
	userGroup := NewUserGroup()
	userGroup.SetValue("user_id", dbUser.GetValue("id"))
	results, err := dbr.searchWithTx(userGroup, false, false, "user_id", tx)
	if err != nil {
		return err
	}
	for _, res := range results {
		_, err := dbr.deleteWithTx(res, tx)
		if err != nil {
			log.Print("DBUser::beforeDelete: error deleting userGroup:", err)
			return err
		}
	}
	// Delete personal group
	personalGroup := NewDBGroup()
	personalGroup.SetValue("id", dbUser.GetValue("group_id"))
	_, err = dbr.deleteWithTx(personalGroup, tx)
	if err != nil {
		log.Print("DBUser::beforeDelete: error deleting personal group:", err)
		return err
	}
	return nil
}

/*
CREATE TABLE IF NOT EXISTS `rprj_groups` (

	`id` varchar(16) NOT NULL DEFAULT '',
	`name` varchar(255) NOT NULL DEFAULT '',
	`description` text,
	PRIMARY KEY (`id`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBGroup struct {
	DBEntity
}

func NewDBGroup() *DBGroup {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{}},
	}
	keys := []string{"id"}
	return &DBGroup{
		DBEntity: *NewDBEntity(
			"DBGroup",
			"groups",
			columns,
			keys,
			[]ForeignKey{},
			make(map[string]any),
		),
	}
}

func (dbGroup *DBGroup) NewInstance() DBEntityInterface {
	return NewDBGroup()
}

func (dbGroup *DBGroup) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	log.Print("DBGroup::beforeInsert called")
	if dbGroup.GetValue("id") == "" {
		groupID, _ := uuid16HexGo()
		dbGroup.SetValue("id", groupID)
	}
	log.Printf("DBGroup::beforeInsert: id=%s, name=%s, description=%s", dbGroup.GetValue("id"), dbGroup.GetValue("name"), dbGroup.GetValue("description	"))

	// Check that group with same name does not already exist
	existingGroup := dbGroup.NewInstance()
	existingGroup.SetValue("name", dbGroup.GetValue("name"))
	results, err := dbr.Search(existingGroup, false, false, "name")
	if err != nil {
		return err
	}
	if len(results) > 0 {
		return fmt.Errorf("group with name '%s' already exists", dbGroup.GetValue("name"))
	}
	return nil
}
func (dbGroup *DBGroup) afterUpdate(dbr *DBRepository, tx *sql.Tx) error {
	// Update user-group associations if user_ids metadata is set
	if !dbGroup.HasMetadata("user_ids") {
		return nil
	}

	groupID := dbGroup.GetValue("id")
	userIDs := dbGroup.GetMetadata("user_ids").([]string)

	// Delete existing associations
	userGroups := dbr.GetInstanceByTableName("users_groups")
	userGroupFilter := userGroups.NewInstance()
	userGroupFilter.SetValue("group_id", groupID)
	results, err := dbr.searchWithTx(userGroupFilter, false, false, "group_id", tx)
	if err != nil {
		return err
	}
	for _, res := range results {
		_, err := dbr.deleteWithTx(res, tx)
		if err != nil {
			log.Print("DBGroup::afterUpdate: error deleting existing userGroup:", err)
			return err
		}
	}

	// Add new associations
	for _, uID := range userIDs {
		userGroup := dbr.GetInstanceByTableName("users_groups")
		userGroup.SetValue("user_id", uID)
		userGroup.SetValue("group_id", groupID)
		_, err = dbr.insertWithTx(userGroup, tx)
		if err != nil {
			log.Print("DBGroup::afterUpdate: error inserting userGroup for updated users:", err)
			return err
		}
	}

	return nil
}
func (dbGroup *DBGroup) beforeDelete(dbr *DBRepository, tx *sql.Tx) error {
	// Check if any users are associated with this group
	userGroup := NewUserGroup()
	userGroup.SetValue("group_id", dbGroup.GetValue("id"))
	results, err := dbr.Search(userGroup, false, false, "group_id")
	if err != nil {
		return err
	}
	// NO: this will be handled by the UI.
	// The right logic would be you can delete a group if is a primary group of a user,
	// but could be not friendly to implement here, so the logic will be in the UI layer.
	// if len(results) > 0 {
	// 	return fmt.Errorf("cannot delete group '%s' because it has associated users", dbGroup.GetValue("name"))
	// }

	// Delete all user-group associations for this group
	for _, res := range results {
		_, err := dbr.deleteWithTx(res, tx)
		if err != nil {
			log.Print("DBGroup::beforeDelete: error deleting userGroup:", err)
			return err
		}
	}

	return nil
}

func (dbGroup *DBGroup) GetValue(columnName string) string {
	return dbGroup.DBEntity.GetValue(columnName)
}
func (dbGroup *DBGroup) SetValue(columnName string, value string) {
	dbGroup.DBEntity.SetValue(columnName, value)
}

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
type UserGroup struct {
	DBEntity
}

func NewUserGroup() *UserGroup {
	columns := []Column{
		{Name: "user_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
	}
	keys := []string{"user_id", "group_id"}
	return &UserGroup{
		DBEntity: *NewDBEntity(
			"UserGroup",
			"users_groups",
			columns,
			keys,
			[]ForeignKey{},
			make(map[string]any),
		),
	}
}
func (dbUserGroup *UserGroup) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	return nil
}
func (dbUserGroup *UserGroup) GetValue(columnName string) string {
	return dbUserGroup.DBEntity.GetValue(columnName)
}
func (dbUserGroup *UserGroup) SetValue(columnName string, value string) {
	dbUserGroup.DBEntity.SetValue(columnName, value)
}
