package dblayer

import (
	"database/sql"
	"fmt"
	"log"
	"time"
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
func (dbVersion *DBVersion) NewInstance() DBEntityInterface {
	return NewDBVersion()
}

/*
CREATE TABLE IF NOT EXISTS oauth_tokens (

	token_id     VARCHAR(64) PRIMARY KEY,
	user_id      VARCHAR(16) NOT NULL,
	access_token TEXT NOT NULL,
	refresh_token TEXT,
	expires_at   DATETIME NOT NULL,
	created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES rra_users(id)

);
*/
type OAuthToken struct {
	DBEntity
}

func NewOAuthToken() *OAuthToken {
	columns := []Column{
		{Name: "token_id", Type: "varchar(64)", Constraints: []string{"PRIMARY KEY"}},
		{Name: "user_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "access_token", Type: "text", Constraints: []string{"NOT NULL"}},
		{Name: "refresh_token", Type: "text", Constraints: []string{}},
		{Name: "expires_at", Type: "datetime", Constraints: []string{"NOT NULL"}},
		{Name: "created_at", Type: "datetime", Constraints: []string{}},
	}
	keys := []string{"token_id"}
	foreignKeys := []ForeignKey{
		{Column: "user_id", RefTable: "users", RefColumn: "id"},
	}
	return &OAuthToken{
		DBEntity: *NewDBEntity(
			"OAuthToken",
			"oauth_tokens",
			columns,
			keys,
			foreignKeys,
			make(map[string]any),
		),
	}
}
func (oAuthToken *OAuthToken) NewInstance() DBEntityInterface {
	return NewOAuthToken()
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
func (dbUser *DBUser) GetValue(columnName string) any {
	return dbUser.DBEntity.GetValue(columnName)
}
func (dbUser *DBUser) SetValue(columnName string, value any) {
	dbUser.DBEntity.SetValue(columnName, value)
}
func (dbUser *DBUser) GetUnencryptedPwd() string {
	// TODO: implement proper password hashing and verification
	return dbUser.GetValue("pwd").(string)
}
func (dbUser *DBUser) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	// 1. Check that user with same login does not already exist
	existingUser := dbUser.NewInstance()
	existingUser.SetValue("login", dbUser.GetValue("login"))
	results, err := dbr.searchWithTx(existingUser, false, false, "login", tx)
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
	group.SetValue("name", dbUser.GetValue("login").(string)+"'s group")
	group.SetValue("description", "Personal group for "+dbUser.GetValue("login").(string))
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
	log.Printf("DBGroup::beforeInsert: id=%s, name=%s, description=%s", dbGroup.GetValue("id"), dbGroup.GetValue("name"), dbGroup.GetValue("description"))

	// Check that group with same name does not already exist
	existingGroup := dbGroup.NewInstance()
	existingGroup.SetValue("name", dbGroup.GetValue("name"))
	results, err := dbr.searchWithTx(existingGroup, false, false, "name", tx)
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

func (dbGroup *DBGroup) GetValue(columnName string) any {
	return dbGroup.DBEntity.GetValue(columnName)
}
func (dbGroup *DBGroup) SetValue(columnName string, value any) {
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
func (dbUserGroup *UserGroup) GetValue(columnName string) any {
	return dbUserGroup.DBEntity.GetValue(columnName)
}
func (dbUserGroup *UserGroup) SetValue(columnName string, value any) {
	dbUserGroup.DBEntity.SetValue(columnName, value)
}

/*
CREATE TABLE IF NOT EXISTS `rra_log` (

	`ip` varchar(16) NOT NULL DEFAULT '',
	`data` date NOT NULL DEFAULT '0000-00-00',
	`ora` time NOT NULL DEFAULT '00:00:00',
	`count` int(11) NOT NULL DEFAULT '0',
	`url` varchar(255) DEFAULT NULL,
	`note` varchar(255) NOT NULL DEFAULT '',
	`note2` text NOT NULL,
	PRIMARY KEY (`ip`,`data`),
	KEY `rra_log_0` (`ip`),
	KEY `rra_log_1` (`data`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBLog struct {
	DBEntity
}

func NewDBLog() *DBLog {
	columns := []Column{
		{Name: "ip", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "data", Type: "date", Constraints: []string{"NOT NULL"}},
		{Name: "ora", Type: "time", Constraints: []string{"NOT NULL"}},
		{Name: "count", Type: "int(11)", Constraints: []string{"NOT NULL"}},
		{Name: "url", Type: "varchar(255)", Constraints: []string{}},
		{Name: "note", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "note2", Type: "text", Constraints: []string{"NOT NULL"}},
	}
	keys := []string{"ip", "data"}
	return &DBLog{
		DBEntity: *NewDBEntity(
			"LogEntry",
			"log",
			columns,
			keys,
			[]ForeignKey{},
			make(map[string]any),
		),
	}
}
func (logEntry *DBLog) NewInstance() DBEntityInterface {
	return NewDBLog()
}
func (logEntry *DBLog) GetValue(columnName string) any {
	return logEntry.DBEntity.GetValue(columnName)
}
func (logEntry *DBLog) SetValue(columnName string, value any) {
	logEntry.DBEntity.SetValue(columnName, value)
}

type DBObjectInterface interface {
	DBEntityInterface
	IsDBObject() bool
	IsDeleted() bool
	CanRead(kind string) bool
	CanWrite(kind string) bool
	CanExecute(kind string) bool
	SetDefaultValues(repo *DBRepository)
}

/*
CREATE TABLE IF NOT EXISTS `rra_objects` (

	`id` varchar(16) NOT NULL DEFAULT '',
	`owner` varchar(16) NOT NULL DEFAULT '',
	`group_id` varchar(16) NOT NULL DEFAULT '',
	`permissions` varchar(9) NOT NULL DEFAULT 'rwx------',
	`creator` varchar(16) NOT NULL DEFAULT '',
	`creation_date` datetime DEFAULT NULL,
	`last_modify` varchar(16) NOT NULL DEFAULT '',
	`last_modify_date` datetime DEFAULT NULL,
	`father_id` varchar(16) DEFAULT NULL,
	`name` varchar(255) NOT NULL DEFAULT '',
	`description` text,
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	PRIMARY KEY (`id`),
	KEY `rra_objects_idx1` (`id`),
	KEY `rra_objects_idx2` (`owner`),
	KEY `rra_objects_idx3` (`name`),
	KEY `rra_objects_idx4` (`creator`),
	KEY `rra_objects_idx5` (`last_modify`),
	KEY `rra_objects_idx6` (`father_id`),
	KEY `rra_timetracks_idx1` (`id`),
	KEY `rra_timetracks_idx2` (`owner`),
	KEY `rra_timetracks_idx3` (`name`),
	KEY `rra_timetracks_idx4` (`creator`),
	KEY `rra_timetracks_idx5` (`last_modify`),
	KEY `rra_timetracks_idx6` (`father_id`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBObject struct {
	DBEntity
}

func NewDBObject() *DBObject {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "varchar(9)", Constraints: []string{"NOT NULL"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{}},
	}
	keys := []string{"id"}
	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "objects", RefColumn: "id"},
	}
	return &DBObject{
		DBEntity: *NewDBEntity(
			"DBObject",
			"objects",
			columns,
			keys,
			foreignKeys,
			make(map[string]any),
		),
	}
}
func (dbObject *DBObject) NewInstance() DBEntityInterface {
	return NewDBObject()
}
func CurrentDateTimeString() string {
	now := time.Now()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
}

func (dbObject *DBObject) IsDBObject() bool {
	return true
}

func (dbObject *DBObject) HasDeletedDate() bool {
	if !dbObject.HasValue("deleted_date") {
		return false
	}
	deletedDate := dbObject.GetValue("deleted_date")
	// NULL = not deleted, any value = deleted
	return deletedDate != nil
}

func (dbObject *DBObject) CanRead(kind string) bool {
	permissions := dbObject.GetValue("permissions").(string)
	if len(permissions) != 9 {
		return false
	}
	switch kind {
	case "U": // User
		return permissions[0] == 'r'
	case "G": // Group
		return permissions[3] == 'r'
	default: // Others
		return permissions[6] == 'r'
	}
}
func (dbObject *DBObject) CanWrite(kind string) bool {
	permissions := dbObject.GetValue("permissions").(string)
	if len(permissions) != 9 {
		return false
	}
	switch kind {
	case "U": // User
		return permissions[1] == 'w'
	case "G": // Group
		return permissions[4] == 'w'
	default: // Others
		return permissions[7] == 'w'
	}
}
func (dbObject *DBObject) CanExecute(kind string) bool {
	permissions := dbObject.GetValue("permissions").(string)
	if len(permissions) != 9 {
		return false
	}
	switch kind {
	case "U": // User
		return permissions[2] == 'x'
	case "G": // Group
		return permissions[5] == 'x'
	default: // Others
		return permissions[8] == 'x'
	}
}
func (dbObject *DBObject) SetDefaultValues(repo *DBRepository) {
	user := repo.GetCurrentUser()
	userID := user.GetValue("id").(string)
	if userID != "" {
		if !dbObject.HasValue("owner") {
			dbObject.SetValue("owner", userID)
		}
		if !dbObject.HasValue("group_id") {
			dbObject.SetValue("group_id", user.GetValue("group_id").(string))
		}
		dbObject.SetValue("creator", userID)
		dbObject.SetValue("last_modify", userID)
	}
	dbObject.SetValue("creation_date", CurrentDateTimeString())
	dbObject.SetValue("last_modify_date", CurrentDateTimeString())
	// dbObject.SetValue("deleted_date", nil) // NULL = not deleted

	if !dbObject.HasValue("father_id") {
		dbObject.SetValue("father_id", nil)

		if dbObject.HasValue("fk_obj_id") && dbObject.GetValue("fk_obj_id") != nil {
			fkobj := repo.ObjectByID(dbObject.GetValue("fk_obj_id").(string), true)
			if fkobj != nil {
				dbObject.SetValue("group_id", fkobj.GetValue("group_id"))
				dbObject.SetValue("permissions", fkobj.GetValue("permissions"))
				dbObject.SetValue("father_id", fkobj.GetValue("id"))
			}
		}
	} else {
		father := repo.ObjectByID(dbObject.GetValue("father_id").(string), true)
		if father != nil {
			dbObject.SetValue("group_id", father.GetValue("group_id"))
			dbObject.SetValue("permissions", father.GetValue("permissions"))
		}
	}
}

func (dbObject *DBObject) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	if !dbObject.HasValue("id") {
		objectID, _ := uuid16HexGo()
		dbObject.SetValue("id", objectID)
	}
	dbObject.SetDefaultValues(dbr)
	if dbr.Verbose {
		log.Println("DBObject.beforeInsert: values=", dbObject.ToJSON())
	}
	return nil
}

func (dbObject *DBObject) beforeUpdate(dbr *DBRepository, tx *sql.Tx) error {
	user := dbr.GetCurrentUser()
	userID := user.GetValue("id").(string)
	if userID != "" {
		dbObject.SetValue("last_modify", userID)
	}
	dbObject.SetValue("last_modify_date", CurrentDateTimeString())
	log.Println("DBObject.beforeUpdate: values=", dbObject.ToJSON())
	return nil
}

func (dbObject *DBObject) beforeDelete(dbr *DBRepository, tx *sql.Tx) error {
	if dbObject.HasDeletedDate() {
		return nil // Already deleted
	}
	user := dbr.GetCurrentUser()
	userID := user.GetValue("id").(string)
	if userID != "" {
		dbObject.SetValue("deleted_by", userID)
	}
	dbObject.SetValue("deleted_date", CurrentDateTimeString())
	if dbr.Verbose {
		log.Println("DBObject.beforeDelete: values=", dbObject.ToJSON())
	}
	return nil
}
