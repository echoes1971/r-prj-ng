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
	Verbose   bool
	DbContext *DBContext
	factory   *DBEFactory

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
	if dbr.DbContext != nil && dbr.DbContext.Schema != "" {
		return dbr.DbContext.Schema + "_" + dbe.GetTableName()
	}
	return dbe.GetTableName()
}

func (dbr *DBRepository) Search(dbe DBEntityInterface, useLike bool, caseSensitive bool, orderBy string) ([]DBEntityInterface, error) {
	return dbr.searchWithTx(dbe, useLike, caseSensitive, orderBy, nil)
}

// searchWithTx is an internal method that performs the search using an existing transaction (if provided)
func (dbr *DBRepository) searchWithTx(dbe DBEntityInterface, useLike bool, caseSensitive bool, orderBy string, tx *sql.Tx) ([]DBEntityInterface, error) {
	if dbr.Verbose {
		log.Print("DBRepository::searchWithTx: dbe=", dbe)
	}

	// 1. Build WHERE clauses
	whereClauses := make([]string, 0)
	args := make([]interface{}, 0) // slice of interface{} for values

	for key, value := range dbe.getDictionary() {
		if useLike {
			// For strings: LIKE '%value%'
			if strings.Contains(dbe.GetColumnType(key), "varchar") || dbe.GetColumnType(key) == "text" {
				if caseSensitive {
					whereClauses = append(whereClauses, key+" LIKE ?")
					args = append(args, "%"+fmt.Sprint(value)+"%")
				} else {
					whereClauses = append(whereClauses, "LOWER("+key+") LIKE LOWER(?)")
					args = append(args, "%"+fmt.Sprint(value)+"%")
				}
			} else {
				// Per numeri/date: exact match
				whereClauses = append(whereClauses, key+" = ?")
				args = append(args, value)
			}
		} else {
			// Exact match
			whereClauses = append(whereClauses, key+" = ?")
			args = append(args, value)
		}
	}

	// 2. Build the final query
	query := "SELECT * FROM " + dbr.buildTableName(dbe)
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
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
			} else {
				resultEntity.SetValue(colName, "")
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

	for key, value := range dbe.getDictionary() {
		columns = append(columns, key)
		placeholders = append(placeholders, "?")
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

	if dbr.Verbose {
		log.Print("DBRepository::insertWithTx: query=", query, " args=", args)
	}

	// 3. Execute the INSERT using the transaction
	result, err := tx.Exec(query, args...)
	if err != nil {
		log.Print("DBRepository::insertWithTx: Exec error:", err)
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
		whereClauses = append(whereClauses, key+" = ?")
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
			setClauses = append(setClauses, key+" = ?")
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
		whereClauses = append(whereClauses, key+" = ?")
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
