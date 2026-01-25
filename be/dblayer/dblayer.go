package dblayer

import (
	"database/sql"
	"fmt"
	"log"
	"rprj/be/models"
	"slices"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// The database can be mysql, sqlite, postgres, etc.
var dbEngine string
var dbUrl string
var DbSchema string
var DbConnection *sql.DB
var Factory *DBEFactory

/* *** DBFiles *** */
var dbFiles_root_directory string = "."
var dbFiles_dest_directory string = "files"

func InitDBLayer(config models.Config) {
	dbEngine = config.DBEngine
	log.Print("DB Engine:", dbEngine)
	dbUrl = config.DBUrl
	DbSchema = strings.ReplaceAll(config.TablePrefix, "_", "")
	log.Print("DB Schema:", DbSchema)
	dbFiles_root_directory = config.RootDirectory
	dbFiles_dest_directory = config.FilesDirectory

	log.Print("Initializing DBEFactory...")

	Factory = NewDBEFactory(false)
	Factory.Register(NewDBVersion())
	Factory.Register(NewOAuthToken())
	Factory.Register(NewDBUser())
	Factory.Register(NewUserGroup())
	Factory.Register(NewDBGroup())
	Factory.Register(NewDBLog())
	Factory.Register(NewDBObject())
	// Contacts
	Factory.Register(NewDBCountry())
	Factory.Register(NewDBCompany())
	Factory.Register(NewDBPerson())
	// CMS
	Factory.Register(NewDBEvent())
	Factory.Register(NewDBFile())
	Factory.Register(NewDBFolder())
	Factory.Register(NewDBLink())
	Factory.Register(NewDBNote())
	Factory.Register(NewDBPage())
	Factory.Register(NewDBNews())
	// Process foreign keys after all registrations
	Factory.ProcessForeignKeys()

	InitDBConnection()
	// log.Print("Initializing DB connection...")
	// var err error
	// DbConnection, err = sql.Open(dbEngine, dbUrl)
	// if err != nil {
	// 	log.Fatal("Error opening DB connection:", err)
	// }
	// err = DbConnection.Ping()
	// if err != nil {
	// 	log.Print("Error pinging DB:", err)
	// }

	// // TODO: make this configurable
	// // Configure connection pool to handle concurrent operations
	// DbConnection.SetMaxOpenConns(25)   // Maximum number of open connections to the database
	// DbConnection.SetMaxIdleConns(10)   // Maximum number of connections in the idle connection pool
	// DbConnection.SetConnMaxLifetime(0) // Maximum amount of time a connection may be reused (0 = unlimited)

}

func InitDBConnection() {
	if DbConnection != nil {
		return
	}
	log.Print("Initializing DB connection...")
	var err error
	// DbConnection, err = sql.Open(dbEngine, dbUrl)
	// if err != nil {
	// 	log.Print(" InitDBConnection: Error opening DB connection:", err)
	// }

	log.Print(" DB Engine: ", dbEngine, " DB URL:", dbUrl)
	dbName := ""
	switch dbEngine {
	case "mysql":
		parts := strings.Split(dbUrl, "/")
		if len(parts) > 1 {
			dbName = strings.Split(parts[1], "?")[0]
		}
	case "postgres":
		parts := strings.Split(dbUrl, "/")
		log.Print(" DB URL parts:", parts)
		if len(parts) > 1 {
			dbName = strings.Split(parts[len(parts)-1], "?")[0]
		}
	}
	log.Print(" DB Name:", dbName)
	dbUrlNoDB := dbUrl
	if dbName != "" {
		dbUrlNoDB = strings.ReplaceAll(dbUrl, "/"+dbName, "/")
	}
	log.Print(" DB URL without DB:", dbUrlNoDB)

	DbConnection, err = sql.Open(dbEngine, dbUrlNoDB)
	if err != nil {
		log.Fatal(" InitDBConnection: Error opening DB connection:", err)
	}

	// Create DB if not exists (for sqlite, the DB file is created automatically)
	sqlCreateDB := ""
	switch dbEngine {
	case "sqlite3":
	case "mysql":
		sqlCreateDB = "CREATE DATABASE IF NOT EXISTS " + dbName + ";"
	case "postgres":
		sqlCreateDB = "CREATE DATABASE " + dbName + ";"
	}
	if sqlCreateDB != "" {
		log.Print("Creating DB if not exists with SQL:", sqlCreateDB)
		_, err = DbConnection.Exec(sqlCreateDB)
		if err != nil {
			log.Print(" Error creating DB:", err)
		}
		// Close and reopen connection to the specific DB
		DbConnection.Close()
		DbConnection, err = sql.Open(dbEngine, dbUrl)
		if err != nil {
			log.Fatal(" InitDBConnection: Error reopening DB connection:", err)
		}
	}

	err = DbConnection.Ping()
	if err != nil {
		log.Fatal("Error pinging DB:", err)
	}

	// TODO: make this configurable
	// Configure connection pool to handle concurrent operations
	DbConnection.SetMaxOpenConns(25)   // Maximum number of open connections to the database
	DbConnection.SetMaxIdleConns(10)   // Maximum number of connections in the idle connection pool
	DbConnection.SetConnMaxLifetime(0) // Maximum amount of time a connection may be reused (0 = unlimited)
}

var objectsColumns = []string{"id", "owner", "group_id", "permissions", "creator", "creation_date", "last_modify", "last_modify_date", "deleted_by", "deleted_date", "father_id", "name", "description"}

// GetCreateTableSQL generates CREATE TABLE SQL for the given entity
// This is a standalone function to properly use polymorphism with IsDBObject()
func GetCreateTableSQL(dbe DBEntityInterface, dbSchema string) string {
	columnDefs := []string{}

	isDBObject := dbe.IsDBObject()
	log.Print("isDBObject=", isDBObject, " typename=", dbe.GetTypeName())
	isDBObjectChild := isDBObject && dbe.GetTypeName() != "DBObject"

	for _, col := range dbe.GetColumnDefinitions() {
		if dbEngine == "postgres" && isDBObjectChild && slices.Contains(objectsColumns, col.Name) {
			continue
		}
		colDef := fmt.Sprintf(" %s %s", col.Name, col.Type)
		if len(col.Constraints) > 0 {
			colDef += " " + strings.Join(col.Constraints, " ")
		}
		columnDefs = append(columnDefs, colDef)
	}
	// Add primary key constraint
	if len(dbe.GetKeys()) > 0 {
		pkDef := fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(dbe.GetKeys(), ", "))
		columnDefs = append(columnDefs, pkDef)
	}
	// TODO: is it worth it as I have fields pointing to multiple tables?
	// // Add foreign key constraints
	// for _, fk := range dbEntity.foreignKeys {
	// 	// TODO: remove this exception when implementing project entities
	// 	notExistingTables := []string{"projects", "tasks"}
	// 	if slices.Contains(notExistingTables, fk.RefTable) {
	// 		continue
	// 	}
	// 	// TODO: with Postgresql, we should skip FK to objects table
	// 	if fk.RefTable == "objects" {
	// 		continue
	// 	}

	// 	fkDef := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s_%s(%s)", fk.Column, dbSchema, fk.RefTable, fk.RefColumn)
	// 	columnDefs = append(columnDefs, fkDef)
	// }
	// IF NOT EXISTS is redundant as we check for table existence before calling this method: but it's kept for future use cases
	inheritanceClause := ""
	if dbEngine == "postgres" && isDBObjectChild {
		inheritanceClause = fmt.Sprintf(" INHERITS (%s_%s)", dbSchema, "objects")
	}
	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s_%s (\n%s\n)%s;", dbSchema, dbe.GetTableName(), strings.Join(columnDefs, ",\n"), inheritanceClause)
	return createTableSQL
}

// The database can be mysql, sqlite, postgres, etc.
func ensureTableExistsAndUpdatedForMysql(dbe DBEntityInterface, Verbose bool) error {
	// Check if table exists
	tableName := DbSchema + "_" + dbe.GetTableName()
	if Verbose {
		log.Print("Checking table: ", tableName)
	}
	// This is internal code, so we can build the query directly
	query := "show tables like '" + tableName + "'"
	var existingTable string
	err := DbConnection.QueryRow(query).Scan(&existingTable)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking existence of table %s: %v", tableName, err)
		return err
	}

	if existingTable == "" {
		// Table does not exist, create it
		// Compose the create table SQL using the DBEntity's information about columns, types, keys, etc.
		createTableSQL := GetCreateTableSQL(dbe, DbSchema)
		if Verbose {
			log.Printf(" Creating table with SQL: %s", createTableSQL)
		}
		_, err := DbConnection.Exec(createTableSQL)
		if err != nil {
			return err
		}
		if Verbose {
			log.Printf(" Created table %s", tableName)
		}
	} else {
		// Table exists, check for schema updates
		// For simplicity, we will not implement schema migration logic here
		if Verbose {
			log.Printf(" Table %s exists", tableName)
		}
		// Fetch table schema from MariaDB and compare with DBEntity definition
		rows, err := DbConnection.Query("DESCRIBE " + tableName)
		if err != nil {
			log.Printf("Error describing table %s: %v", tableName, err)
			return err
		}
		defer rows.Close()
		columnsInDB := make(map[string]map[string]string)
		for rows.Next() {
			var field, colType, null, key, extra string
			var defaultValue any
			if err := rows.Scan(&field, &colType, &null, &key, &defaultValue, &extra); err != nil {
				log.Printf("Error scanning row for table %s: %v", tableName, err)
				return err
			}
			// log.Print("field=", field, " colType=", colType, " null=", null, " key=", key, " defaultValue=", defaultValue, " extra=", extra)
			// convert defaultValue to string
			defaultValueStr := ""
			if defaultValue != nil {
				defaultValueStr = fmt.Sprintf("%v", defaultValue)
			}
			columnsInDB[field] = map[string]string{
				"Type":    colType,
				"Null":    null,
				"Key":     key,
				"Default": defaultValueStr,
				"Extra":   extra,
			}
		}

		// Compare columnsInDB with dbe.GetColumnDefinitions() and identify differences
		columnDefs := dbe.GetColumnDefinitions()
		for colName, colDef := range columnDefs {
			if dbColDef, exists := columnsInDB[colName]; exists {
				// Column exists, check for differences
				if !strings.EqualFold(dbColDef["Type"], colDef.Type) {
					log.Printf(" Column %s type mismatch: DB=%s, Expected=%s", colName, dbColDef["Type"], colDef.Type)
					alterTableSQL := fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s", tableName, colName, colDef.Type)
					log.Printf(" Alter column with: %s", alterTableSQL)
					// Implement ALTER TABLE to modify column type if needed
				}
				// Check other attributes as needed (Null, Key, Default, Extra)
			} else {
				// Column does not exist, add it
				addColumnSQL := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, colName, colDef.Type)
				log.Printf("Adding missing column with SQL: %s", addColumnSQL)
				// _, err := DbConnection.Exec(addColumnSQL)
				// if err != nil {
				// 	log.Printf("Error adding column %s to table %s: %v", colName, tableName, err)
				// 	return err
				// }
				log.Printf("Added column %s to table %s", colName, tableName)
			}
		}
		// Implement schema comparison and migration logic as needed
		// This can be complex and is often handled by dedicated migration tools

	}

	return nil
}
func ensureTableExistsAndUpdatedForSqlite(dbe DBEntityInterface, Verbose bool) error {
	// Check if table exists
	tableName := DbSchema + "_" + dbe.GetTableName()
	if Verbose {
		log.Print("Checking table: ", tableName)
	}
	// Query SQLite's metadata table
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
	var existingTable string
	err := DbConnection.QueryRow(query, tableName).Scan(&existingTable)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking existence of table %s: %v", tableName, err)
		return err
	}

	if existingTable == "" {
		// Table does not exist, create it
		createTableSQL := GetCreateTableSQL(dbe, DbSchema)
		if Verbose {
			log.Printf(" Creating table with SQL: %s", createTableSQL)
		}
		_, err := DbConnection.Exec(createTableSQL)
		if err != nil {
			return err
		}
		if Verbose {
			log.Printf(" Created table %s", tableName)
		}
	} else {
		// Table exists, check for schema updates
		if Verbose {
			log.Printf(" Table %s exists", tableName)
		}
		// Fetch table schema from SQLite using PRAGMA table_info
		rows, err := DbConnection.Query("PRAGMA table_info(" + tableName + ")")
		if err != nil {
			log.Printf("Error describing table %s: %v", tableName, err)
			return err
		}
		defer rows.Close()
		columnsInDB := make(map[string]map[string]string)
		for rows.Next() {
			var cid int
			var name, colType string
			var notNull, pk int
			var dfltValue sql.NullString
			if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
				log.Printf("Error scanning row for table %s: %v", tableName, err)
				return err
			}
			defaultValueStr := ""
			if dfltValue.Valid {
				defaultValueStr = dfltValue.String
			}
			columnsInDB[name] = map[string]string{
				"Type":    colType,
				"NotNull": fmt.Sprintf("%d", notNull),
				"Default": defaultValueStr,
				"PK":      fmt.Sprintf("%d", pk),
			}
		}

		// Compare columnsInDB with dbe.GetColumnDefinitions() and identify differences
		columnDefs := dbe.GetColumnDefinitions()
		for colName, colDef := range columnDefs {
			if dbColDef, exists := columnsInDB[colName]; exists {
				// Column exists, check for differences
				if !strings.EqualFold(dbColDef["Type"], colDef.Type) {
					log.Printf(" Column %s type mismatch: DB=%s, Expected=%s", colName, dbColDef["Type"], colDef.Type)
					// Note: SQLite doesn't support MODIFY COLUMN, would need table recreation
					log.Printf(" Warning: SQLite doesn't support ALTER COLUMN TYPE, manual migration needed")
				}
			} else {
				// Column does not exist, add it
				addColumnSQL := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, colName, colDef.Type)
				log.Printf("Adding missing column with SQL: %s", addColumnSQL)
				// _, err := DbConnection.Exec(addColumnSQL)
				// if err != nil {
				// 	log.Printf("Error adding column %s to table %s: %v", colName, tableName, err)
				// 	return err
				// }
				log.Printf("Added column %s to table %s", colName, tableName)
			}
		}
	}

	return nil
}
func ensureTableExistsAndUpdatedForPostgres(dbe DBEntityInterface, Verbose bool) error {
	// Check if table exists
	tableName := DbSchema + "_" + dbe.GetTableName()
	if Verbose {
		log.Print("Checking table: ", tableName)
	}
	// This is internal code, so we can build the query directly
	query := "SELECT to_regclass('" + tableName + "')"
	var existingTable sql.NullString
	err := DbConnection.QueryRow(query).Scan(&existingTable)
	if err != nil {
		log.Printf("Error checking existence of table %s: %v", tableName, err)
		return err
	}

	if !existingTable.Valid {
		// Table does not exist, create it
		// Compose the create table SQL using the DBEntity's information about columns, types, keys, etc.
		createTableSQL := GetCreateTableSQL(dbe, DbSchema)
		// replace datetime with timestamp in createTableSQL for Postgres
		createTableSQL = strings.ReplaceAll(createTableSQL, " DATETIME", " TIMESTAMP")
		createTableSQL = strings.ReplaceAll(createTableSQL, " datetime", " TIMESTAMP")
		log.Printf(" Creating table with SQL: %s", createTableSQL)
		_, err := DbConnection.Exec(createTableSQL)
		if err != nil {
			return err
		}
		log.Printf(" Created table %s", tableName)
	} else {
		// Table exists, check for schema updates
		// For simplicity, we will not implement schema migration logic here
		log.Printf(" Table %s exists", tableName)
		// Fetch table schema from Postgres and compare with DBEntity definition
		// err := DbConnection.QueryRow("SELECT to_regclass('" + tableName + "')").Scan(&existingTable)
		// if err != nil {
		// 	log.Printf("Error describing table %s: %v", tableName, err)
		// 	return err
		// }
		rows, err := DbConnection.Query("SELECT column_name, data_type, is_nullable, column_default, character_maximum_length FROM information_schema.columns WHERE table_name = $1", tableName)
		if err != nil {
			log.Printf("Error describing table %s: %v", tableName, err)
			return err
		}

		defer rows.Close()
		columnsInDB := make(map[string]map[string]string)
		for rows.Next() {
			var columnName, dataType, isNullable string
			var columnDefault sql.NullString
			var charMaxLength sql.NullInt64
			if err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault, &charMaxLength); err != nil {
				log.Printf("Error scanning row for table %s: %v", tableName, err)
				return err
			}
			// log.Print("columnName=", columnName, " dataType=", dataType, " isNullable=", isNullable, " columnDefault=", columnDefault, " charMaxLength=", charMaxLength)
			defaultValueStr := ""
			if columnDefault.Valid {
				defaultValueStr = columnDefault.String
			}
			columnsInDB[columnName] = map[string]string{
				"Type":    dataType,
				"Null":    isNullable,
				"Default": defaultValueStr,
				"MaxLen":  fmt.Sprintf("%d", charMaxLength.Int64),
			}
		}

		// Compare columnsInDB with dbe.GetColumnDefinitions() and identify differences
		columnDefs := dbe.GetColumnDefinitions()
		for colName, colDef := range columnDefs {
			colNameLowerCase := strings.ToLower(colName)
			if dbColDef, exists := columnsInDB[colNameLowerCase]; exists {
				// Column exists, check for differences
				db_type := dbColDef["Type"]
				db_type = strings.ToLower(db_type)
				db_type = strings.ReplaceAll(db_type, "character varying", "varchar")
				db_type = strings.ReplaceAll(db_type, "character", "char")
				db_type = strings.ReplaceAll(db_type, "timestamp without time zone", "datetime")
				db_type = strings.ReplaceAll(db_type, "timestamp", "datetime")
				db_type = strings.ReplaceAll(db_type, "time without time zone", "time")
				db_type = strings.ReplaceAll(db_type, "integer", "int")
				if MaxLen, ok := dbColDef["MaxLen"]; ok && MaxLen != "0" {
					db_type += fmt.Sprintf("(%s)", MaxLen)
				}
				if !strings.EqualFold(db_type, colDef.Type) {
					log.Printf(" Column %s type mismatch: DB=%s, Expected=%s", colName, dbColDef["Type"], colDef.Type)
					alterTableSQL := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", tableName, colName, colDef.Type)
					log.Printf(" Alter column with: %s", alterTableSQL)
					// Implement ALTER TABLE to modify column type if needed
				}
				// Check other attributes as needed (Null, Key, Default, Extra)
			} else {
				// Column does not exist, add it
				addColumnSQL := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, colName, colDef.Type)
				log.Printf("Adding missing column with SQL: %s", addColumnSQL)
				// _, err := DbConnection.Exec(addColumnSQL)
				// if err != nil {
				// 	log.Printf("Error adding column %s to table %s: %v", colName, tableName, err)
				// 	return err
				// }
				log.Printf("Added column %s to table %s", colName, tableName)
			}
		}
		// Implement schema comparison and migration logic as needed
		// This can be complex and is often handled by dedicated migration tools

	}

	return nil
}

func InitDBData() {
	log.Print("Initializing DB data...")

	// Setup
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   DbSchema,
	}
	log.Print("InitDBData: dbContext=", dbContext)

	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false

	// Load initial data from JSON
	initialData, err := LoadInitialData()
	if err != nil {
		log.Printf("Failed to load initial data: %v\n", err)
		return
	}

	// Import tables in order (preserves foreign key constraints)
	for _, table := range initialData.Tables {
		log.Printf("Populating table '%s' with %d rows...\n", table.Name, len(table.Data))
		results, err, shouldReturn := populateTable(repo, table.Name, table.Columns, table.Data)
		if shouldReturn {
			return
		}
		if err == nil && len(results) > 0 {
			log.Printf(" Successfully populated '%s'\n", table.Name)
		}
	}

	// DBVersion
	dbVersion := repo.GetInstanceByTableName("dbversion")
	results, err := repo.Search(dbVersion, false, false, "")
	if err != nil {
		log.Printf(" Failed to find or create DB version entry: %v\n", err)
		return
	}
	if len(results) == 0 {
		newVersion := repo.GetInstanceByTableName("dbversion")
		newVersion.SetValue("model_name", "rprj")
		newVersion.SetValue("version", 2)
		_, err := repo.Insert(newVersion)
		if err != nil {
			log.Printf(" Failed to create DB version entry: %v\n", err)
			return
		}
		log.Printf(" Created DB version entry.\n")
	} else {
		log.Printf(" DB version entry exists with version %s.\n", results[0].GetValue("version").(string))
	}

	log.Print("DB data initialization completed.")
}

func populateTable(repo *DBRepository, tablename string, listColumns []string, listData [][]string) ([]DBEntityInterface, error, bool) {
	dbe := repo.GetInstanceByTableName(tablename)
	results, err := repo.Search(dbe, false, false, "")
	if err != nil {
		log.Printf(" Failed to find or create %s: %v\n", tablename, err)
		return nil, nil, true
	}
	log.Println(" Found ", len(results), " existing entries in ", tablename)
	if len(results) == 0 {
		// Populate table
		log.Printf(" Populating %s...\n", tablename)
		for _, entryData := range listData {
			newEntry := repo.GetInstanceByTableName(tablename)
			for i, columnName := range listColumns {
				if entryData[i] == "" || entryData[i] == "NULL" {
					continue
				}
				newEntry.SetValue(columnName, entryData[i])
			}
			if tablename == "folders" {
				log.Println(" newEntry:", newEntry.ToJSON())
			}
			log.Printf(" Inserting entry: %s\n", entryData[0])
			_, err := repo.Insert(newEntry)
			if err != nil {
				for _, col := range listColumns {
					log.Printf("  %s = %v\n", col, newEntry.GetValue(col))
				}
				log.Printf(" Failed to insert data: %v\n", err)
				return nil, nil, true
			}
		}
	} else if tablename == "users_groups" && len(results) < len(listData) {
		for _, entryData := range listData {
			newEntry := repo.GetInstanceByTableName(tablename)
			for i, columnName := range listColumns {
				newEntry.SetValue(columnName, entryData[i])
			}
			log.Printf(" Inserting missing entry: %s\n", newEntry.ToJSON())
			_, err := repo.Insert(newEntry)
			if err != nil {
				log.Printf(" Failed to insert data: %v\n", err)
			}
		}
	} else if len(results) < len(listData) {
		// Populate missing entries
		log.Printf(" Populating missing entries in %s...\n", tablename)
		existingEntries := make(map[string]bool)
		for _, res := range results {
			name := res.GetValue("id").(string)
			existingEntries[name] = true
		}
		for _, entryData := range listData {
			entryID := entryData[0]
			if _, exists := existingEntries[entryID]; !exists {
				newEntry := repo.GetInstanceByTableName(tablename)
				for i, columnName := range listColumns {
					if entryData[i] == "" || entryData[i] == "NULL" {
						continue
					}
					newEntry.SetValue(columnName, entryData[i])
				}
				log.Printf(" Inserting missing entry: %s %s\n", entryID, entryData[1])
				_, err := repo.Insert(newEntry)
				if err != nil {
					for _, col := range listColumns {
						log.Printf("  %s = %v\n", col, newEntry.GetValue(col))
					}
					log.Printf(" Failed to insert data: %v\n", err)
					return nil, nil, true
				}
			}
		}
	}
	log.Printf(" %s initialized with %d entries.\n", tablename, len(listData))
	return results, err, false
}

// Iterate over all registered DBEntity types and create tables if they do not exist or update their schema
func EnsureDBSchema(Verbose bool) {
	var classInstances []DBEntityInterface
	for _, className := range Factory.GetAllClassNames() {
		dbe := Factory.GetInstanceByClassName(className)
		if dbe != nil {
			classInstances = append(classInstances, dbe)
		}
	}
	// Sort classInstances based on dependencies (foreign keys) if needed
	// Sort based on dependencies using a topological sort
	sorted := make([]DBEntityInterface, 0, len(classInstances))
	visited := make(map[string]bool)
	temp := make(map[string]bool)

	var visit func(DBEntityInterface, string) bool
	visit = func(dbe DBEntityInterface, prefix string) bool {
		className := dbe.GetTypeName()
		if temp[className] {
			// Circular dependency detected
			log.Printf(prefix+"- Warning: Circular dependency detected for %s", className)
			return false
		}
		if visited[className] {
			return true
		}

		temp[className] = true
		// Get foreign key dependencies
		foreignKeys := dbe.GetForeignKeys()
		currentTableName := dbe.GetTableName()
		for _, fk := range foreignKeys {
			// Skip self-referencing foreign keys (e.g., DBObject.parent_id -> DBObject.id)
			if fk.RefTable == currentTableName {
				log.Printf(prefix+"- Skipping self-referencing FK: %s.%s -> %s.%s", currentTableName, fk.Column, fk.RefTable, fk.RefColumn)
				continue
			}
			// Find the referenced table's DBEntity
			for _, depDbe := range classInstances {
				if depDbe.GetTableName() == fk.RefTable {
					// log.Printf(prefix+"- %s depends on %s via FK %s -> %s", className, depDbe.GetTypeName(), fk.Column, fk.RefColumn)
					if !visit(depDbe, prefix+"  ") {
						return false
					}
					break
				}
			}
		}
		temp[className] = false
		visited[className] = true
		sorted = append(sorted, dbe)
		return true
	}

	for _, dbe := range classInstances {
		// log.Print("Processing ", dbe.GetTypeName())
		if !visited[dbe.GetTypeName()] {
			// log.Printf("- Visiting %s for topological sort", dbe.GetTypeName())
			visit(dbe, "  ")
		}
		// visitedJSON, _ := json.Marshal(visited)
		// visitedJSON, _ := json.MarshalIndent(visited, "  ", "  ")
		// log.Print(" Visited:", string(visitedJSON))
	}
	// reverse sorted to get correct order
	// for i, j := 0, len(sorted)-1; i < j; i, j = i+1, j-1 {
	// 	sorted[i], sorted[j] = sorted[j], sorted[i]
	// }
	// slices.Reverse(sorted)

	// Print sorted class names for debugging
	if Verbose {
		log.Print("DB Entities creation order:")
		for _, dbe := range sorted {
			log.Printf(" - %s\n", dbe.GetTableName())
			for _, fk := range dbe.GetForeignKeys() {
				log.Printf("    FK: %s -> %s(%s)\n", fk.Column, fk.RefTable, fk.RefColumn)
			}
		}
	}
	classInstances = sorted

	for _, dbe := range classInstances {
		var err error
		className := dbe.GetTypeName()
		switch dbEngine {
		case "mysql":
			err = ensureTableExistsAndUpdatedForMysql(dbe, Verbose)
		case "sqlite3":
			err = ensureTableExistsAndUpdatedForSqlite(dbe, Verbose)
		case "postgres":
			err = ensureTableExistsAndUpdatedForPostgres(dbe, Verbose)
		default:
			log.Fatal("Unsupported dbEngine:", dbEngine)
		}
		if err != nil {
			log.Fatal("Error ensuring table for ", className, ":", err)
		}
	}
}
func CloseDBConnection() {
	if DbConnection != nil {
		log.Print("Closing DB connection...")
		DbConnection.Close()
	}
}

/* **** Compatibility functions **** */

func uuid2hex(str string) string {
	if str == "" {
		return str
	}
	if len(str) < 4 {
		return str
	}
	if str[0:4] == "uuid" {
		return str
	}
	hex := ""
	for i := 0; i < len(str); i++ {
		hex += stringFormat("%x", str[i])
	}
	return "uuid" + hex
}
func hex2uuid(a_str string) string {
	if len(a_str) < 4 || a_str[0:4] != "uuid" {
		return a_str
	}
	str := a_str[4:]
	bin := ""
	for i := 0; i < len(str); i += 2 {
		var b byte
		fmtSscanf(str[i:i+2], "%02x", &b)
		bin += string(b)
	}
	return bin
}
func stringFormat(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}
func fmtSscanf(str string, format string, a ...interface{}) {
	fmt.Sscanf(str, format, a...)
}

// Below is the equivalent PHP code for the uuid2hex and hex2uuid functions:
/*
   static function uuid2hex($str) {
       if($str===null) return $str;
       $str_len = strlen($str);
       if($str_len<4) return $str;
       if(substr($str,0,4)=='uuid')
           return $str;
       $hex = "";
       $i = 0;
       do {
           $hex .= dechex(ord($str[$i]));
           // $hex .= dechex(ord($str{$i}));
           $i++;
       } while ($i<$str_len);
       return 'uuid'.$hex;
   }
   static function hex2uuid($a_str) {
       if(substr($a_str,0,4)!='uuid')
           return $a_str;
       $str=substr($a_str,4);
       $bin = "";
       $i = 0;
       $str_len = strlen($str);
       do {
           $bin .= chr(hexdec($str[$i].$str[($i + 1)]));
           // $bin .= chr(hexdec($str{$i}.$str{($i + 1)}));
           $i += 2;
       } while ($i < $str_len);
       return $bin;
   }
*/
