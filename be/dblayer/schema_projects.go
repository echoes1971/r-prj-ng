package dblayer

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

/*
CREATE TABLE `rprj_projects` (

	`id` varchar(16) NOT NULL,
	`owner` varchar(16) NOT NULL,
	`group_id` varchar(16) NOT NULL,
	`permissions` char(9) NOT NULL DEFAULT 'rwx------',
	`creator` varchar(16) NOT NULL,
	`creation_date` datetime DEFAULT NULL,
	`last_modify` varchar(16) NOT NULL,
	`last_modify_date` datetime DEFAULT NULL,
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime default null,
	`father_id` varchar(16) DEFAULT NULL,
	`name` varchar(255) NOT NULL,
	`description` text DEFAULT NULL,
	PRIMARY KEY (`id`),
	KEY `rprj_projects_0` (`id`),
	KEY `rprj_projects_1` (`owner`),
	KEY `rprj_projects_2` (`group_id`),
	KEY `rprj_projects_3` (`creator`),
	KEY `rprj_projects_4` (`last_modify`),
	KEY `rprj_projects_5` (`deleted_by`),
	KEY `rprj_projects_6` (`father_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

	class DBProject extends DBEObject {
		var $_typeName="DBProject";
		function getTableName() {
			return "projects";
		}
		function __construct($tablename=null, $names=null, $values=null, $attrs=null, $keys=null) {
			parent::__construct($tablename, $names, $values, $attrs, $keys, DBEObject::$_mycolumns);

		// 		$this->_columns['p_iva']=array('varchar(16)');
		}
		function getFK() {
			if($this->_fk==null) {
				$this->_fk=array();
			}
			if(count($this->_fk)==0) {
				parent::getFK();
		// 			$this->_fk[] = new ForeignKey('fk_obj_id','companies','id');
			}
			return $this->_fk;
		}
		//
		// @TODO da finire?
		//
		function _before_delete(&$dbmgr) {
			parent::_before_delete($dbmgr);

			// Cancello i legami con le compagnie
			// Cancello i legami con le persone
			// Cancello i legami con i progetti
		}
	}
*/
type DBProject struct {
	DBObject
}

func NewDBProject() *DBProject {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "char(9)", Constraints: []string{"NOT NULL", "DEFAULT 'rwx------'"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{"default null"}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{"DEFAULT NULL"}},
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
	return &DBProject{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBProject",
				"projects",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbProject *DBProject) NewInstance() DBEntityInterface {
	return NewDBProject()
}

/**
 * Defines the roles assumed by people in individual projects
 */
type DBProjectCompanyRole struct {
	DBObject
}

func NewDBProjectCompanyRole() *DBProjectCompanyRole {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "char(9)", Constraints: []string{"NOT NULL", "DEFAULT 'rwx------'"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{"default null"}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{"DEFAULT NULL"}},
		{Name: "order_position", Type: "int(11)", Constraints: []string{"DEFAULT 0"}},
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
	return &DBProjectCompanyRole{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBProjectCompanyRole",
				"projects_companies_roles",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbProjectCompanyRole *DBProjectCompanyRole) NewInstance() DBEntityInterface {
	return NewDBProjectCompanyRole()
}
func (dbProjectCompanyRole *DBProjectCompanyRole) GetOrderBy() []string {
	return []string{"order_position", "id"}
}

// TODO: test this. Plus, does it make sense to just get the max over the whole table?
func (dbProjectCompanyRole *DBProjectCompanyRole) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	if dbr.Verbose {
		log.Print("DBProjectCompanyRole.beforeInsert called")
	}
	err := dbProjectCompanyRole.DBObject.beforeInsert(dbr, tx)
	if err != nil {
		log.Print("DBProjectCompanyRole.beforeInsert: error in parent beforeInsert:", err)
		return err
	}

	query := "SELECT COALESCE(MAX(order_position), 0) AS order_position FROM " + DbSchema + "_" + dbProjectCompanyRole.GetTableName()
	results := dbr.Select("DBEObject", dbProjectCompanyRole.GetTableName(), query)
	if len(results) == 1 {
		intVal := 0
		fmt.Scanf(results[0].GetValue("order_position").(string), "%d", &intVal)
		dbProjectCompanyRole.SetValue("order_position", intVal+1)
	} else {
		dbProjectCompanyRole.SetValue("order_position", 1)
	}
	return nil
}

/*
CREATE TABLE `rprj_projects_companies` (
  `project_id` varchar(16) NOT NULL DEFAULT '',
  `company_id` varchar(16) NOT NULL DEFAULT '',
  `projects_companies_role_id` varchar(16) NOT NULL DEFAULT '',
  PRIMARY KEY (`project_id`,`company_id`,`projects_companies_role_id`),
  KEY `rprj_projects_companies_0` (`project_id`),
  KEY `rprj_projects_companies_1` (`company_id`),
  KEY `rprj_projects_companies_2` (`projects_companies_role_id`),
  KEY `rprj_projects_companies_3` (`project_id`),
  KEY `rprj_projects_companies_4` (`company_id`),
  KEY `rprj_projects_companies_5` (`projects_companies_role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/

/**
 * Associates companies to projects with a certain role
 */
type DBProjectCompany struct {
	DBAssociation
}

func NewDBProjectCompany() *DBProjectCompany {
	columns := []Column{
		{Name: "project_id", Type: "varchar(16)", Constraints: []string{"NOT NULL", "DEFAULT ''"}},
		{Name: "company_id", Type: "varchar(16)", Constraints: []string{"NOT NULL", "DEFAULT ''"}},
		{Name: "projects_companies_role_id", Type: "varchar(16)", Constraints: []string{"NOT NULL", "DEFAULT ''"}},
	}
	keys := []string{"project_id", "company_id", "projects_companies_role_id"}
	foreignKeys := []ForeignKey{
		{Column: "project_id", RefTable: "projects", RefColumn: "id"},
		{Column: "company_id", RefTable: "companies", RefColumn: "id"},
		{Column: "projects_companies_role_id", RefTable: "projects_companies_roles", RefColumn: "id"},
	}
	return &DBProjectCompany{
		DBAssociation: DBAssociation{
			DBEntity: *NewDBEntity(
				"DBProjectCompany",
				"projects_companies",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbProjectCompany *DBProjectCompany) NewInstance() DBEntityInterface {
	return NewDBProjectCompany()
}

/*
CREATE TABLE `rprj_projects_people_roles` (
  `id` varchar(16) NOT NULL,
  `owner` varchar(16) NOT NULL,
  `group_id` varchar(16) NOT NULL,
  `permissions` char(9) NOT NULL DEFAULT 'rwx------',
  `creator` varchar(16) NOT NULL,
  `creation_date` datetime DEFAULT NULL,
  `last_modify` varchar(16) NOT NULL,
  `last_modify_date` datetime DEFAULT NULL,
  `deleted_by` varchar(16) DEFAULT NULL,
  `deleted_date` datetime default null,
  `father_id` varchar(16) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `order_position` int(11) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `rprj_projects_people_roles_0` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/
/**
 * Defines the roles assumed by people in individual projects
 */
type DBProjectPeopleRole struct {
	DBAssociation
}

func NewDBProjectPeopleRole() *DBProjectPeopleRole {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "char(9)", Constraints: []string{"NOT NULL", "DEFAULT 'rwx------'"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{"default null"}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{"DEFAULT NULL"}},
		{Name: "order_position", Type: "int(11)", Constraints: []string{"DEFAULT 0"}},
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
	return &DBProjectPeopleRole{
		DBAssociation: DBAssociation{
			DBEntity: *NewDBEntity(
				"DBProjectPeopleRole",
				"projects_people_roles",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbProjectPeopleRole *DBProjectPeopleRole) NewInstance() DBEntityInterface {
	return NewDBProjectPeopleRole()
}
func (dbProjectPeopleRole *DBProjectPeopleRole) GetOrderBy() []string {
	return []string{"order_position", "id"}
}
func (dbProjectPeopleRole *DBProjectPeopleRole) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	if dbr.Verbose {
		log.Print("DBProjectPeopleRole.beforeInsert called")
	}
	err := dbProjectPeopleRole.DBAssociation.beforeInsert(dbr, tx)
	if err != nil {
		log.Print("DBProjectPeopleRole.beforeInsert: error in parent beforeInsert:", err)
		return err
	}

	query := "SELECT COALESCE(MAX(order_position), 0) AS order_position FROM " + DbSchema + "_" + dbProjectPeopleRole.GetTableName()
	results := dbr.Select("DBEObject", dbProjectPeopleRole.GetTableName(), query)
	if len(results) == 1 {
		intVal := 0
		fmt.Scanf(results[0].GetValue("order_position").(string), "%d", &intVal)
		dbProjectPeopleRole.SetValue("order_position", intVal+1)
	} else {
		dbProjectPeopleRole.SetValue("order_position", 1)
	}
	return nil
}

/*
CREATE TABLE `rprj_projects_people` (
  `project_id` varchar(16) NOT NULL DEFAULT '',
  `people_id` varchar(16) NOT NULL DEFAULT '',
  `projects_people_role_id` varchar(16) NOT NULL DEFAULT '',
  PRIMARY KEY (`project_id`,`people_id`,`projects_people_role_id`),
  KEY `rprj_projects_people_0` (`project_id`),
  KEY `rprj_projects_people_1` (`people_id`),
  KEY `rprj_projects_people_2` (`projects_people_role_id`),
  KEY `rprj_projects_people_3` (`project_id`),
  KEY `rprj_projects_people_4` (`people_id`),
  KEY `rprj_projects_people_5` (`projects_people_role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/
/**
 * Associates people to projects with a certain role
 */
type DBProjectPeople struct {
	DBAssociation
}

func NewDBProjectPeople() *DBProjectPeople {
	columns := []Column{
		{Name: "project_id", Type: "varchar(16)", Constraints: []string{"NOT NULL", "DEFAULT ''"}},
		{Name: "people_id", Type: "varchar(16)", Constraints: []string{"NOT NULL", "DEFAULT ''"}},
		{Name: "projects_people_role_id", Type: "varchar(16)", Constraints: []string{"NOT NULL", "DEFAULT ''"}},
	}
	keys := []string{"project_id", "people_id", "projects_people_role_id"}
	foreignKeys := []ForeignKey{
		{Column: "project_id", RefTable: "projects", RefColumn: "id"},
		{Column: "people_id", RefTable: "people", RefColumn: "id"},
		{Column: "projects_people_role_id", RefTable: "projects_people_roles", RefColumn: "id"},
	}
	return &DBProjectPeople{
		DBAssociation: DBAssociation{
			DBEntity: *NewDBEntity(
				"DBProjectPeople",
				"projects_people",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbProjectPeople *DBProjectPeople) NewInstance() DBEntityInterface {
	return NewDBProjectPeople()
}

/*
CREATE TABLE `rprj_projects_projects_roles` (
  `id` varchar(16) NOT NULL,
  `owner` varchar(16) NOT NULL,
  `group_id` varchar(16) NOT NULL,
  `permissions` char(9) NOT NULL DEFAULT 'rwx------',
  `creator` varchar(16) NOT NULL,
  `creation_date` datetime DEFAULT NULL,
  `last_modify` varchar(16) NOT NULL,
  `last_modify_date` datetime DEFAULT NULL,
  `deleted_by` varchar(16) DEFAULT NULL,
  `deleted_date` datetime default null,
  `father_id` varchar(16) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `order_position` int(11) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `rprj_projects_projects_roles_0` (`id`),
  KEY `rprj_projects_projects_roles_1` (`owner`),
  KEY `rprj_projects_projects_roles_2` (`group_id`),
  KEY `rprj_projects_projects_roles_3` (`creator`),
  KEY `rprj_projects_projects_roles_4` (`last_modify`),
  KEY `rprj_projects_projects_roles_5` (`deleted_by`),
  KEY `rprj_projects_projects_roles_6` (`father_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/
/**
 * Defines the roles assumed by projects in individual projects
 */
type DBProjectProjectsRole struct {
	DBAssociation
}

func NewDBProjectProjectsRole() *DBProjectProjectsRole {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "char(9)", Constraints: []string{"NOT NULL", "DEFAULT 'rwx------'"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{"default null"}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{"DEFAULT NULL"}},
		{Name: "order_position", Type: "int(11)", Constraints: []string{"DEFAULT 0"}},
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
	return &DBProjectProjectsRole{
		DBAssociation: DBAssociation{
			DBEntity: *NewDBEntity(
				"DBProjectProjectsRole",
				"projects_projects_roles",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbProjectProjectsRole *DBProjectProjectsRole) NewInstance() DBEntityInterface {
	return NewDBProjectProjectsRole()
}
func (dbProjectProjectsRole *DBProjectProjectsRole) GetOrderBy() []string {
	return []string{"order_position", "id"}
}
func (dbProjectProjectsRole *DBProjectProjectsRole) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	if dbr.Verbose {
		log.Print("DBProjectProjectsRole.beforeInsert called")
	}
	err := dbProjectProjectsRole.DBAssociation.beforeInsert(dbr, tx)
	if err != nil {
		log.Print("DBProjectProjectsRole.beforeInsert: error in parent beforeInsert:", err)
		return err
	}

	query := "SELECT COALESCE(MAX(order_position), 0) AS order_position FROM " + DbSchema + "_" + dbProjectProjectsRole.GetTableName()
	results := dbr.Select("DBEObject", dbProjectProjectsRole.GetTableName(), query)
	if len(results) == 1 {
		intVal := 0
		fmt.Scanf(results[0].GetValue("order_position").(string), "%d", &intVal)
		dbProjectProjectsRole.SetValue("order_position", intVal+1)
	} else {
		dbProjectProjectsRole.SetValue("order_position", 1)
	}
	return nil
}

/*
CREATE TABLE `rprj_projects_projects` (
  `project_id` varchar(16) NOT NULL DEFAULT '',
  `project2_id` varchar(16) NOT NULL DEFAULT '',
  `projects_projects_role_id` varchar(16) NOT NULL DEFAULT '',
  PRIMARY KEY (`project_id`,`project2_id`,`projects_projects_role_id`),
  KEY `rprj_projects_projects_0` (`project_id`),
  KEY `rprj_projects_projects_1` (`project2_id`),
  KEY `rprj_projects_projects_2` (`projects_projects_role_id`),
  KEY `rprj_projects_projects_3` (`project_id`),
  KEY `rprj_projects_projects_4` (`project2_id`),
  KEY `rprj_projects_projects_5` (`projects_projects_role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/
/**
 * Associates projects to projects with a certain role
 */
type DBProjectProjects struct {
	DBAssociation
}

func NewDBProjectProjects() *DBProjectProjects {
	columns := []Column{
		{Name: "project_id", Type: "varchar(16)", Constraints: []string{"NOT NULL", "DEFAULT ''"}},
		{Name: "project2_id", Type: "varchar(16)", Constraints: []string{"NOT NULL", "DEFAULT ''"}},
		{Name: "projects_projects_role_id", Type: "varchar(16)", Constraints: []string{"NOT NULL", "DEFAULT ''"}},
	}
	keys := []string{"project_id", "project2_id", "projects_projects_role_id"}
	foreignKeys := []ForeignKey{
		{Column: "project_id", RefTable: "projects", RefColumn: "id"},
		{Column: "project2_id", RefTable: "projects", RefColumn: "id"},
		{Column: "projects_projects_role_id", RefTable: "projects_projects_roles", RefColumn: "id"},
	}
	return &DBProjectProjects{
		DBAssociation: DBAssociation{
			DBEntity: *NewDBEntity(
				"DBProjectProjects",
				"projects_projects",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbProjectProjects *DBProjectProjects) NewInstance() DBEntityInterface {
	return NewDBProjectProjects()
}

/*
CREATE TABLE `rprj_timetracks` (
  `id` varchar(16) NOT NULL,
  `owner` varchar(16) NOT NULL,
  `group_id` varchar(16) NOT NULL,
  `permissions` char(9) NOT NULL DEFAULT 'rwx------',
  `creator` varchar(16) NOT NULL,
  `creation_date` datetime DEFAULT NULL,
  `last_modify` varchar(16) NOT NULL,
  `last_modify_date` datetime DEFAULT NULL,
  `deleted_by` varchar(16) DEFAULT NULL,
  `deleted_date` datetime default null,
  `father_id` varchar(16) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `fk_obj_id` varchar(16) DEFAULT NULL,
  `fk_progetto` varchar(16) DEFAULT NULL,
  `dalle_ore` datetime DEFAULT NULL,
  `alle_ore` datetime DEFAULT NULL,
  `ore_intervento` datetime DEFAULT NULL,
  `ore_viaggio` datetime DEFAULT NULL,
  `km_viaggio` int(11) NOT NULL DEFAULT 0,
  `luogo_di_intervento` int(11) NOT NULL DEFAULT 0,
  `stato` int(11) NOT NULL DEFAULT 0,
  `costo_per_ora` float NOT NULL DEFAULT 0,
  `costo_valuta` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `rprj_timetracks_0` (`id`),
  KEY `rprj_timetracks_1` (`owner`),
  KEY `rprj_timetracks_2` (`group_id`),
  KEY `rprj_timetracks_3` (`creator`),
  KEY `rprj_timetracks_4` (`last_modify`),
  KEY `rprj_timetracks_5` (`deleted_by`),
  KEY `rprj_timetracks_6` (`father_id`),
  KEY `rprj_timetracks_7` (`fk_obj_id`),
  KEY `rprj_timetracks_8` (`fk_obj_id`),
  KEY `rprj_timetracks_9` (`fk_obj_id`),
  KEY `rprj_timetracks_10` (`fk_progetto`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/

/*
*

  - Time tracking entries associated with projects

    TODO:
    1. translate fields in English
    2. sql steps to translate data in existing DBs to new fields
*/
type DBTimeTrack struct {
	DBObject
}

func NewDBTimeTrack() *DBTimeTrack {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "char(9)", Constraints: []string{"NOT NULL", "DEFAULT 'rwx------'"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{"default null"}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{"DEFAULT NULL"}},
		{Name: "fk_obj_id", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "fk_progetto", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "dalle_ore", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "alle_ore", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "ore_intervento", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "ore_viaggio", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "km_viaggio", Type: "int(11)", Constraints: []string{"NOT NULL", "DEFAULT 0"}},
		{Name: "luogo_di_intervento", Type: "int(11)", Constraints: []string{"NOT NULL", "DEFAULT 0"}},
		{Name: "stato", Type: "int(11)", Constraints: []string{"NOT NULL", "DEFAULT 0"}},
		{Name: "costo_per_ora", Type: "float", Constraints: []string{"NOT NULL", "DEFAULT 0"}},
		{Name: "costo_valuta", Type: "varchar(255)", Constraints: []string{"DEFAULT NULL"}},
	}
	keys := []string{"id"}
	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "objects", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "projects", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "folders", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "todo", RefColumn: "id"},
		{Column: "fk_progetto", RefTable: "projects", RefColumn: "id"},
	}
	return &DBTimeTrack{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBTimeTrack",
				"timetracks",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbTimeTrack *DBTimeTrack) NewInstance() DBEntityInterface {
	return NewDBTimeTrack()
}
func (dbTimeTrack *DBTimeTrack) GetOrderBy() []string {
	return []string{"fk_progetto", "stato desc", "dalle_ore desc", "fk_obj_id", "name"}
}

/*
CREATE TABLE `rprj_todo` (

	`id` varchar(16) NOT NULL,
	`owner` varchar(16) NOT NULL,
	`group_id` varchar(16) NOT NULL,
	`permissions` char(9) NOT NULL DEFAULT 'rwx------',
	`creator` varchar(16) NOT NULL,
	`creation_date` datetime DEFAULT NULL,
	`last_modify` varchar(16) NOT NULL,
	`last_modify_date` datetime DEFAULT NULL,
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime default null,
	`father_id` varchar(16) DEFAULT NULL,
	`name` varchar(255) NOT NULL,
	`description` text DEFAULT NULL,
	`priority` int(11) NOT NULL DEFAULT 0,
	`data_segnalazione` datetime DEFAULT NULL,
	`fk_segnalato_da` varchar(16) DEFAULT NULL,
	`fk_cliente` varchar(16) DEFAULT NULL,
	`fk_progetto` varchar(16) DEFAULT NULL,
	`fk_funzionalita` varchar(16) DEFAULT NULL,
	`fk_tipo` varchar(16) DEFAULT NULL,
	`stato` int(11) NOT NULL DEFAULT 0,
	`descrizione` text NOT NULL,
	`intervento` text NOT NULL,
	`data_chiusura` datetime DEFAULT NULL,
	PRIMARY KEY (`id`),
	KEY `rprj_todo_0` (`id`),
	KEY `rprj_todo_1` (`owner`),
	KEY `rprj_todo_2` (`group_id`),
	KEY `rprj_todo_3` (`creator`),
	KEY `rprj_todo_4` (`last_modify`),
	KEY `rprj_todo_5` (`deleted_by`),
	KEY `rprj_todo_6` (`father_id`),
	KEY `rprj_todo_7` (`fk_segnalato_da`),
	KEY `rprj_todo_8` (`fk_cliente`),
	KEY `rprj_todo_9` (`fk_progetto`),
	KEY `rprj_todo_10` (`father_id`),
	KEY `rprj_todo_11` (`father_id`),
	KEY `rprj_todo_12` (`fk_tipo`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/
type DBTodo struct {
	DBObject
}

/*
*
TODO:
 1. translate fields in English
 2. sql steps to translate data in existing DBs to new fields
*/
func NewDBTodo() *DBTodo {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "char(9)", Constraints: []string{"NOT NULL", "DEFAULT 'rwx------'"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{"default null"}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{"DEFAULT NULL"}},
		{Name: "priority", Type: "int(11)", Constraints: []string{"NOT NULL", "DEFAULT 0"}},
		{Name: "data_segnalazione", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "fk_segnalato_da", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "fk_cliente", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "fk_progetto", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "fk_funzionalita", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "fk_tipo", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "stato", Type: "int(11)", Constraints: []string{"NOT NULL", "DEFAULT 0"}},
		{Name: "descrizione", Type: "text", Constraints: []string{"NOT NULL"}},
		{Name: "intervento", Type: "text", Constraints: []string{"NOT NULL"}},
		{Name: "data_chiusura", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
	}
	keys := []string{"id"}
	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "objects", RefColumn: "id"},

		{Column: "fk_segnalato_da", RefTable: "people", RefColumn: "id"},
		{Column: "fk_cliente", RefTable: "companies", RefColumn: "id"},
		{Column: "fk_progetto", RefTable: "projects", RefColumn: "id"},

		{Column: "father_id", RefTable: "folders", RefColumn: "id"},
		{Column: "father_id", RefTable: "todo", RefColumn: "id"},

		{Column: "fk_tipo", RefTable: "todo_tipo", RefColumn: "id"},
	}
	return &DBTodo{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBTodo",
				"todo",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbTodo *DBTodo) NewInstance() DBEntityInterface {
	return NewDBTodo()
}
func (dbTodo *DBTodo) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	if dbr.Verbose {
		log.Print("DBTodo.beforeInsert called")
	}
	err := dbTodo.DBObject.beforeInsert(dbr, tx)
	if err != nil {
		log.Print("DBTodo.beforeInsert: error in parent beforeInsert:", err)
		return err
	}

	if !dbTodo.HasValue("data_segnalazione") {
		dbTodo.SetValue("data_segnalazione", time.Now().Format("2006-01-02 15:04:05"))
	}
	// $data_segnalazione = $this->getValue('data_segnalazione');
	// if($data_segnalazione==null || $data_segnalazione=="" || $data_segnalazione=="00:00" || $data_segnalazione=="0000/00/00 00:00") {
	//   $this->setValue('data_segnalazione', $this->_getTodayString());
	// }

	err = dbTodo.RULE_check_closure()
	if err != nil {
		log.Print("DBTodo.beforeInsert: error in _RULE_check_closure:", err)
		return err
	}

	return nil
}

func (dbTodo *DBTodo) beforeUpdate(dbr *DBRepository, tx *sql.Tx) error {
	if dbr.Verbose {
		log.Print("DBTodo.beforeUpdate called")
	}
	err := dbTodo.DBObject.beforeUpdate(dbr, tx)
	if err != nil {
		log.Print("DBTodo.beforeUpdate: error in parent beforeUpdate:", err)
		return err
	}

	err = dbTodo.RULE_check_closure()
	if err != nil {
		log.Print("DBTodo.beforeUpdate: error in _RULE_check_closure:", err)
		return err
	}
	err = dbTodo.RULE_check_reopening()
	if err != nil {
		log.Print("DBTodo.beforeUpdate: error in _RULE_check_reopening:", err)
		return err
	}

	return nil
}

/**
 * IF data_chiusura without value AND stato=100% ==> data_chiusura=today
 */
func (dbTodo *DBTodo) RULE_check_closure() error {
	if !dbTodo.HasValue("data_chiusura") {
		statoVal, ok := dbTodo.GetValue("stato").(int)
		if ok && statoVal >= 100 {
			dbTodo.SetValue("data_chiusura", time.Now().Format("2006-01-02 15:04:05"))
			dbTodo.SetValue("stato", 100)
		}
	}
	// $data_chiusura = $this->getValue('data_chiusura');
	// if($data_chiusura==null || $data_chiusura=="" || $data_chiusura=="0000/00/00 00:00") {
	// 	$stato = $this->getValue('stato');
	// 	if(is_numeric($stato) && intval($stato)>=100) {
	// 			$this->setValue('data_chiusura', $this->_getTodayString());
	// 			$this->setValue('stato', 100);
	// 	}
	// }

	return nil
}

/**
 * IF data_chiusura has value AND stato<100% => data_chiusura=null
 */
func (dbTodo *DBTodo) RULE_check_reopening() error {
	if !dbTodo.HasValue("data_chiusura") {
		return nil
	}
	statoVal, ok := dbTodo.GetValue("stato").(int)
	if ok && statoVal < 100 {
		dbTodo.SetValue("data_chiusura", nil)
	}
	return nil

	// $data_chiusura = $this->getValue('data_chiusura');
	// if($data_chiusura==null || $data_chiusura=="" || $data_chiusura=="0000/00/00 00:00") {
	// } else {
	// 	$stato = $this->getValue('stato');
	// 	if(is_numeric($stato) && intval($stato)<100) {
	// 			$this->setValue('data_chiusura', '');
	// 	}
	// }
}

/*
CREATE TABLE `rprj_todo_tipo` (

	`id` varchar(16) NOT NULL,
	`owner` varchar(16) NOT NULL,
	`group_id` varchar(16) NOT NULL,
	`permissions` char(9) NOT NULL DEFAULT 'rwx------',
	`creator` varchar(16) NOT NULL,
	`creation_date` datetime DEFAULT NULL,
	`last_modify` varchar(16) NOT NULL,
	`last_modify_date` datetime DEFAULT NULL,
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime default null,
	`father_id` varchar(16) DEFAULT NULL,
	`name` varchar(255) NOT NULL,
	`description` text DEFAULT NULL,
	`order_position` int(11) DEFAULT 0,
	PRIMARY KEY (`id`),
	KEY `rprj_todo_tipo_0` (`id`),
	KEY `rprj_todo_tipo_1` (`owner`),
	KEY `rprj_todo_tipo_2` (`group_id`),
	KEY `rprj_todo_tipo_3` (`creator`),
	KEY `rprj_todo_tipo_4` (`last_modify`),
	KEY `rprj_todo_tipo_5` (`deleted_by`),
	KEY `rprj_todo_tipo_6` (`father_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/
type DBTodoTipo struct {
	DBObject
}

func NewDBTodoTipo() *DBTodoTipo {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "char(9)", Constraints: []string{"NOT NULL", "DEFAULT 'rwx------'"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{"default null"}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{"DEFAULT NULL"}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{"DEFAULT NULL"}},
		{Name: "order_position", Type: "int(11)", Constraints: []string{"DEFAULT 0"}},
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
	return &DBTodoTipo{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBTodoTipo",
				"todo_tipo",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbTodoTipo *DBTodoTipo) NewInstance() DBEntityInterface {
	return NewDBTodoTipo()
}
func (dbTodoTipo *DBTodoTipo) GetOrderBy() []string {
	return []string{"order_position", "id"}
}
func (dbTodoTipo *DBTodoTipo) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	if dbr.Verbose {
		log.Print("DBTodoTipo.beforeInsert called")
	}
	err := dbTodoTipo.DBObject.beforeInsert(dbr, tx)
	if err != nil {
		log.Print("DBTodoTipo.beforeInsert: error in parent beforeInsert:", err)
		return err
	}

	query := "SELECT COALESCE(MAX(order_position), 0) AS order_position FROM " + DbSchema + "_" + dbTodoTipo.GetTableName()
	results := dbr.Select("DBEObject", dbTodoTipo.GetTableName(), query)
	if len(results) == 1 {
		intVal := 0
		fmt.Scanf(results[0].GetValue("order_position").(string), "%d", &intVal)
		dbTodoTipo.SetValue("order_position", intVal+1)
	} else {
		dbTodoTipo.SetValue("order_position", 1)
	}
	return nil
}
