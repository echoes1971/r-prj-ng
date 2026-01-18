package dblayer

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type DBContext struct {
	UserID   string
	GroupIDs []string

	Schema string // prefix to add to a table name
}

func (dbctx *DBContext) IsInGroup(groupID string) bool {
	for _, gid := range dbctx.GroupIDs {
		if gid == groupID {
			return true
		}
	}
	return false
}
func (dbctx *DBContext) AddGroup(groupID string) {
	for _, gid := range dbctx.GroupIDs {
		if gid == groupID {
			return
		}
	}
	dbctx.GroupIDs = append(dbctx.GroupIDs, groupID)
}
func (dbctx *DBContext) IsUser(userID string) bool {
	return dbctx.UserID == userID
}

type DBRepository struct {
	Verbose     bool
	DbContext   *DBContext
	factory     *DBEFactory
	currentUser *DBUser

	/* Can be a connection to mysql, postgresql, sqlite, etc. */
	DbConnection *sql.DB
}

func NewDBRepository(dbContext *DBContext, factory *DBEFactory, dbConnection *sql.DB) *DBRepository {
	return &DBRepository{
		Verbose:      false,
		DbContext:    dbContext,
		factory:      factory,
		DbConnection: dbConnection,
	}
}

// placeholder returns the correct SQL placeholder for the given parameter index
// PostgreSQL uses $1, $2, $3..., while MySQL/SQLite use ?
func (dbr *DBRepository) placeholder(index int) string {
	if strings.Contains(dbEngine, "postgres") {
		return fmt.Sprintf("$%d", index)
	}
	return "?"
}

func (dbr *DBRepository) GetDBVersion() int {
	search := dbr.GetInstanceByTableName("dbversion")
	if search == nil {
		return -1
	}
	search.SetValue("model_name", DbSchema)
	foundEntities, err := dbr.Search(search, false, false, "")
	if err != nil || len(foundEntities) == 0 {
		return -1
	}
	dbVersion, ok := foundEntities[0].(*DBVersion)
	if !ok {
		return -1
	}
	return dbVersion.GetValue("version").(int)
}
func (dbr *DBRepository) SetDBVersion(version int) error {
	search := dbr.GetInstanceByTableName("dbversion")
	if search == nil {
		return fmt.Errorf("DBRepository::SetDBVersion: cannot create dbversion instance")
	}
	search.SetValue("model_name", DbSchema)
	foundEntities, err := dbr.Search(search, false, false, "")
	if err != nil {
		return err
	}
	if len(foundEntities) == 0 {
		// Insert new
		newVersion := dbr.GetInstanceByTableName("dbversion")
		if newVersion == nil {
			return fmt.Errorf("DBRepository::SetDBVersion: cannot create dbversion instance for insert")
		}
		newVersion.SetValue("model_name", DbSchema)
		newVersion.SetValue("version", version)
		_, err := dbr.Insert(newVersion)
		return err
	} else {
		// Update existing
		dbVersion, ok := foundEntities[0].(*DBVersion)
		if !ok {
			return fmt.Errorf("DBRepository::SetDBVersion: cannot cast found entity to DBVersion")
		}
		dbVersion.SetValue("version", version)
		_, err := dbr.Update(dbVersion)
		return err
	}
}

func (dbr *DBRepository) GetInstanceByClassName(classname string) DBEntityInterface {
	return dbr.factory.GetInstanceByClassName(classname)
}
func (dbr *DBRepository) GetInstanceByTableName(tablename string) DBEntityInterface {
	return dbr.factory.GetInstanceByTableName(tablename)
}

func (dbr *DBRepository) buildTableName(dbe DBEntityInterface) string {
	tablename := dbe.GetTableName()
	// Handle lightweight objects fetched by ObjectByID. Too much of spaghetti code?
	if tablename == "objects" && dbe.HasMetadata("classname") {
		classname := dbe.GetMetadata("classname").(string)
		tmp := dbr.factory.GetInstanceByClassName(classname)
		if tmp != nil {
			tablename = tmp.GetTableName()
		}
	}

	if dbr.DbContext != nil && dbr.DbContext.Schema != "" {
		return dbr.DbContext.Schema + "_" + tablename
	}
	return tablename
}

func (dbr *DBRepository) Search(dbe DBEntityInterface, useLike bool, caseSensitive bool, orderBy string) ([]DBEntityInterface, error) {
	return dbr.searchWithTx(dbe, useLike, caseSensitive, orderBy, nil)
}

// searchWithTx is an internal method that performs the search using an existing transaction (if provided)
func (dbr *DBRepository) searchWithTx(dbe DBEntityInterface, useLike bool, caseSensitive bool, orderBy string, tx *sql.Tx) ([]DBEntityInterface, error) {
	if dbr.Verbose {
		log.Print("DBRepository::searchWithTx: dbe=", dbe.ToString())
	}

	// 1. Build WHERE clauses
	clauses := make([]string, 0)
	args := make([]interface{}, 0) // slice of interface{} for values

	clausesOr := make([]string, 0)
	argsOr := make([]interface{}, 0)
	if conditions, hasOr := dbe.GetMetadata("or").([]interface{}); hasOr {
		for _, condition := range conditions {
			log.Print("DBRepository::searchWithTx: OR condition =", condition)
			if condMap, ok := condition.(map[string]interface{}); ok {
				// log.Print("DBRepository::searchWithTx: OR condMap =", condMap)
				for column, expr := range condMap {
					for operator, value := range expr.(map[string]interface{}) {
						clauseString := column + " " + operator[1:] + " " + dbr.placeholder(len(argsOr)+1)
						log.Print("DBRepository::searchWithTx: OR clauseString =", clauseString, " value=", value)
						clausesOr = append(clausesOr, "("+clauseString+")")
						argsOr = append(argsOr, value)
					}
				}
			}
		}
	}

	// Default search: AND all populated fields
	for key, value := range dbe.getDictionary() {
		if dbr.Verbose {
			log.Print("DBRepository::searchWithTx: processing key=", key, " value=", value)
		}
		if value == nil {
			switch key {
			case "father_id":
				clauses = append(clauses, "("+key+" IS NULL OR "+key+" = '0')")
			default:
				clauses = append(clauses, key+" IS NULL")
			}
			continue
		}
		if key == "father_id" && value == "0" {
			clauses = append(clauses, key+" = '0'")
			// whereClauses = append(whereClauses, "("+key+" IS NULL OR "+key+" = '0')")
			continue
		}
		if dbr.Verbose {
			log.Print("DBRepository::searchWithTx: key=", key, " value=", value)
		}
		if rangeSlice, ok := value.([]string); ok && len(rangeSlice) == 2 {
			if dbr.Verbose {
				log.Print("DBRepository::searchWithTx: detected range for key=", key, " from=", rangeSlice[0], " to=", rangeSlice[1])
			}
			// Range search
			if rangeSlice[0] != "" {
				clauses = append(clauses, key+" >= "+dbr.placeholder(len(args)+1))
				args = append(args, rangeSlice[0])
			}
			if rangeSlice[1] != "" {
				clauses = append(clauses, key+" <= "+dbr.placeholder(len(args)+1))
				args = append(args, rangeSlice[1])
			}
			continue
		}
		if valueStr, ok := value.(string); ok && valueStr == "" {
			continue // skip empty strings
		}
		if useLike {
			// For strings: LIKE '%value%'
			if strings.Contains(dbe.GetColumnType(key), "varchar") || dbe.GetColumnType(key) == "text" {
				if caseSensitive {
					clauses = append(clauses, key+" LIKE "+dbr.placeholder(len(args)+1))
					args = append(args, "%"+fmt.Sprint(value)+"%")
				} else {
					clauses = append(clauses, "LOWER("+key+") LIKE LOWER("+dbr.placeholder(len(args)+1)+")")
					args = append(args, "%"+fmt.Sprint(value)+"%")
				}
			} else {
				// Per numeri/date: exact match
				clauses = append(clauses, key+" = "+dbr.placeholder(len(args)+1))
				args = append(args, value)
			}
		} else {
			// Exact match
			if dbr.Verbose {
				log.Print("DBRepository::searchWithTx: adding exact match for key=", key, " value=", value)
			}
			clauses = append(clauses, key+" = "+dbr.placeholder(len(args)+1))
			args = append(args, value)
		}
	}

	// 2. Build the final query
	query := "SELECT * FROM " + dbr.buildTableName(dbe)
	if len(clauses) > 0 || len(clausesOr) > 0 {
		query += " WHERE "
		if len(clauses) > 0 {
			query += strings.Join(clauses, " AND ")
		}
		if len(clausesOr) > 0 {
			if len(clauses) > 0 {
				query += " AND "
			}
			query += "(" + strings.Join(clausesOr, " OR ") + ")"
			args = append(args, argsOr...)
		}
	}
	if orderBy != "" {
		query += " ORDER BY " + orderBy
	}

	if dbr.Verbose {
		log.Print("DBRepository::searchWithTx: query=", query, " args=", args)
	}

	// 3. Execute the query (use transaction if provided, otherwise use connection)
	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(query, args...)
	} else {
		rows, err = dbr.DbConnection.Query(query, args...)
	}
	if err != nil {
		log.Print("DBRepository::searchWithTx: Query error:", err)
		return nil, err
	}
	defer rows.Close()

	// 4. Process results
	results := make([]DBEntityInterface, 0)
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		// Create a new instance of the DBEntity
		resultEntity := dbe.NewInstance()
		resultEntity.SetMetadata("classname", dbe.GetTypeName())

		// Prepare a slice of interfaces to hold column values
		columnValues := make([]interface{}, len(columns))
		columnValuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			columnValuePtrs[i] = &columnValues[i]
		}

		// Scan the row into the column value pointers
		if err := rows.Scan(columnValuePtrs...); err != nil {
			return nil, err
		}

		// Map column values to the result entity's dictionary
		for i, colName := range columns {
			val := columnValues[i]
			if b, ok := val.([]byte); ok {
				resultEntity.SetValue(colName, string(b))
			} else if val != nil {
				resultEntity.SetValue(colName, fmt.Sprint(val))
				// } else {
				// 	resultEntity.SetValue(colName, "")
			}
		}

		results = append(results, resultEntity)
	}

	if dbr.Verbose {
		log.Printf("DBRepository::Search: found %d results", len(results))
	}

	// 5. Return results

	return results, nil
}

// CreateObject creates and inserts a new entity with the provided values
// Usage: repo.CreateObject("files", map[string]any{"name": "Test", "filename": "test.jpg"})
func (dbr *DBRepository) CreateObject(tableName string, values map[string]any, metadata map[string]any) (DBEntityInterface, error) {
	instance := dbr.factory.GetInstanceByTableNameWithValues(tableName, values, metadata)
	if instance == nil {
		return nil, fmt.Errorf("unknown table name: %s", tableName)
	}
	return dbr.Insert(instance)
}

// UpdateObject updates an existing entity with the provided values (only updates specified fields)
// Usage: repo.UpdateObject("files", "file-id-123", map[string]any{"name": "New Name", "description": "Updated"})
func (dbr *DBRepository) UpdateObject(tableName string, id string, values map[string]any, metadata map[string]any) (DBEntityInterface, error) {
	// Get the existing entity
	existing := dbr.GetEntityByID(tableName, id)
	if existing == nil {
		return nil, fmt.Errorf("entity not found: %s with id %s", tableName, id)
	}

	// TODO if we want to update only the selected fields, we should use factory.GetInstanceByTableNameWithValues

	// Update only the provided values
	for key, value := range values {
		existing.SetValue(key, value)
	}
	for key, value := range metadata {
		existing.SetMetadata(key, value)
	}

	return dbr.Update(existing)
}

// Insert inserts a new entity into the database within a transaction
func (dbr *DBRepository) Insert(dbe DBEntityInterface) (DBEntityInterface, error) {
	// Start a transaction
	tx, err := dbr.DbConnection.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Use internal method with transaction
	result, err := dbr.insertWithTx(dbe, tx)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

// insertWithTx is an internal method that performs the insert using an existing transaction
func (dbr *DBRepository) insertWithTx(dbe DBEntityInterface, tx *sql.Tx) (DBEntityInterface, error) {
	if dbr.Verbose {
		log.Print("DBRepository::insertWithTx: dbe=", dbe.ToString())
	}

	// Call beforeInsert hook (which can use dbr.insertWithTx for nested inserts)
	err := dbe.beforeInsert(dbr, tx)
	if err != nil {
		log.Print("DBRepository::insertWithTx: beforeInsert error:", err)
		return nil, err
	}

	// 1. Build INSERT query dynamically based on populated fields
	columns := make([]string, 0)
	placeholders := make([]string, 0)
	args := make([]interface{}, 0)

	paramIndex := 1
	for key, value := range dbe.getDictionary() {
		columns = append(columns, key)
		// Use $1, $2, $3... for PostgreSQL, ? for MySQL/SQLite
		if strings.Contains(dbEngine, "postgres") {
			placeholders = append(placeholders, fmt.Sprintf("$%d", paramIndex))
			paramIndex++
		} else {
			placeholders = append(placeholders, "?")
		}
		args = append(args, value)
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("no fields to insert")
	}

	// 2. Build the INSERT query
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		dbr.buildTableName(dbe),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	// IF the engine is PostgreSQL, add ; at the end of the query
	if strings.Contains(dbEngine, "postgres") {
		query += ";"
	}

	if dbr.Verbose {
		log.Print("DBRepository::insertWithTx: query=", query, " args=", args)
	}

	// 3. Execute the INSERT using the transaction
	result, err := tx.Exec(query, args...)
	if err != nil {
		log.Print("DBRepository::insertWithTx: Exec error:", err)
		argsString := make([]string, len(args))
		for i, a := range args {
			argsString[i] = "'" + fmt.Sprint(a) + "',"
		}
		log.Print("DBRepository::insertWithTx: Error - query=", query, " args=", argsString)
		return nil, err
	}

	// 4. Get the last inserted ID (if auto-increment)
	lastID, err := result.LastInsertId()
	if err == nil && dbr.Verbose {
		log.Printf("DBRepository::insertWithTx: inserted with ID=%d", lastID)
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && dbr.Verbose {
		log.Printf("DBRepository::insertWithTx: rows affected=%d", rowsAffected)
	}

	err = dbe.afterInsert(dbr, tx)
	if err != nil {
		log.Print("DBRepository::insertWithTx: afterInsert error:", err)
		return nil, err
	}

	return dbe, nil
}

func (dbr *DBRepository) Delete(dbe DBEntityInterface) (DBEntityInterface, error) {
	// Start a transaction
	tx, err := dbr.DbConnection.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Use internal method with transaction
	result, err := dbr.deleteWithTx(dbe, tx)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

// deleteWithTx is an internal method that performs the delete using an existing transaction
func (dbr *DBRepository) deleteWithTx(dbe DBEntityInterface, tx *sql.Tx) (DBEntityInterface, error) {
	if dbr.Verbose {
		log.Print("DBRepository::deleteWithTx: dbe=", dbe.ToString())
	}

	if dbe.IsDBObject() {
		dbObj := dbe.(DBObjectInterface)
		// IF has not deleted date
		if !dbObj.HasDeletedDate() {
			// Call beforeDelete
			err := dbObj.beforeDelete(dbr, tx)
			if err != nil {
				log.Print("DBRepository::deleteWithTx: beforeDelete error:", err)
				return nil, err
			}
			// Build UPDATE query dynamicallyto set deleted_date and deleted_by
			query := fmt.Sprintf("UPDATE %s SET deleted_date = %s, deleted_by = %s WHERE id='%s'",
				dbr.buildTableName(dbe), dbr.placeholder(1), dbr.placeholder(2), dbe.GetValue("id"))
			if dbr.Verbose {
				log.Print("DBRepository::deleteWithTx: Soft delete query=", query)
			}

			_, err = tx.Exec(query, dbe.GetValue("deleted_date"), dbe.GetValue("deleted_by"))
			if err != nil {
				log.Print("DBRepository::deleteWithTx: Exec error:", err)
				return nil, err
			}
			err = dbe.afterDelete(dbr, tx)
			if err != nil {
				log.Print("DBRepository::deleteWithTx: afterDelete error:", err)
				return nil, err
			}
			return dbe, nil
		}
		// If deleted_date is set, proceed with hard delete below
	}

	err := dbe.beforeDelete(dbr, tx)
	if err != nil {
		log.Print("DBRepository::deleteWithTx: beforeDelete error:", err)
		return nil, err
	}

	// 1. Build DELETE query dynamically based on primary keys
	whereClauses := make([]string, 0)
	args := make([]interface{}, 0)

	for _, key := range dbe.GetKeys() {
		value := dbe.GetValue(key)
		whereClauses = append(whereClauses, key+" = "+dbr.placeholder(len(args)+1))
		args = append(args, value)
	}

	if len(whereClauses) == 0 {
		return nil, fmt.Errorf("no primary keys defined for delete")
	}

	// 2. Build the DELETE query
	query := fmt.Sprintf("DELETE FROM %s WHERE %s",
		dbr.buildTableName(dbe),
		strings.Join(whereClauses, " AND "))

	if dbr.Verbose {
		log.Print("DBRepository::deleteWithTx: query=", query, " args=", args)
	}

	// 3. Execute the DELETE using the transaction
	result, err := tx.Exec(query, args...)
	if err != nil {
		log.Print("DBRepository::deleteWithTx: Exec error:", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && dbr.Verbose {
		log.Printf("DBRepository::deleteWithTx: rows affected=%d", rowsAffected)
	}

	err = dbe.afterDelete(dbr, tx)
	if err != nil {
		log.Print("DBRepository::deleteWithTx: afterDelete error:", err)
		return nil, err
	}

	return dbe, nil
}

// Update updates an existing entity in the database within a transaction
func (dbr *DBRepository) Update(dbe DBEntityInterface) (DBEntityInterface, error) {
	// Start a transaction
	tx, err := dbr.DbConnection.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Use internal method with transaction
	result, err := dbr.updateWithTx(dbe, tx)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

// updateWithTx is an internal method that performs the update using an existing transaction
func (dbr *DBRepository) updateWithTx(dbe DBEntityInterface, tx *sql.Tx) (DBEntityInterface, error) {
	if dbr.Verbose {
		log.Print("DBRepository::updateWithTx: dbe=", dbe.ToString())
	}

	// Call beforeUpdate hook (which can use dbr methods for nested operations)
	err := dbe.beforeUpdate(dbr, tx)
	if err != nil {
		log.Print("DBRepository::updateWithTx: beforeUpdate error:", err)
		return nil, err
	}

	// 1. Build UPDATE query dynamically based on populated fields (excluding primary keys)
	setClauses := make([]string, 0)
	args := make([]interface{}, 0)

	// Get primary keys to exclude them from SET clause
	primaryKeys := make(map[string]bool)
	for _, key := range dbe.GetKeys() {
		primaryKeys[key] = true
	}

	// Build SET clause with non-primary-key fields
	for key, value := range dbe.getDictionary() {
		if !primaryKeys[key] {
			setClauses = append(setClauses, key+" = "+dbr.placeholder(len(args)+1))
			args = append(args, value)
		}
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// 2. Build WHERE clause with primary keys
	whereClauses := make([]string, 0)
	for _, key := range dbe.GetKeys() {
		value := dbe.GetValue(key)
		whereClauses = append(whereClauses, key+" = "+dbr.placeholder(len(args)+1))
		args = append(args, value)
	}

	if len(whereClauses) == 0 {
		return nil, fmt.Errorf("no primary keys defined for update")
	}

	// 3. Build the UPDATE query
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		dbr.buildTableName(dbe),
		strings.Join(setClauses, ", "),
		strings.Join(whereClauses, " AND "))

	if dbr.Verbose {
		log.Print("DBRepository::updateWithTx: query=", query, " args=", args)
	}

	// 4. Execute the UPDATE using the transaction
	result, err := tx.Exec(query, args...)
	if err != nil {
		log.Print("DBRepository::updateWithTx: Exec error:", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && dbr.Verbose {
		log.Printf("DBRepository::updateWithTx: rows affected=%d", rowsAffected)
	}

	// Call afterUpdate hook
	err = dbe.afterUpdate(dbr, tx)
	if err != nil {
		log.Print("DBRepository::updateWithTx: afterUpdate error:", err)
		return nil, err
	}

	return dbe, nil
}

func (dbr *DBRepository) ExecuteSQL(sqlString string, args ...interface{}) (sql.Result, error) {
	if dbr.Verbose {
		log.Print("DBRepository::ExecuteSQL: sqlString=", sqlString, " args=", args)
	}
	result, err := dbr.DbConnection.Exec(sqlString, args...)
	if err != nil {
		log.Print("DBRepository::ExecuteSQL: Exec error:", err)
		return nil, err
	}
	return result, nil
}

func (dbr *DBRepository) GetCurrentUser() DBEntityInterface {
	if dbr.currentUser != nil {
		return dbr.currentUser
	}
	user := dbr.GetInstanceByTableName("users")
	if user == nil {
		return nil
	}
	user.SetValue("id", dbr.DbContext.UserID)
	foundUsers, err := dbr.Search(user, false, false, "")
	if err != nil || len(foundUsers) == 0 {
		return nil
	}
	dbr.currentUser = foundUsers[0].(*DBUser)
	return foundUsers[0]
}

// **** Objects Management ****

func (dbr *DBRepository) ObjectByID(objectID string, ignoreDeleted bool) DBEntityInterface {
	registeredTypes := dbr.factory.GetAllClassNames()
	var queries []string

	for _, className := range registeredTypes {
		// TODO: enable this once we have subclasses
		// if className == "DBObject" {
		// 	continue
		// }
		dbe := dbr.GetInstanceByClassName(className)
		if dbe == nil {
			continue
		}
		if !dbe.IsDBObject() {
			continue
		}
		query := "SELECT '" + className + "' as classname, id,owner,group_id,permissions,creator," +
			"creation_date,last_modify,last_modify_date," +
			"deleted_by,deleted_date," +
			"father_id,name,description" +
			" from " + dbr.buildTableName(dbe) +
			" WHERE id = '" + objectID + "'"
		if ignoreDeleted {
			query += " AND deleted_date IS NULL"
		}
		queries = append(queries, query)
	}
	searchString := strings.Join(queries, " UNION ")
	if dbr.Verbose {
		log.Print("DBRepository::ObjectByID: searchString=", searchString)
	}
	results := dbr.Select("DBObject", searchString)
	if len(results) == 0 {
		return nil
	}
	if len(results) > 1 {
		log.Printf("DBRepository::ObjectByID: Warning, multiple objects found with ID=%s", objectID)
	}
	return results[0]
}
func (dbr *DBRepository) FullObjectById(objectID string, ignoreDeleted bool) DBEntityInterface {
	obj := dbr.ObjectByID(objectID, ignoreDeleted)
	if obj == nil {
		return nil
	}
	dbObj, ok := obj.(*DBObject)
	if !ok {
		return nil
	}

	classname, ok := dbObj.GetMetadata("classname").(string)
	if !ok {
		return nil
	}
	fullObj := dbr.GetInstanceByClassName(classname)
	if fullObj == nil {
		return nil
	}
	fullObj.SetValue("id", objectID)
	foundEntities, err := dbr.Search(fullObj, false, false, "")
	if err != nil || len(foundEntities) == 0 {
		return nil
	}
	return foundEntities[0]
}
func (dbr *DBRepository) SearchByName(name string, orderBy string, ignoreDeleted bool) []DBEntityInterface {
	registeredTypes := dbr.factory.GetAllClassNames()
	var queries []string

	for _, className := range registeredTypes {
		dbe := dbr.GetInstanceByClassName(className)
		if dbe == nil {
			continue
		}
		if !dbe.IsDBObject() {
			continue
		}
		query := "SELECT '" + className + "' as classname, id,owner,group_id,permissions,creator," +
			"creation_date,last_modify,last_modify_date," +
			"deleted_by,deleted_date," +
			"father_id,name,description" +
			" from " + dbr.buildTableName(dbe) +
			" WHERE name like '%" + name + "%'"
		if ignoreDeleted {
			query += " AND deleted_date IS NULL"
		}
		queries = append(queries, query)
	}
	searchString := strings.Join(queries, " UNION ")
	if orderBy != "" {
		searchString += " ORDER BY " + orderBy
	}
	if dbr.Verbose {
		log.Print("DBRepository::SearchByName: searchString=", searchString)
	}
	results := dbr.Select("DBObject", searchString)
	return results
}

// SearchByNameAndDescription searches for DBObjects by name or description
// Returns all objects where name OR description contains the search text
func (dbr *DBRepository) SearchByNameAndDescription(searchText string, orderBy string, ignoreDeleted bool) []DBEntityInterface {
	registeredTypes := dbr.factory.GetAllClassNames()
	var queries []string

	for _, className := range registeredTypes {
		dbe := dbr.GetInstanceByClassName(className)
		if dbe == nil {
			continue
		}
		if !dbe.IsDBObject() {
			continue
		}
		query := "SELECT '" + className + "' as classname, id,owner,group_id,permissions,creator," +
			"creation_date,last_modify,last_modify_date," +
			"deleted_by,deleted_date," +
			"father_id,name,description" +
			" from " + dbr.buildTableName(dbe) +
			" WHERE ( name like '%" + searchText + "%' OR description like '%" + searchText + "%' ) "
		if ignoreDeleted {
			query += " AND deleted_date IS NULL"
		}
		queries = append(queries, query)
	}
	searchString := strings.Join(queries, " UNION ")
	if orderBy != "" {
		searchString += " ORDER BY " + orderBy
	}
	if dbr.Verbose {
		log.Print("DBRepository::SearchByNameAndDescription: searchString=", searchString)
	}
	results := dbr.Select("DBObject", searchString)
	return results
}

// GetChildren returns all direct children of a folder (objects with father_id = parentID)
// Filters results by read permissions
func (dbr *DBRepository) GetChildren(parentID string, ignoreDeleted bool) []DBEntityInterface {

	// Get the container object
	container := dbr.FullObjectById(parentID, true)
	if dbr.Verbose {
		log.Print("DBRepository.GetChildren: container=", container.ToJSON())
	}
	// Get childs_sort_order if container is DBFolder
	childs_sort_order := []string{}
	if container != nil && container.GetTypeName() == "DBFolder" {
		childs_sort_order = container.(*DBFolder).GetChildsSortOrder()
		if dbr.Verbose {
			log.Print("DBRepository.GetChildren: childs_sort_order=", childs_sort_order)
		}
	}

	registeredTypes := dbr.factory.GetAllClassNames()
	var queries []string

	for _, className := range registeredTypes {
		dbe := dbr.GetInstanceByClassName(className)
		if dbe == nil {
			continue
		}
		if !dbe.IsDBObject() {
			continue
		}

		clause := "father_id='" + parentID + "'"

		switch className {
		// case "DBCompany":
		// case "DBObject":
		// 	continue
		case "DBPerson":
			clause = "father_id='" + parentID + "'" + " or " + "fk_companies_id='" + parentID + "'"
		default:
			clause = "father_id='" + parentID + "'"
			// clause = "father_id='" + parentID + "' OR fk_obj_id='" + parentID + "'"
		}

		query := "SELECT '" + className + "' as classname, id,owner,group_id,permissions,creator," +
			"creation_date,last_modify,last_modify_date," +
			"deleted_by,deleted_date," +
			"father_id,name,description" +
			" from " + dbr.buildTableName(dbe) +
			" WHERE " + clause
		if ignoreDeleted {
			query += " AND deleted_date IS NULL"
		}
		queries = append(queries, query)
	}
	searchString := strings.Join(queries, " UNION ")
	searchString += " ORDER BY name"

	if dbr.Verbose {
		log.Print("DBRepository::GetChildren: searchString=", searchString)
	}
	results := dbr.Select("DBObject", searchString)

	// If childs_sort_order is defined, sort results accordingly and append any missing items at the end
	if len(childs_sort_order) > 0 {
		sortedResults := make([]DBEntityInterface, 0)
		seenIDs := make(map[string]bool)
		// First, add items in the order defined by childs_sort_order
		for _, childID := range childs_sort_order {
			for _, obj := range results {
				if obj.GetValue("id") == childID {
					sortedResults = append(sortedResults, obj)
					seenIDs[childID] = true
					break
				}
			}
		}
		// Then, add any remaining items that were not in childs_sort_order
		for _, obj := range results {
			objID := obj.GetValue("id").(string)
			if !seenIDs[objID] {
				sortedResults = append(sortedResults, obj)
			}
		}
		results = sortedResults
	}

	// Filter by read permissions
	return dbr.FilterByReadPermission(results)
}

// GetBreadcrumb returns the path from root to the specified object
// Each element is a DBObject with id, name, and father_id
func (dbr *DBRepository) GetBreadcrumb(objectID string, ignoreDeleted bool) []DBEntityInterface {
	breadcrumb := make([]DBEntityInterface, 0)
	currentID := objectID

	for currentID != "" && currentID != "0" {
		obj := dbr.ObjectByID(currentID, ignoreDeleted)
		if obj == nil {
			break
		}

		// Check read permission
		if !dbr.CheckReadPermission(obj) {
			break
		}

		breadcrumb = append(breadcrumb, obj)

		// Get father_id for next iteration
		fatherID := obj.GetValue("father_id")
		if fatherID == nil {
			break
		}
		fatherIDStr, ok := fatherID.(string)
		if !ok || fatherIDStr == "" {
			break
		}
		currentID = fatherIDStr
	}

	// Reverse the slice to get root -> object order
	for i := 0; i < len(breadcrumb)/2; i++ {
		j := len(breadcrumb) - 1 - i
		breadcrumb[i], breadcrumb[j] = breadcrumb[j], breadcrumb[i]
	}

	return breadcrumb
}

func (dbr *DBRepository) Select(returnedClassName string, sqlString string, args ...interface{}) []DBEntityInterface {
	if dbr.Verbose {
		log.Print("DBRepository::Select: sqlString=", sqlString, " args=", args)
	}
	rows, err := dbr.DbConnection.Query(sqlString, args...)
	if err != nil {
		log.Print("DBRepository::Select: Query error:", err)
		return nil
	}
	defer rows.Close()

	results := make([]DBEntityInterface, 0)
	columns, err := rows.Columns()
	if err != nil {
		log.Print("DBRepository::Select: Columns error:", err)
		return nil
	}

	for rows.Next() {
		// Read classname first
		// var className string
		columnValues := make([]interface{}, len(columns))
		columnValuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			columnValuePtrs[i] = &columnValues[i]
		}
		if err := rows.Scan(columnValuePtrs...); err != nil {
			log.Print("DBRepository::Select: Scan error:", err)
			return nil
		}
		// for i, colName := range columns {
		// 	if colName == "classname" {
		// 		if b, ok := columnValues[i].([]byte); ok {
		// 			className = string(b)
		// 		} else if columnValues[i] != nil {
		// 			className = fmt.Sprint(columnValues[i])
		// 		}
		// 		break
		// 	}
		// }
		dbe := dbr.GetInstanceByClassName(returnedClassName)
		if dbe == nil {
			log.Printf("DBRepository::Select: Warning, cannot create instance of class %s", returnedClassName)
			continue
		}
		// Map column values to the result entity's dictionary
		for i, colName := range columns {
			if colName == "classname" {
				myval := columnValues[i]
				if b, ok := myval.([]byte); ok {
					dbe.SetMetadata("classname", string(b))
				} else if myval != nil {
					dbe.SetMetadata("classname", fmt.Sprint(myval))
				}
				continue
			}
			// if colName == "deleted_date" {
			// 	// Special handling for deleted_date to allow nil
			// 	log.Println("DBRepository::Select: deleted_date=", columnValues[i])
			// }
			val := columnValues[i]
			if val == nil {
				// dbe.SetValue(colName, nil)
				continue
			}
			if b, ok := val.([]byte); ok {
				dbe.SetValue(colName, string(b))
			} else if val != nil {
				dbe.SetValue(colName, fmt.Sprint(val))
			} else {
				dbe.SetValue(colName, "")
			}
		}
		results = append(results, dbe)
	}

	if dbr.Verbose {
		log.Printf("DBRepository::Select: found %d results", len(results))
	}

	return results
}

// CheckReadPermission checks if the current user can read a DBObject
// Returns true if:
// - User is the owner
// - User is in the object's group and group has read permission
// - Object has public read permission
func (dbr *DBRepository) CheckReadPermission(dbe DBEntityInterface) bool {
	if !dbe.IsDBObject() {
		return true // Non-DBObjects have no permission restrictions
	}

	owner, ok := dbe.GetValue("owner").(string)
	if !ok {
		return false
	}

	// User is owner
	if dbr.DbContext.IsUser(owner) {
		permissions, ok := dbe.GetValue("permissions").(string)
		if !ok || len(permissions) != 9 {
			return false
		}
		return permissions[0] == 'r' // User read permission
	}

	groupID, ok := dbe.GetValue("group_id").(string)
	if !ok {
		return false
	}

	// User is in group
	if dbr.DbContext.IsInGroup(groupID) {
		permissions, ok := dbe.GetValue("permissions").(string)
		if !ok || len(permissions) != 9 {
			return false
		}
		return permissions[3] == 'r' // Group read permission
	}

	// Check public read permission
	permissions, ok := dbe.GetValue("permissions").(string)
	if !ok || len(permissions) != 9 {
		return false
	}
	return permissions[6] == 'r' // Others read permission
}

// CheckWritePermission checks if the current user can write (modify/delete) a DBObject
// Returns true if:
// - User is the owner and has write permission
// - User is in the object's group and group has write permission
// - Object has public write permission
func (dbr *DBRepository) CheckWritePermission(dbe DBEntityInterface) bool {
	if !dbe.IsDBObject() {
		return true // Non-DBObjects have no permission restrictions
	}

	owner, ok := dbe.GetValue("owner").(string)
	if !ok {
		return false
	}

	// User is owner
	if dbr.DbContext.IsUser(owner) {
		permissions, ok := dbe.GetValue("permissions").(string)
		if !ok || len(permissions) != 9 {
			return false
		}
		return permissions[1] == 'w' // User write permission
	}

	groupID, ok := dbe.GetValue("group_id").(string)
	if !ok {
		return false
	}

	// User is in group
	if dbr.DbContext.IsInGroup(groupID) {
		permissions, ok := dbe.GetValue("permissions").(string)
		if !ok || len(permissions) != 9 {
			return false
		}
		return permissions[4] == 'w' // Group write permission
	}

	// Check public write permission
	permissions, ok := dbe.GetValue("permissions").(string)
	if !ok || len(permissions) != 9 {
		return false
	}
	return permissions[7] == 'w' // Others write permission
}

// FilterByReadPermission filters a slice of DBEntityInterface, keeping only objects the user can read
func (dbr *DBRepository) FilterByReadPermission(entities []DBEntityInterface) []DBEntityInterface {
	filtered := make([]DBEntityInterface, 0, len(entities))
	for _, entity := range entities {
		if dbr.CheckReadPermission(entity) {
			filtered = append(filtered, entity)
		}
	}
	return filtered
}

func (dbr *DBRepository) GetEntityByID(tableName string, id string) DBEntityInterface {
	return dbr.GetEntityByIDWithTx(tableName, id, nil)
}

// GetEntityByID retrieves a generic entity (non-DBObject) by table name and ID
func (dbr *DBRepository) GetEntityByIDWithTx(tableName string, id string, tx *sql.Tx) DBEntityInterface {
	dbe := dbr.GetInstanceByTableName(tableName)
	if dbe == nil {
		return nil
	}
	dbe.SetValue("id", id)
	results, err := dbr.searchWithTx(dbe, false, false, "", tx)
	if err != nil || len(results) == 0 {
		return nil
	}
	return results[0]
}
