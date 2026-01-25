package dblayer

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

/**
 *	DBEntity is the base class for all database entities, with 2 main subclasses:
 * 	- DBAssociation represents associations between two entities
 * 	- DBObject represents objects with ownership and permissions
 */

/* Generate a random UUID-like string of 16 hex characters */
func uuid16HexGo() (string, error) {
	b := make([]byte, 8) // 8 bytes = 16 hex chars
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

type ForeignKey struct {
	Column    string
	RefTable  string
	RefColumn string
}

type Column struct {
	Name        string
	Type        string
	Constraints []string
}

type DBEntityInterface interface {
	NewInstance() DBEntityInterface
	GetColumnType(columnName string) string
	GetColumnDefinitions() map[string]Column
	GetTypeName() string
	GetTableName() string
	GetKeys() []string
	GetForeignKeys() []ForeignKey
	GetOrderBy() []string
	GetOrderByString() string
	GetForeignKeysForTable(tableName string) []ForeignKey
	GetForeignKeyDefinition(columnName string) *ForeignKey
	SetValue(columnName string, value any)
	GetValue(columnName string) any
	HasValue(columnName string) bool
	GetAllValues() map[string]any
	SetMetadata(key string, value any)
	GetMetadata(key string) any
	HasMetadata(key string) bool
	GetAllMetadata() map[string]any
	ReadFKFrom(dbe *DBEntity)
	WriteToFK(dbe *DBEntity)
	IsPrimaryKey(columnName string) bool
	IsForeignKey(columnName string) bool
	GetDictionaryKeys() []string
	GetDictionaryValues() []string
	GetKeySetDictionary() map[string]string
	RemoveKeysFromDictionary()
	IsNew() bool
	IsDBObject() bool
	ToString() string
	ToJSON() string

	getDictionary() map[string]any
	beforeInsert(dbRepository *DBRepository, tx *sql.Tx) error
	afterInsert(dbRepository *DBRepository, tx *sql.Tx) error
	beforeUpdate(dbRepository *DBRepository, tx *sql.Tx) error
	afterUpdate(dbRepository *DBRepository, tx *sql.Tx) error
	beforeDelete(dbRepository *DBRepository, tx *sql.Tx) error
	afterDelete(dbRepository *DBRepository, tx *sql.Tx) error
}
type DBEntity struct {
	typename    string
	tablename   string
	columns     map[string]Column
	keys        []string
	foreignKeys []ForeignKey
	dictionary  map[string]any
	metadata    map[string]any // Extra data for business logic, not persisted to DB
}

func NewDBEntity(typename string, tablename string, columns []Column, keys []string, foreignKeys []ForeignKey, dictionary map[string]any) *DBEntity {
	columnsMap := make(map[string]Column)
	for _, col := range columns {
		columnsMap[col.Name] = col
	}
	return &DBEntity{
		typename:    typename,
		tablename:   tablename,
		columns:     columnsMap,
		keys:        keys,
		foreignKeys: foreignKeys,
		dictionary:  dictionary,
	}
}

/* Override */
func (dbEntity *DBEntity) NewInstance() DBEntityInterface {
	columns := make([]Column, 0, len(dbEntity.columns))
	for _, col := range dbEntity.columns {
		columns = append(columns, col)
	}
	return NewDBEntity(dbEntity.typename, dbEntity.tablename, columns, dbEntity.keys, dbEntity.foreignKeys, make(map[string]any))
}

func (dbEntity *DBEntity) GetColumnType(columnName string) string {
	if col, exists := dbEntity.columns[columnName]; exists {
		return col.Type
	}
	return ""
}
func (dbEntity *DBEntity) GetColumnDefinitions() map[string]Column {
	return dbEntity.columns
}
func (dbEntity *DBEntity) GetTypeName() string {
	return dbEntity.typename
}
func (dbEntity *DBEntity) GetTableName() string {
	return dbEntity.tablename
}
func (dbEntity *DBEntity) GetKeys() []string {
	return dbEntity.keys
}
func (dbEntity *DBEntity) GetForeignKeys() []ForeignKey {
	return dbEntity.foreignKeys
}
func (dbEntity *DBEntity) GetOrderBy() []string {
	return dbEntity.GetKeys()
}
func (dbEntity *DBEntity) GetOrderByString() string {
	return strings.Join(dbEntity.GetOrderBy(), ", ")
}
func (dbEntity *DBEntity) GetForeignKeysForTable(tableName string) []ForeignKey {
	var foreignKeysForTable []ForeignKey
	for _, fk := range dbEntity.foreignKeys {
		if fk.RefTable == tableName {
			foreignKeysForTable = append(foreignKeysForTable, fk)
		}
	}
	return foreignKeysForTable
}
func (dbEntity *DBEntity) GetForeignKeyDefinition(columnName string) *ForeignKey {
	for _, fk := range dbEntity.foreignKeys {
		if fk.Column == columnName {
			return &fk
		}
	}
	return nil
}

// TODO? Manage different types of values (int, date, etc.)
func (dbEntity *DBEntity) SetValue(columnName string, value any) {
	// if _, exists := dbEntity.dictionary[columnName]; exists {
	dbEntity.dictionary[columnName] = value
	// }
}
func (dbEntity *DBEntity) GetValue(columnName string) any {
	if val, exists := dbEntity.dictionary[columnName]; exists {
		return val
	}
	return nil
}
func (dbEntity *DBEntity) HasValue(columnName string) bool {
	_, exists := dbEntity.dictionary[columnName]
	return exists
}
func (dbEntity *DBEntity) GetAllValues() map[string]any {
	// Return a copy to prevent external modification
	valuesCopy := make(map[string]any)
	for k, v := range dbEntity.dictionary {
		valuesCopy[k] = v
	}
	return valuesCopy
}

// SetMetadata sets a metadata value that won't be persisted to the database
// Useful for passing extra data to beforeInsert/beforeUpdate hooks
func (dbEntity *DBEntity) SetMetadata(key string, value any) {
	if dbEntity.metadata == nil {
		dbEntity.metadata = make(map[string]any)
	}
	dbEntity.metadata[key] = value
}

// GetMetadata retrieves a metadata value
func (dbEntity *DBEntity) GetMetadata(key string) any {
	if dbEntity.metadata == nil {
		return nil
	}
	return dbEntity.metadata[key]
}

func (dbEntity *DBEntity) HasMetadata(key string) bool {
	if dbEntity.metadata == nil {
		return false
	}
	_, exists := dbEntity.metadata[key]
	return exists
}

func (dbEntity *DBEntity) GetAllMetadata() map[string]any {
	if dbEntity.metadata == nil {
		return nil
	}
	// Return a copy to prevent external modification
	metadataCopy := make(map[string]any)
	for k, v := range dbEntity.metadata {
		metadataCopy[k] = v
	}
	return metadataCopy
}

func (dbEntity *DBEntity) ReadFKFrom(dbe *DBEntity) {
	fks := dbEntity.GetForeignKeysForTable(dbe.GetTableName())
	for _, fk := range fks {
		value := dbe.GetValue(fk.RefColumn)
		dbEntity.SetValue(fk.Column, value)
	}
}
func (dbEntity *DBEntity) WriteToFK(dbe *DBEntity) {
	fks := dbEntity.GetForeignKeysForTable(dbe.GetTableName())
	for _, fk := range fks {
		value := dbEntity.GetValue(fk.Column)
		dbe.SetValue(fk.RefColumn, value)
	}
}
func (dbEntity *DBEntity) IsPrimaryKey(columnName string) bool {
	for _, key := range dbEntity.keys {
		if key == columnName {
			return true
		}
	}
	return false
}
func (dbEntity *DBEntity) IsForeignKey(columnName string) bool {
	for _, fk := range dbEntity.foreignKeys {
		if fk.Column == columnName {
			return true
		}
	}
	return false
}

/*
Returns the dictionary keys which means all values set in the entity
*/
func (dbEntity *DBEntity) GetDictionaryKeys() []string {
	keys := make([]string, 0, len(dbEntity.dictionary))
	for key := range dbEntity.dictionary {
		keys = append(keys, key)
	}
	// Sort the keys alphabetically
	sort.Strings(keys)
	return keys
}
func (dbEntity *DBEntity) GetDictionaryValues() []string {
	keys := dbEntity.GetDictionaryKeys() // If I use this, the sorting of the keys may be unnecessary
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, dbEntity.dictionary[key].(string))
	}
	return values
}

func (dbEntity *DBEntity) getDictionary() map[string]any {
	return dbEntity.dictionary
}

/*
Returns a dictionary of the keys set in the entity
*/
func (dbEntity *DBEntity) GetKeySetDictionary() map[string]string {
	result := make(map[string]string)
	for _, key := range dbEntity.keys {
		if val, exists := dbEntity.dictionary[key]; exists {
			result[key] = val.(string)
		}
	}
	return result
}

/*
Remove keys from dictionary
*/
func (dbEntity *DBEntity) RemoveKeysFromDictionary() {
	for _, key := range dbEntity.keys {
		delete(dbEntity.dictionary, key)
	}
}

/*
Returns true if all primary keys have not been set in the dictionary
*/
func (dbEntity *DBEntity) IsNew() bool {
	for _, key := range dbEntity.keys {
		if _, exists := dbEntity.dictionary[key]; exists {
			return false
		}
	}
	return true
}
func (dbEntity *DBEntity) IsDBObject() bool {
	return false
}

func (dbEntity *DBEntity) ToString() string {
	return fmt.Sprintf("%s(%v)", dbEntity.typename, dbEntity.ToJSON())
}
func (dbEntity *DBEntity) ToJSON() string {

	data := map[string]any{
		"data": dbEntity.dictionary,
	}
	if dbEntity.metadata != nil {
		data["metadata"] = dbEntity.metadata
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

func (dbEntity *DBEntity) beforeInsert(dbRepository *DBRepository, tx *sql.Tx) error {
	// Implement any logic needed before inserting the entity into the database
	return nil
}

func (dbEntity *DBEntity) afterInsert(dbRepository *DBRepository, tx *sql.Tx) error {
	// Implement any logic needed after inserting the entity into the database
	return nil
}

func (dbEntity *DBEntity) beforeUpdate(dbRepository *DBRepository, tx *sql.Tx) error {
	// Implement any logic needed before updating the entity in the database
	return nil
}

func (dbEntity *DBEntity) afterUpdate(dbRepository *DBRepository, tx *sql.Tx) error {
	// Implement any logic needed after updating the entity in the database
	return nil
}

func (dbEntity *DBEntity) beforeDelete(dbRepository *DBRepository, tx *sql.Tx) error {
	// Implement any logic needed before deleting the entity from the database
	return nil
}

func (dbEntity *DBEntity) afterDelete(dbRepository *DBRepository, tx *sql.Tx) error {
	// Implement any logic needed after deleting the entity from the database
	return nil
}

/**
 *	An association between two DBEntities
 *  It MUST have 2 foreign keys: from_table_id and to_table_id
 */
type DBAssociationInterface interface {
	GetFromTableName() string
	GetToTableName() string
	IsDBObject() bool
}

type DBAssociation struct {
	DBEntity
}

func (dbAssoc *DBAssociation) GetFromTableName() (string, error) {
	// The first foreign key represents the "from" table
	if len(dbAssoc.foreignKeys) < 2 {
		return "", fmt.Errorf("DBAssociation::GetFromTableName: not enough foreign keys")
	}
	return dbAssoc.DBEntity.foreignKeys[0].RefTable, nil
}
func (dbAssoc *DBAssociation) GetToTableName() (string, error) {
	// The second foreign key represents the "to" table
	if len(dbAssoc.foreignKeys) < 2 {
		return "", fmt.Errorf("DBAssociation::GetToTableName: not enough foreign keys")
	}
	return dbAssoc.DBEntity.foreignKeys[1].RefTable, nil
}
func (dbAssoc *DBAssociation) IsDBObject() bool {
	return false
}

type DBObjectInterface interface {
	DBEntityInterface
	IsDBObject() bool
	HasDeletedDate() bool
	CanRead(kind string) bool
	CanWrite(kind string) bool
	CanExecute(kind string) bool
	SetDefaultValues(repo *DBRepository)
	beforeInsert(dbr *DBRepository, tx *sql.Tx) error
	beforeUpdate(dbr *DBRepository, tx *sql.Tx) error
	beforeDelete(dbr *DBRepository, tx *sql.Tx) error
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
		{Name: "permissions", Type: "char(9)", Constraints: []string{"NOT NULL"}},
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
	userID := ""
	if user != nil {
		userID = user.GetValue("id").(string)
	}
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

	if !dbObject.HasValue("permissions") {
		dbObject.SetValue("permissions", "rwx------")
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
