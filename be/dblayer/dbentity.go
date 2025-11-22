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
)

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
	GetCreateTableSQL(dbSchema string) string

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

func (dbEntity *DBEntity) GetCreateTableSQL(dbSchema string) string {
	columnDefs := []string{}
	for _, col := range dbEntity.columns {
		colDef := fmt.Sprintf(" %s %s", col.Name, col.Type)
		if len(col.Constraints) > 0 {
			colDef += " " + strings.Join(col.Constraints, " ")
		}
		columnDefs = append(columnDefs, colDef)
	}
	// Add primary key constraint
	if len(dbEntity.keys) > 0 {
		pkDef := fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(dbEntity.keys, ", "))
		columnDefs = append(columnDefs, pkDef)
	}
	// Add foreign key constraints
	for _, fk := range dbEntity.foreignKeys {
		fkDef := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", fk.Column, fk.RefTable, fk.RefColumn)
		columnDefs = append(columnDefs, fkDef)
	}
	createTableSQL := fmt.Sprintf("CREATE TABLE %s_%s (\n%s\n);", dbSchema, dbEntity.tablename, strings.Join(columnDefs, ",\n"))
	return createTableSQL
}

func (dbEntity *DBEntity) beforeInsert(dbRepository *DBRepository, tx *sql.Tx) error {
	// Implement any logic needed before inserting the entity into the database
	log.Print("DBEntity::beforeInsert: ", dbEntity.ToString())
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
