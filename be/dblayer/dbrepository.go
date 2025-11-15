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

func (dbr *DBRepository) GetInstanceByClassName(classname string) *DBEntity {
	return dbr.factory.GetInstanceByClassName(classname)
}
func (dbr *DBRepository) GetInstanceByTableName(tablename string) *DBEntity {
	return dbr.factory.GetInstanceByTableName(tablename)
}

func (dbr *DBRepository) buildTableName(dbe *DBEntity) string {
	if dbr.DbContext != nil && dbr.DbContext.Schema != "" {
		return dbr.DbContext.Schema + "_" + dbe.GetTableName()
	}
	return dbe.GetTableName()
}

func (dbr *DBRepository) Search(dbe *DBEntity, useLike bool, caseSensitive bool, orderBy string) ([]*DBEntity, error) {
	if dbr.Verbose {
		log.Print("DBRepository::Search: dbe=", dbe)
	}

	// 1. Build WHERE clauses
	whereClauses := make([]string, 0)
	args := make([]interface{}, 0) // slice of interface{} for values

	for key, value := range dbe.dictionary {
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
		log.Print("DBRepository::Search: query=", query, " args=", args)
	}

	// 3. Execute the query
	rows, err := dbr.DbConnection.Query(query, args...)
	if err != nil {
		log.Print("DBRepository::Search: Query error:", err)
		return nil, err
	}
	defer rows.Close()

	// 4. Process results
	results := make([]*DBEntity, 0)
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
