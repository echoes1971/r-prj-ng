package dblayer

/*
CREATE TABLE IF NOT EXISTS `rra_countrylist` (

	`id` varchar(16) NOT NULL DEFAULT '',
	`Common_Name` varchar(255) DEFAULT NULL,
	`Formal_Name` varchar(255) DEFAULT NULL,
	`Type` varchar(255) DEFAULT NULL,
	`Sub_Type` varchar(255) DEFAULT NULL,
	`Sovereignty` varchar(255) DEFAULT NULL,
	`Capital` varchar(255) DEFAULT NULL,
	`ISO_4217_Currency_Code` varchar(255) DEFAULT NULL,
	`ISO_4217_Currency_Name` varchar(255) DEFAULT NULL,
	`ITU_T_Telephone_Code` varchar(255) DEFAULT NULL,
	`ISO_3166_1_2_Letter_Code` varchar(255) DEFAULT NULL,
	`ISO_3166_1_3_Letter_Code` varchar(255) DEFAULT NULL,
	`ISO_3166_1_Number` varchar(255) DEFAULT NULL,
	`IANA_Country_Code_TLD` varchar(255) DEFAULT NULL,
	PRIMARY KEY (`id`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/

type DBCountry struct {
	DBEntity
}

func NewDBCountry() DBEntityInterface {
	// Define columns
	columns := []Column{
		{"id", "VARCHAR(16)", []string{"NOT NULL"}},
		{"Common_Name", "VARCHAR(255)", []string{}},
		{"Formal_Name", "VARCHAR(255)", []string{}},
		{"Type", "VARCHAR(255)", []string{}},
		{"Sub_Type", "VARCHAR(255)", []string{}},
		{"Sovereignty", "VARCHAR(255)", []string{}},
		{"Capital", "VARCHAR(255)", []string{}},
		{"ISO_4217_Currency_Code", "VARCHAR(255)", []string{}},
		{"ISO_4217_Currency_Name", "VARCHAR(255)", []string{}},
		{"ITU_T_Telephone_Code", "VARCHAR(255)", []string{}},
		{"ISO_3166_1_2_Letter_Code", "VARCHAR(255)", []string{}},
		{"ISO_3166_1_3_Letter_Code", "VARCHAR(255)", []string{}},
		{"ISO_3166_1_Number", "VARCHAR(255)", []string{}},
		{"IANA_Country_Code_TLD", "VARCHAR(255)", []string{}},
	}
	// Define keys
	keys := []string{
		"id",
	}
	return &DBCountry{
		DBEntity: *NewDBEntity(
			"DBCountry",
			"countrylist",
			columns,
			keys,
			[]ForeignKey{},
			make(map[string]any),
		),
	}
}
func (dbCountry *DBCountry) NewInstance() DBEntityInterface {
	columns := make([]Column, 0, len(dbCountry.columns))
	for _, col := range dbCountry.columns {
		columns = append(columns, col)
	}
	return &DBCountry{
		DBEntity: DBEntity{
			typename:   dbCountry.typename,
			tablename:  dbCountry.tablename,
			columns:    dbCountry.columns,
			keys:       dbCountry.keys,
			dictionary: make(map[string]any),
		},
	}
}

/*
CREATE TABLE IF NOT EXISTS `rra_companies` (

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
	`street` varchar(255) DEFAULT NULL,
	`zip` varchar(255) DEFAULT NULL,
	`city` varchar(255) DEFAULT NULL,
	`state` varchar(255) DEFAULT NULL,
	`fk_countrylist_id` varchar(16) DEFAULT NULL,
	`phone` varchar(255) DEFAULT NULL,
	`fax` varchar(255) DEFAULT NULL,
	`email` varchar(255) DEFAULT NULL,
	`url` varchar(255) DEFAULT NULL,
	`p_iva` varchar(16) DEFAULT NULL,
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	PRIMARY KEY (`id`),
	KEY `rra_companies_idx1` (`id`),
	KEY `rra_companies_idx2` (`owner`),
	KEY `rra_companies_idx3` (`name`),
	KEY `rra_companies_idx4` (`creator`),
	KEY `rra_companies_idx5` (`last_modify`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBCompany struct {
	DBObject
}

func NewDBCompany() *DBCompany {
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
		{Name: "street", Type: "varchar(255)", Constraints: []string{}},
		{Name: "zip", Type: "varchar(255)", Constraints: []string{}},
		{Name: "city", Type: "varchar(255)", Constraints: []string{}},
		{Name: "state", Type: "varchar(255)", Constraints: []string{}},
		{Name: "fk_countrylist_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "phone", Type: "varchar(255)", Constraints: []string{}},
		{Name: "fax", Type: "varchar(255)", Constraints: []string{}},
		{Name: "email", Type: "varchar(255)", Constraints: []string{}},
		{Name: "url", Type: "varchar(255)", Constraints: []string{}},
		{Name: "p_iva", Type: "varchar(16)", Constraints: []string{}},
	}
	keys := []string{"id"}
	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "objects", RefColumn: "id"},
		{Column: "fk_countrylist_id", RefTable: "countrylist", RefColumn: "id"},
	}
	return &DBCompany{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBCompany",
				"companies",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbCompany *DBCompany) NewInstance() DBEntityInterface {
	return NewDBCompany()
}

/*
CREATE TABLE IF NOT EXISTS `rra_people` (

	`id` varchar(16) NOT NULL DEFAULT '',
	`owner` varchar(16) NOT NULL DEFAULT '',
	`group_id` varchar(16) NOT NULL DEFAULT '',
	`permissions` varchar(9) NOT NULL DEFAULT 'rwx------',
	`creator` varchar(16) NOT NULL DEFAULT '',
	`creation_date` datetime DEFAULT NULL,
	`last_modify` varchar(16) NOT NULL DEFAULT '',
	`last_modify_date` datetime DEFAULT NULL,
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	`father_id` varchar(16) DEFAULT NULL,
	`name` varchar(255) NOT NULL DEFAULT '',
	`description` text,
	`street` varchar(255) DEFAULT NULL,
	`zip` varchar(255) DEFAULT NULL,
	`city` varchar(255) DEFAULT NULL,
	`state` varchar(255) DEFAULT NULL,
	`fk_countrylist_id` varchar(16) DEFAULT NULL,
	`fk_companies_id` varchar(16) DEFAULT NULL,
	`fk_users_id` varchar(16) DEFAULT NULL,
	`phone` varchar(255) DEFAULT NULL,
	`mobile` varchar(255) DEFAULT NULL,
	`fax` varchar(255) DEFAULT NULL,
	`email` varchar(255) DEFAULT NULL,
	`url` varchar(255) DEFAULT NULL,
	`office_phone` varchar(255) DEFAULT NULL,
	`codice_fiscale` varchar(20) DEFAULT NULL,
	`p_iva` varchar(16) DEFAULT NULL,
	PRIMARY KEY (`id`),
	KEY `rra_people_idx1` (`id`),
	KEY `rra_people_idx2` (`owner`),
	KEY `rra_people_idx3` (`name`),
	KEY `rra_people_idx4` (`creator`),
	KEY `rra_people_idx5` (`last_modify`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBPerson struct {
	DBObject
}

func NewDBPerson() *DBPerson {
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
		{Name: "street", Type: "varchar(255)", Constraints: []string{}},
		{Name: "zip", Type: "varchar(255)", Constraints: []string{}},
		{Name: "city", Type: "varchar(255)", Constraints: []string{}},
		{Name: "state", Type: "varchar(255)", Constraints: []string{}},
		{Name: "fk_countrylist_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "fk_companies_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "fk_users_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "phone", Type: "varchar(255)", Constraints: []string{}},
		{Name: "mobile", Type: "varchar(255)", Constraints: []string{}},
		{Name: "fax", Type: "varchar(255)", Constraints: []string{}},
		{Name: "email", Type: "varchar(255)", Constraints: []string{}},
		{Name: "url", Type: "varchar(255)", Constraints: []string{}},
		{Name: "office_phone", Type: "varchar(255)", Constraints: []string{}},
		{Name: "codice_fiscale", Type: "varchar(20)", Constraints: []string{}},
		{Name: "p_iva", Type: "varchar(16)", Constraints: []string{}},
	}
	keys := []string{"id"}
	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "objects", RefColumn: "id"},
		{Column: "fk_countrylist_id", RefTable: "countrylist", RefColumn: "id"},
		{Column: "fk_companies_id", RefTable: "companies", RefColumn: "id"},
		{Column: "fk_users_id", RefTable: "users", RefColumn: "id"},
	}
	return &DBPerson{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBPerson",
				"people",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbPerson *DBPerson) NewInstance() DBEntityInterface {
	return NewDBPerson()
}
