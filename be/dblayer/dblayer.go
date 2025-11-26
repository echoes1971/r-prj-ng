package dblayer

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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

func InitDBLayer(dbEngineName, dbUrlAddress, schema string) {
	dbEngine = dbEngineName
	dbUrl = dbUrlAddress
	DbSchema = strings.ReplaceAll(schema, "_", "")
	log.Print("DB Schema:", DbSchema)

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

	log.Print("Initializing DB connection...")
	var err error
	DbConnection, err = sql.Open(dbEngine, dbUrl)
	if err != nil {
		log.Fatal("Error opening DB connection:", err)
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

func InitDBConnection() {
	if DbConnection != nil {
		return
	}
	log.Print("Initializing DB connection...")
	var err error
	DbConnection, err = sql.Open(dbEngine, dbUrl)
	if err != nil {
		log.Fatal("InitDBConnection: Error opening DB connection:", err)
	}
	err = DbConnection.Ping()
	if err != nil {
		log.Fatal("Error pinging DB:", err)
	}
}

// The database can be mysql, sqlite, postgres, etc.
func ensureTableExistsAndUpdatedForMysql(dbe DBEntityInterface) error {
	// Check if table exists
	tableName := DbSchema + "_" + dbe.GetTableName()
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
		createTableSQL := dbe.GetCreateTableSQL(DbSchema)
		log.Printf("Creating table with SQL: %s", createTableSQL)
		_, err := DbConnection.Exec(createTableSQL)
		if err != nil {
			return err
		}
		log.Printf("Created table %s", tableName)
	} else {
		// Table exists, check for schema updates
		// For simplicity, we will not implement schema migration logic here
		log.Printf("Table %s already exists", tableName)
	}

	return nil
}
func ensureTableExistsAndUpdatedForSqlite(dbe DBEntityInterface) error {
	// TODO: Implement for Sqlite
	return nil
}

func InitDBData() {

	// Check if the anonymous user exists, if not create it and associate it to the group "Guests"
	// Setup
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   "rprj",
	}

	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false
	// Check for "Guests" group
	guestGroup := repo.GetInstanceByTableName("groups")
	guestGroup.SetValue("name", "Guests")
	results, err := repo.Search(guestGroup, false, false, "")
	if err != nil {
		log.Printf("Failed to find or create 'Guests' group: %v\n", err)
		return
	}
	var guestGroupID string
	if len(results) == 1 {
		guestGroupID = results[0].GetValue("id").(string)
		log.Printf("Found existing 'Guests' group with ID %s\n", guestGroupID)
	} else {
		// Create the group
		newGroup := repo.GetInstanceByTableName("groups")
		newGroup.SetValue("name", "Guests")
		newGroup.SetValue("description", "Default group for anonymous users")
		created, err := repo.Insert(newGroup)
		if err != nil {
			log.Printf("Failed to create 'Guests' group: %v\n", err)
			return
		}
		guestGroupID = created.GetValue("id").(string)
		log.Printf("Created 'Guests' group with ID %s\n", guestGroupID)
	}

	// Check for anonymous user
	anonUser := repo.GetInstanceByTableName("users")
	anonUser.SetValue("login", "anonymous")
	results, err = repo.Search(anonUser, false, false, "")
	if err != nil {
		log.Printf("Failed to find or create 'anonymous' user: %v\n", err)
		return
	}
	if len(results) == 1 {
		log.Printf("Found existing 'anonymous' user with ID %s\n", results[0].GetValue("id").(string))
	} else {
		// Create the user
		newUser := repo.GetInstanceByTableName("users")
		newUser.SetValue("id", "-7")
		newUser.SetValue("login", "anonymous")
		newUser.SetValue("pwd", "") // No password for anonymous user
		newUser.SetValue("fullname", "Anonymous User")
		newUser.SetMetadata("group_ids", []string{guestGroupID})
		created, err := repo.Insert(newUser)
		if err != nil {
			log.Printf("Failed to create 'anonymous' user: %v\n", err)
			return
		}
		log.Printf("Created 'anonymous' user with ID %s\n", created.GetValue("id").(string))
	}
}

// Iterate over all registered DBEntity types and create tables if they do not exist or update their schema
func EnsureDBSchema() {
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

	var visit func(DBEntityInterface) bool
	visit = func(dbe DBEntityInterface) bool {
		className := dbe.GetTypeName()
		if temp[className] {
			// Circular dependency detected
			log.Printf("Warning: Circular dependency detected for %s", className)
			return false
		}
		if visited[className] {
			return true
		}

		temp[className] = true
		// Get foreign key dependencies
		foreignKeys := dbe.GetForeignKeys()
		for _, fk := range foreignKeys {
			// Find the referenced table's DBEntity
			for _, depDbe := range classInstances {
				if depDbe.GetTableName() == fk.RefTable {
					if !visit(depDbe) {
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
		if !visited[dbe.GetTypeName()] {
			visit(dbe)
		}
	}
	// reverse sorted to get correct order
	// for i, j := 0, len(sorted)-1; i < j; i, j = i+1, j-1 {
	// 	sorted[i], sorted[j] = sorted[j], sorted[i]
	// }
	// slices.Reverse(sorted)

	// Print sorted class names for debugging
	log.Print("DB Entities creation order:")
	for _, dbe := range sorted {
		log.Printf(" - %s\n", dbe.GetTableName())
		for _, fk := range dbe.GetForeignKeys() {
			log.Printf("    FK: %s -> %s(%s)\n", fk.Column, fk.RefTable, fk.RefColumn)
		}
	}

	classInstances = sorted

	for _, dbe := range classInstances {
		var err error
		className := dbe.GetTypeName()
		switch dbEngine {
		case "mysql":
			err = ensureTableExistsAndUpdatedForMysql(dbe)
		case "sqlite":
			err = ensureTableExistsAndUpdatedForSqlite(dbe)
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
