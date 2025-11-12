package models

/** *********************************** RRA Framework: start. *********************************** */

/*
CREATE TABLE IF NOT EXISTS `rra_dbversion` (

	`version` int(11) NOT NULL DEFAULT '0',
	KEY `rra_dbversion_0` (`version`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBVersion struct {
	Version int
}

// Struttura che rappresenta la tabella
type DBUser struct {
	ID       string
	Login    string
	Pwd      string
	PwdSalt  string
	Fullname string
	GroupID  string
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
	ID          string
	Name        string
	Description string
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
type DBUserGroup struct {
	UserID  string
	GroupID string
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
type DBOAuthToken struct {
	TokenID      string
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpiresAt    string // DATETIME in formato stringa
	CreatedAt    string // DATETIME in formato stringa
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
	IP    string
	Data  string // date in formato stringa
	Ora   string // time in formato stringa
	Count int
	URL   string
	Note  string
	Note2 string
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
	ID             string
	Owner          string
	GroupID        string
	Permissions    string
	Creator        string
	CreationDate   string // DATETIME in formato stringa
	LastModify     string
	LastModifyDate string // DATETIME in formato stringa
	FatherID       string
	Name           string
	Description    string
	DeletedBy      string
	DeletedDate    string // DATETIME in formato stringa
}

/** *********************************** RRA Framework: end. *********************************** */

/** *********************************** RRA Contacts: start. *********************************** */

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
	ID                   string
	CommonName           string
	FormalName           string
	Type                 string
	SubType              string
	Sovereignty          string
	Capital              string
	ISO4217CurrencyCode  string
	ISO4217CurrencyName  string
	ITUTTelephoneCode    string
	ISO31661_2LetterCode string
	ISO31661_3LetterCode string
	ISO31661Number       string
	IANACountryCodeTLD   string
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
	ID              string
	Owner           string
	GroupID         string
	Permissions     string
	Creator         string
	CreationDate    string // DATETIME in formato stringa
	LastModify      string
	LastModifyDate  string // DATETIME in formato stringa
	FatherID        string
	Name            string
	Description     string
	Street          string
	ZIP             string
	City            string
	State           string
	FKCountryListID string
	Phone           string
	Fax             string
	Email           string
	URL             string
	PIVA            string
	DeletedBy       string
	DeletedDate     string // DATETIME in formato stringa
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
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	PRIMARY KEY (`id`),
	KEY `rra_people_idx1` (`id`),
	KEY `rra_people_idx2` (`owner`),
	KEY `rra_people_idx3` (`name`),
	KEY `rra_people_idx4` (`creator`),
	KEY `rra_people_idx5` (`last_modify`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBPerson struct {
	ID              string
	Owner           string
	GroupID         string
	Permissions     string
	Creator         string
	CreationDate    string // DATETIME in formato stringa
	LastModify      string
	LastModifyDate  string // DATETIME in formato stringa
	FatherID        string
	Name            string
	Description     string
	Street          string
	ZIP             string
	City            string
	State           string
	FKCountryListID string
	FKCompaniesID   string
	FKUsersID       string
	Phone           string
	Mobile          string
	Fax             string
	Email           string
	URL             string
	OfficePhone     string
	CodiceFiscale   string
	PIVA            string
	DeletedBy       string
	DeletedDate     string // DATETIME in formato stringa
}

/** *********************************** RRA Contacts: end. *********************************** */

/** *********************************** CMS: start. *********************************** */

/*
CREATE TABLE IF NOT EXISTS `rra_events` (

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
	`fk_obj_id` varchar(16) DEFAULT NULL,
	`start_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	`end_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	`all_day` char(1) NOT NULL DEFAULT '1',
	`url` varchar(255) DEFAULT NULL,
	`alarm` char(1) DEFAULT '0',
	`alarm_minute` int(11) DEFAULT '0',
	`alarm_unit` char(1) DEFAULT '0',
	`before_event` char(1) DEFAULT '0',
	`category` varchar(255) DEFAULT '',
	`recurrence` char(1) DEFAULT '0',
	`recurrence_type` char(1) DEFAULT '0',
	`daily_every_x` int(11) DEFAULT '0',
	`weekly_every_x` int(11) DEFAULT '0',
	`weekly_day_of_the_week` char(1) DEFAULT '0',
	`monthly_every_x` int(11) DEFAULT '0',
	`monthly_day_of_the_month` int(11) DEFAULT '0',
	`monthly_week_number` int(11) DEFAULT '0',
	`monthly_week_day` char(1) DEFAULT '0',
	`yearly_month_number` int(11) DEFAULT '0',
	`yearly_month_day` int(11) DEFAULT '0',
	`yearly_week_number` int(11) DEFAULT '0',
	`yearly_week_day` char(1) DEFAULT '0',
	`yearly_day_of_the_year` int(11) DEFAULT '0',
	`recurrence_times` int(11) DEFAULT '0',
	`recurrence_end_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	PRIMARY KEY (`id`),
	KEY `rra_events_0` (`id`),
	KEY `rra_events_1` (`owner`),
	KEY `rra_events_2` (`group_id`),
	KEY `rra_events_3` (`creator`),
	KEY `rra_events_4` (`last_modify`),
	KEY `rra_events_5` (`deleted_by`),
	KEY `rra_events_6` (`father_id`),
	KEY `rra_events_7` (`fk_obj_id`),
	KEY `rra_events_8` (`fk_obj_id`),
	KEY `rra_events_9` (`fk_obj_id`),
	KEY `rra_events_10` (`fk_obj_id`),
	KEY `rra_events_idx2` (`start_date`),
	KEY `rra_events_idx3` (`end_date`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBEvent struct {
	ID                   string
	Owner                string
	GroupID              string
	Permissions          string
	Creator              string
	CreationDate         string // DATETIME in formato stringa
	LastModify           string
	LastModifyDate       string // DATETIME in formato stringa
	DeletedBy            string
	DeletedDate          string // DATETIME in formato stringa
	FatherID             string
	Name                 string
	Description          string
	FKObjID              string
	StartDate            string // DATETIME in formato stringa
	EndDate              string // DATETIME in formato stringa
	AllDay               string
	URL                  string
	Alarm                string
	AlarmMinute          int
	AlarmUnit            string
	BeforeEvent          string
	Category             string
	Recurrence           string
	RecurrenceType       string
	DailyEveryX          int
	WeeklyEveryX         int
	WeeklyDayOfTheWeek   string
	MonthlyEveryX        int
	MonthlyDayOfTheMonth int
	MonthlyWeekNumber    int
	MonthlyWeekDay       string
	YearlyMonthNumber    int
	YearlyMonthDay       int
	YearlyWeekNumber     int
	YearlyWeekDay        string
	YearlyDayOfTheYear   int
	RecurrenceTimes      int
	RecurrenceEndDate    string // DATETIME in formato stringa
}

/*
CREATE TABLE IF NOT EXISTS `rra_files` (

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
	`fk_obj_id` varchar(16) DEFAULT NULL,
	`path` text,
	`filename` text NOT NULL,
	`checksum` varchar(40) DEFAULT NULL,
	`mime` varchar(255) DEFAULT NULL,
	`alt_link` varchar(255) NOT NULL DEFAULT '',
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	PRIMARY KEY (`id`),
	KEY `rra_files_idx1` (`id`),
	KEY `rra_files_idx2` (`owner`),
	KEY `rra_files_idx3` (`name`),
	KEY `rra_files_idx4` (`creator`),
	KEY `rra_files_idx5` (`last_modify`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBFile struct {
	ID             string
	Owner          string
	GroupID        string
	Permissions    string
	Creator        string
	CreationDate   string // DATETIME in formato stringa
	LastModify     string
	LastModifyDate string // DATETIME in formato stringa
	FatherID       string
	Name           string
	Description    string
	FKObjID        string
	Path           string
	Filename       string
	Checksum       string
	Mime           string
	AltLink        string
	DeletedBy      string
	DeletedDate    string // DATETIME in formato stringa
}

/*
CREATE TABLE IF NOT EXISTS `rra_folders` (

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
	`fk_obj_id` varchar(16) DEFAULT NULL,
	`childs_sort_order` text,
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	PRIMARY KEY (`id`),
	KEY `rra_folders_idx1` (`id`),
	KEY `rra_folders_idx2` (`owner`),
	KEY `rra_folders_idx3` (`name`),
	KEY `rra_folders_idx4` (`creator`),
	KEY `rra_folders_idx5` (`last_modify`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBFolder struct {
	ID              string
	Owner           string
	GroupID         string
	Permissions     string
	Creator         string
	CreationDate    string // DATETIME in formato stringa
	LastModify      string
	LastModifyDate  string // DATETIME in formato stringa
	FatherID        string
	Name            string
	Description     string
	FKObjID         string
	ChildsSortOrder string
	DeletedBy       string
	DeletedDate     string // DATETIME in formato stringa
}

/*
CREATE TABLE IF NOT EXISTS `rra_links` (

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
	`href` varchar(255) NOT NULL DEFAULT '',
	`target` varchar(255) DEFAULT '_blank',
	`fk_obj_id` varchar(16) DEFAULT NULL,
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	PRIMARY KEY (`id`),
	KEY `rra_links_idx1` (`id`),
	KEY `rra_links_idx2` (`owner`),
	KEY `rra_links_idx3` (`name`),
	KEY `rra_links_idx4` (`creator`),
	KEY `rra_links_idx5` (`last_modify`),
	KEY `rra_links_idx6` (`father_id`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBLink struct {
	ID             string
	Owner          string
	GroupID        string
	Permissions    string
	Creator        string
	CreationDate   string // DATETIME in formato stringa
	LastModify     string
	LastModifyDate string // DATETIME in formato stringa
	FatherID       string
	Name           string
	Description    string
	Href           string
	Target         string
	FKObjID        string
	DeletedBy      string
	DeletedDate    string // DATETIME in formato stringa
}

/*
CREATE TABLE IF NOT EXISTS `rra_notes` (

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
	`fk_obj_id` varchar(16) DEFAULT NULL,
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	PRIMARY KEY (`id`),
	KEY `rra_notes_idx1` (`id`),
	KEY `rra_notes_idx2` (`owner`),
	KEY `rra_notes_idx3` (`name`),
	KEY `rra_notes_idx4` (`creator`),
	KEY `rra_notes_idx5` (`last_modify`),
	KEY `rra_pages_idx1` (`id`),
	KEY `rra_pages_idx2` (`owner`),
	KEY `rra_pages_idx3` (`name`),
	KEY `rra_pages_idx4` (`creator`),
	KEY `rra_pages_idx5` (`last_modify`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBNote struct {
	ID             string
	Owner          string
	GroupID        string
	Permissions    string
	Creator        string
	CreationDate   string // DATETIME in formato stringa
	LastModify     string
	LastModifyDate string // DATETIME in formato stringa
	FatherID       string
	Name           string
	Description    string
	FKObjID        string
	DeletedBy      string
	DeletedDate    string // DATETIME in formato stringa
}

/*
CREATE TABLE IF NOT EXISTS `rra_pages` (
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
  `html` text,
  `fk_obj_id` varchar(16) DEFAULT NULL,
  `language` varchar(5) DEFAULT 'en_us',
  `deleted_by` varchar(16) DEFAULT NULL,
  `deleted_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/

type DBPage struct {
	ID             string
	Owner          string
	GroupID        string
	Permissions    string
	Creator        string
	CreationDate   string // DATETIME in formato stringa
	LastModify     string
	LastModifyDate string // DATETIME in formato stringa
	FatherID       string
	Name           string
	Description    string
	HTML           string
	FKObjID        string
	Language       string
	DeletedBy      string
	DeletedDate    string // DATETIME in formato stringa
}

/*
CREATE TABLE IF NOT EXISTS `rra_news` (

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
	`html` text,
	`fk_obj_id` varchar(16) DEFAULT NULL,
	`language` varchar(5) DEFAULT 'en_us',
	`deleted_by` varchar(16) DEFAULT NULL,
	`deleted_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
	PRIMARY KEY (`id`),
	KEY `rra_news_0` (`id`),
	KEY `rra_news_1` (`owner`),
	KEY `rra_news_2` (`group_id`),
	KEY `rra_news_3` (`creator`),
	KEY `rra_news_4` (`last_modify`),
	KEY `rra_news_5` (`father_id`),
	KEY `rra_news_6` (`fk_obj_id`),
	KEY `rra_news_7` (`fk_obj_id`),
	KEY `rra_news_8` (`fk_obj_id`),
	KEY `rra_news_9` (`fk_obj_id`)

) ENGINE=MyISAM DEFAULT CHARSET=latin1;
*/
type DBNews struct {
	ID             string
	Owner          string
	GroupID        string
	Permissions    string
	Creator        string
	CreationDate   string // DATETIME in formato stringa
	LastModify     string
	LastModifyDate string // DATETIME in formato stringa
	FatherID       string
	Name           string
	Description    string
	HTML           string
	FKObjID        string
	Language       string
	DeletedBy      string
	DeletedDate    string // DATETIME in formato stringa
}

/** *********************************** CMS: end. *********************************** */

/** *********************************** RRA Projects: start. *********************************** */
/** TODO */
/** *********************************** RRA Projects: end. *********************************** */
