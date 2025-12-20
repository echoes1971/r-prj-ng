package dblayer

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

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
	DBObject
}

func NewDBEvent() *DBEvent {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "varchar(9)", Constraints: []string{"NOT NULL"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{}},
		{Name: "fk_obj_id", Type: "varchar(16)", Constraints: []string{}},

		{Name: "start_date", Type: "datetime", Constraints: []string{"NOT NULL"}},
		{Name: "end_date", Type: "datetime", Constraints: []string{"NOT NULL"}},

		{Name: "all_day", Type: "char(1)", Constraints: []string{"NOT NULL"}}, // Bool - An all day event?

		{Name: "url", Type: "varchar(255)", Constraints: []string{}}, // An Url associated to the event

		{Name: "alarm", Type: "char(1)", Constraints: []string{}},        // Bool - Signal an alarm before?
		{Name: "alarm_minute", Type: "int(11)", Constraints: []string{}}, // Num. time units
		{Name: "alarm_unit", Type: "char(1)", Constraints: []string{}},   // Time unit 0-2 => minutes, hours, days
		{Name: "before_event", Type: "char(1)", Constraints: []string{}}, // 0=before event starts 1=after

		{Name: "category", Type: "varchar(255)", Constraints: []string{}}, // Event category, e.g. Personal, Work, etc.

		{Name: "recurrence", Type: "char(1)", Constraints: []string{}},      // Bool - Recurrence active?
		{Name: "recurrence_type", Type: "char(1)", Constraints: []string{}}, // 0=Daily, 1=Weekly, 2=monthly, 3=yearly
		// 0: daily
		{Name: "daily_every_x", Type: "int(11)", Constraints: []string{}}, // every_x_days
		// 1: weekly
		{Name: "weekly_every_x", Type: "int(11)", Constraints: []string{}},         // every x weeks
		{Name: "weekly_day_of_the_week", Type: "char(1)", Constraints: []string{}}, // 0=monday ... 6=sunday
		// 2: monthly
		{Name: "monthly_every_x", Type: "int(11)", Constraints: []string{}}, // every x months
		//  2.1: n-th day of the month
		{Name: "monthly_day_of_the_month", Type: "int(11)", Constraints: []string{}}, // 0=do not, -5...-1,1 ... 31
		//  2.2: n-th week on monday
		{Name: "monthly_week_number", Type: "int(11)", Constraints: []string{}}, // 0=do not, 1...5
		{Name: "monthly_week_day", Type: "char(1)", Constraints: []string{}},    // 0=monday ... 6=sunday
		// 3: yearly
		//  3.1: every day XX of month MM
		{Name: "yearly_month_number", Type: "int(11)", Constraints: []string{}}, // 0=do not, 1...12
		{Name: "yearly_month_day", Type: "int(11)", Constraints: []string{}},    // 0=do not, 1...31
		//  3.2: every first monday of june
		{Name: "yearly_month_number", Type: "int(11)", Constraints: []string{}}, // 0=do not, 1...12
		{Name: "yearly_week_number", Type: "int(11)", Constraints: []string{}},  // 0=do not 1...5
		{Name: "yearly_week_day", Type: "char(1)", Constraints: []string{}},     // 0=monday ... 6=sunday
		// 3.3: every n-th day of the year
		{Name: "yearly_day_of_the_year", Type: "int(11)", Constraints: []string{}}, // 0=do not, 1...366
		// Recurrence range
		{Name: "recurrence_times", Type: "int(11)", Constraints: []string{}}, // 0=always 1...N times
		// Recurrence until <date>
		{Name: "recurrence_end_date", Type: "datetime", Constraints: []string{"NOT NULL"}}, // 0=always 1...N times
	}
	keys := []string{"id"}
	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "folders", RefColumn: "id"},
		{Column: "father_id", RefTable: "pages", RefColumn: "id"},
		{Column: "father_id", RefTable: "news", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "companies", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "folders", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "people", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "projects", RefColumn: "id"},
	}
	return &DBEvent{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBEvent",
				"events",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbEvent *DBEvent) NewInstance() DBEntityInterface {
	return NewDBEvent()
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
	DBObject
}

func NewDBFile() *DBFile {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "varchar(9)", Constraints: []string{"NOT NULL"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{}},
		{Name: "fk_obj_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "path", Type: "text", Constraints: []string{}},
		{Name: "filename", Type: "text", Constraints: []string{"NOT NULL"}},
		{Name: "checksum", Type: "varchar(40)", Constraints: []string{}},
		{Name: "mime", Type: "varchar(255)", Constraints: []string{}},
		{Name: "alt_link", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{}},
	}
	keys := []string{"id"}
	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "objects", RefColumn: "id"},
		{Column: "father_id", RefTable: "folders", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "pages", RefColumn: "id"},
		{Column: "father_id", RefTable: "pages", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "news", RefColumn: "id"},
		{Column: "father_id", RefTable: "news", RefColumn: "id"},
	}
	return &DBFile{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBFile",
				"files",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbFile *DBFile) NewInstance() DBEntityInterface {
	return NewDBFile()
}

func (dbFile *DBFile) GetOrderBy() []string {
	return []string{"path", "filename"}
}

func (dbFile *DBFile) generateFilename(aId any, aFilename any) string {
	var filename string
	if aFilename == nil {
		filename = dbFile.GetValue("filename").(string)
	} else {
		filename = aFilename.(string)
	}
	var id any
	if aId == nil {
		id = dbFile.GetValue("id")
	} else {
		id = aId
	}
	prefix := "r_" + id.(string) + "_"
	if len(filename) >= len(prefix) && filename[0:len(prefix)] == prefix {
		filename = filename[len(prefix):]
	}
	return prefix + filename
}

// function generaFilename($aId=null, $aFilename=null) {
// 	$nomefile = $aFilename==null?$this->getValue('filename'):$aFilename;
// 	$id=$aId==null?$this->getValue('id'):$aId;
// 	$prefisso = 'r_'.$id.'_';
// 	if(strpos($nomefile,$prefisso)!==false)
// 		$nomefile=str_replace($prefisso,"",$nomefile);
// 	return $prefisso.$nomefile;
// }

func (dbFile *DBFile) generateObjectPath(a_dbe DBEntityInterface) string {
	var dest_path string
	mydbe := a_dbe
	if mydbe == nil {
		mydbe = dbFile
	}
	if mydbe.GetValue("path") != nil && mydbe.GetValue("path").(string) != "" {
		dest_path = mydbe.GetValue("path").(string)
	} else {
		dest_path = ""
	}
	father_id := mydbe.GetValue("father_id")
	if father_id != nil && father_id.(string) != "" && father_id.(string) != "0" {
		if dest_path != "" {
			dest_path = father_id.(string) + "/" + dest_path
		} else {
			dest_path = father_id.(string)
		}
	}
	return dest_path
}

// function generaObjectPath($a_dbe=null) {
// 	$_dbe = $a_dbe!=null ? $a_dbe : $this;
// 	$dest_path = $_dbe->getValue('path')>'' ? $_dbe->getValue('path') : '';
// 	$father_id = $_dbe->getValue('father_id');
// 	if($father_id>0) $dest_path = $father_id.($dest_path>''?'/':'').$dest_path;
// 	return $dest_path;
// }

func (dbFile *DBFile) GetFullpath(a_dbe DBEntityInterface) string {
	var mydbe *DBFile
	if a_dbe != nil {
		mydbe = a_dbe.(*DBFile)
	} else {
		mydbe = dbFile
	}
	dest_path := mydbe.generateObjectPath(nil)
	// $dest_dir=realpath($GLOBALS['root_directory'].'/'.$mydbe->dest_directory);
	dest_dir := dbFiles_root_directory + "/" + dbFiles_dest_directory
	if dest_path != "" {
		dest_dir = dest_dir + "/" + dest_path
	}
	ret := dest_dir + "/" + mydbe.GetValue("filename").(string)
	return ret
}

// function getFullpath($a_dbe=null) {
// 	$mydbe = $a_dbe!=null ? $a_dbe : $this;
// 	$dest_path = $mydbe->generaObjectPath();
// 	$dest_dir=realpath($GLOBALS['root_directory'].'/'.$mydbe->dest_directory);
// 	if($dest_path>'') $dest_dir.="/$dest_path";
// 	$ret = "$dest_dir/".$mydbe->getValue('filename');
// 	return $ret;
// }

// // Image management: start.

// function getThumbnailFilename() { return $this->getValue('filename')."_thumb.jpg"; }
func (dbFile *DBFile) getThumbnailFilename(fullpath string) string {
	if fullpath != "" {
		return fullpath + "_thumb.jpg"
	}
	return dbFile.GetValue("filename").(string) + "_thumb.jpg"
}
func (dbFile *DBFile) GetThumbnailFullpath(a_dbe DBEntityInterface) string {
	var mydbe *DBFile
	if a_dbe != nil {
		mydbe = a_dbe.(*DBFile)
	} else {
		mydbe = dbFile
	}
	dest_path := mydbe.GetFullpath(nil)
	thumbPath := mydbe.getThumbnailFilename(dest_path)
	return thumbPath
}

// function IsImage() { $_mime = $this->getValue('mime'); return $_mime>” && substr($_mime,0,5)=='image'; }
func (dbFile *DBFile) IsImage() bool {
	mime := dbFile.GetValue("mime")
	if mime != nil && mime.(string) != "" && len(mime.(string)) >= 5 && mime.(string)[0:5] == "image" {
		return true
	}
	return false
}

// computeSHA1 calculates the SHA1 hash of a file
func (dbFile *DBFile) computeSHA1(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// detectMimeType detects the MIME type of a file (equivalent to finfo_open in PHP)
func (dbFile *DBFile) detectMimeType(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "application/octet-stream", err
	}
	defer file.Close()

	// Read first 512 bytes for detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "application/octet-stream", err
	}

	// Use http.DetectContentType (equivalent to finfo_open(FILEINFO_MIME) in PHP)
	mimeType := http.DetectContentType(buffer[:n])
	return mimeType, nil
}

func (dbFile *DBFile) createThumbnail(fullpath string) string {
	thumbPath := fullpath + "_thumb.jpg"

	// Open file
	file, err := os.Open(fullpath)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", fullpath, err)
		return ""
	}
	defer file.Close()

	// Decode image (supports JPEG, PNG, GIF)
	img, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Error decoding image %s: %v\n", fullpath, err)
		return ""
	}

	// Calculate dimensions maintaining aspect ratio
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())
	maxDim := width
	if height > width {
		maxDim = height
	}

	// Resize to fit within 100x100 (same as PHP version)
	var thumb image.Image
	if maxDim > 100 {
		thumb = resize.Thumbnail(100, 100, img, resize.Lanczos3)
	} else {
		thumb = img // No need to resize if already small
	}

	// Create output file
	out, err := os.Create(thumbPath)
	if err != nil {
		log.Printf("Error creating thumbnail file %s: %v\n", thumbPath, err)
		return ""
	}
	defer out.Close()

	// Save as JPEG with quality 85
	err = jpeg.Encode(out, thumb, &jpeg.Options{Quality: 85})
	if err != nil {
		log.Printf("Error encoding thumbnail %s: %v\n", thumbPath, err)
		return ""
	}

	log.Printf("Created thumbnail for %s (format: %s) at %s\n", fullpath, format, thumbPath)
	return thumbPath
}

func (dbFile *DBFile) deleteThumbnail(fullpath string) {
	thumbPath := fullpath + "_thumb.jpg"
	err := os.Remove(thumbPath)
	if err != nil {
		log.Printf("Error deleting thumbnail at %s: %v\n", thumbPath, err)
	} else {
		log.Printf("Deleted thumbnail at %s\n", thumbPath)
	}
}

// function createThumbnail($fullpath, $pix_width=100,$pix_height=100) {
// 	$gis = getimagesize($fullpath);
// 	$type = $gis[2];
// 	$imorig=null;
// 		if(!function_exists('imagecreatefromjpeg')) {
// 			echo "<h1>";
// 			echo "RUN:<br/>";
// 			echo "sudo aptitude install php5-gd</br>";
// 			echo "sudo /etc/init.d/apache2 restart<br/>";
// 			echo "</h1>";
// 		}
// 	switch($type) {
// 		case "1": $imorig = imagecreatefromgif($fullpath); break;
// 		case "2": $imorig = imagecreatefromjpeg($fullpath);break;
// 		case "3": $imorig = imagecreatefrompng($fullpath); break;
// 		default:  $imorig = imagecreatefromjpeg($fullpath);
// 	}
// 	$w = imagesx($imorig);
// 	$h = imagesy($imorig);
// 	$max_pixel = $w>$h ? $w : $h;
// 	$scale = $max_pixel / ($w>$h ? $pix_width : $pix_height);
// 	$pix_width = intval($w / $scale);
// 	$pix_height = intval($h / $scale);
// 	$im = imagecreatetruecolor($pix_width,$pix_height);
// 	if(imagecopyresampled($im,$imorig , 0,0,0,0,$pix_width,$pix_height,$w,$h)) {
// 		if(imagejpeg($im, $fullpath."_thumb.jpg")) {
// 			return $fullpath."_thumb.jpg";
// 		}
// 	}
// 	return "";
// }
// function deleteThumbnail($fullpath) {
// 	unlink($fullpath."_thumb.jpg");
// }
// // Image management: end.

func (dbFile *DBFile) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	if dbr.Verbose {
		log.Print("DBFile.beforeInsert called")
	}
	err := dbFile.DBObject.beforeInsert(dbr, tx)
	if err != nil {
		return err
	}
	// This seems redundant with SetDefaultValues, but keeping it for compatibility
	// I don't know, it seems to hide the effects of SetDefaultValues... maybe it should be removed?
	fatherId := dbFile.GetValue("father_id")
	if dbFile.HasValue("father_id") && fatherId != nil && fatherId != "" && fatherId != "0" {
		father := dbr.GetEntityByIDWithTx("folders", fatherId.(string), tx)
		if father != nil {
			if fatherFolder, ok := father.(*DBFolder); ok {
				if fatherFolder.HasValue("fk_obj_id") && fatherFolder.GetValue("fk_obj_id") != "" && fatherFolder.GetValue("fk_obj_id") != "0" {
					dbFile.SetValue("fk_obj_id", fatherFolder.GetValue("fk_obj_id"))
				}
			}
		}
	}
	// Adding prefix to filename
	if dbFile.GetValue("filename") != nil && dbFile.GetValue("filename").(string) != "" {
		dest_path := dbFile.generateObjectPath(nil)
		from_dir := dbFiles_root_directory + "/" + dbFiles_dest_directory
		dest_dir := dbFiles_root_directory + "/" + dbFiles_dest_directory
		if dest_path != "" {
			dest_dir = dest_dir + "/" + dest_path
		}
		// Create destination directory if it does not exist
		os.MkdirAll(dest_dir, os.FileMode(0755))
		// Using basename equivalent
		new_filename := dbFile.generateFilename(dbFile.GetValue("id"), filepath.Base(dbFile.GetValue("filename").(string)))
		err := os.Rename(from_dir+"/"+dbFile.GetValue("filename").(string), dest_dir+"/"+new_filename)
		if err != nil {
			return err
		}
		if dbFile.GetValue("name") == nil || dbFile.GetValue("name").(string) == "" {
			dbFile.SetValue("name", filepath.Base(dbFile.GetValue("filename").(string)))
		}
		dbFile.SetValue("filename", new_filename)
	}
	if dbFile.GetValue("filename") == nil || strings.TrimSpace(dbFile.GetValue("filename").(string)) == "" {
		// return fmt.Errorf("filename cannot be empty after processing")
		dbFile.SetValue("filename", "")
		return nil
	}
	// Checksum
	fullpath := dbFile.GetFullpath(nil)
	if _, err := os.Stat(fullpath); err == nil {
		// File exists
		checksum, err := dbFile.computeSHA1(fullpath)
		if err != nil {
			return err
		}
		dbFile.SetValue("checksum", checksum)
	} else {
		dbFile.SetValue("checksum", "File '"+dbFile.GetValue("filename").(string)+"' not found!")
	}
	// Mime type
	if _, err := os.Stat(fullpath); err == nil {
		mimeType, err := dbFile.detectMimeType(fullpath)
		if err != nil {
			return err
		}
		dbFile.SetValue("mime", mimeType)
	} else {
		dbFile.SetValue("mime", "text/plain")
	}
	// Image
	if dbFile.IsImage() {
		dbFile.createThumbnail(fullpath)
	}
	return nil
}

// function _before_insert(&$dbmgr) {
// 	parent::_before_insert($dbmgr);
// 	// Eredita la 'radice' dal padre
// 	$father_id = $this->getValue('father_id');
// 	if($father_id>0) {
// 		$query="select fk_obj_id from ". $dbmgr->buildTableName($this)." where id='".DBEntity::uuid2hex($this->getValue('father_id'))."'";
// 		$tmp = $dbmgr->select("DBE",$this->getTableName(),$query);
// 		if(count($tmp)==1) {
// 			$this->setValue('fk_obj_id', $tmp[0]->getValue('fk_obj_id'));
// 		}
// 	}
// 	// Aggiungo il prefisso al nome del file
// 	if($this->getValue('filename')>'') {
// 		$dest_path = $this->generaObjectPath();
// 		$from_dir=realpath($GLOBALS['root_directory'].'/'.$this->dest_directory);
// 		$dest_dir=realpath($GLOBALS['root_directory'].'/'.$this->dest_directory);
// 		if($dest_path>'') $dest_dir.="/$dest_path";
// 		if(!file_exists($dest_dir)) mkdir($dest_dir, 0755);
// 		// con basename() ottengo solo il nome del file senza il path relativo nel quale e' stato caricato
// 		$nuovo_filename = $this->generaFilename($this->getValue('id'), basename($this->getValue('filename')));
// 		rename($from_dir."/".$this->getValue('filename'), $dest_dir."/".$nuovo_filename);
// 		if(!($this->getValue('name')>'')) $this->setValue('name',basename($this->getValue('filename')));
// 		$this->setValue('filename', $nuovo_filename);
// 	}
// 	// Checksum
// 	$_fullpath = $this->getFullpath();
// 	if(file_exists($_fullpath)) {
// 		$newchecksum = sha1_file($_fullpath);
// 		$this->setValue('checksum',$newchecksum);
// 	} else {
// 		$this->setValue('checksum',"File '".$this->getValue('filename')."' not found!");
// 	}
// 	// Mime type
// 	if(file_exists($_fullpath)) {
// 		if(function_exists('finfo_open')) {
// 			$finfo = finfo_open(FILEINFO_MIME);
// 			if(!$finfo) {
// 				if(function_exists('mime_content_type'))
// 					$this->setValue('mime',mime_content_type($_fullpath));
// 				else
// 					$this->setValue('mime','text/plain');
// 				return;
// 			}
// 			$this->setValue('mime',finfo_file($finfo,$_fullpath));
// 			finfo_close($finfo);
// 		} elseif(function_exists('mime_content_type'))
// 			$this->setValue('mime',mime_content_type($_fullpath));
// 		else
// 			$this->setValue('mime','text/plain');
// 	} else {
// 		$this->setValue('mime','text/plain');
// 	}
// 	// Image
// 	if($this->isImage())
// 		$this->createThumbnail($_fullpath);
// }

func (dbFile *DBFile) beforeUpdate(dbr *DBRepository, tx *sql.Tx) error {
	if dbr.Verbose {
		log.Print("DBFile.beforeUpdate called")
	}
	err := dbFile.DBObject.beforeUpdate(dbr, tx)
	if err != nil {
		return err
	}
	// Inherit 'root' from parent
	fatherId := dbFile.GetValue("father_id")
	if dbFile.HasValue("father_id") && fatherId != nil && fatherId != "" && fatherId != "0" {
		father := dbr.GetEntityByIDWithTx("folders", fatherId.(string), tx)
		if father != nil {
			if fatherFolder, ok := father.(*DBFolder); ok {
				if fatherFolder.HasValue("fk_obj_id") && fatherFolder.GetValue("fk_obj_id") != "" && fatherFolder.GetValue("fk_obj_id") != "0" {
					dbFile.SetValue("fk_obj_id", fatherFolder.GetValue("fk_obj_id"))
				}
			}
		}
	}
	// Check if I already have a saved file
	myself := dbr.GetEntityByIDWithTx("files", dbFile.GetValue("id").(string), tx).(*DBFile)
	if myself == nil {
		// Error: should not happen
		return nil
	}
	myself_has_a_file := myself.GetValue("filename") != nil && myself.GetValue("filename").(string) != ""
	// dbFile_has_a_file := dbFile.GetValue("filename") != nil && dbFile.GetValue("filename").(string) != ""
	// If I had a file and now I don't, delete the old one
	if myself_has_a_file {
		// TODO very ugly nesting
		if dbFile.GetValue("filename") != nil && dbFile.GetValue("filename").(string) != "" && myself.GetValue("filename").(string) != dbFile.GetValue("filename").(string) {
			// Different filenames ==> delete the old one
			dest_path := myself.generateObjectPath(nil)
			dest_dir := dbFiles_root_directory + "/" + dbFiles_dest_directory
			if dest_path != "" {
				dest_dir = dest_dir + "/" + dest_path
			}
			dest_file := dest_dir + "/" + myself.GetValue("filename").(string)
			if _, err := os.Stat(dest_file); os.IsNotExist(err) {
				// Do nothing
			} else {
				err := os.Remove(dest_file)
				if err != nil {
					return err
				}
				// Image
				if dbFile.IsImage() {
					dbFile.deleteThumbnail(dest_file)
				}
			}
		}
	}
	// Adding prefix to filename
	if dbFile.GetValue("filename") != nil && dbFile.GetValue("filename").(string) != "" {
		from_dir := dbFiles_root_directory + "/" + dbFiles_dest_directory
		dest_path := dbFile.generateObjectPath(nil)
		dest_dir := dbFiles_root_directory + "/" + dbFiles_dest_directory
		if dest_path != "" {
			dest_dir = dest_dir + "/" + dest_path
		}
		// Create destination directory if it does not exist
		os.MkdirAll(dest_dir, os.FileMode(0755))
		new_filename := dbFile.generateFilename(dbFile.GetValue("id"), filepath.Base(dbFile.GetValue("filename").(string)))
		log.Print("DBFile.beforeUpdate: moving file from ", from_dir+"/"+myself.GetValue("filename").(string), " to ", dest_dir+"/"+new_filename)
		// Move the file only if it exists
		if _, err := os.Stat(from_dir + "/" + dbFile.GetValue("filename").(string)); err == nil {
			err := os.Rename(from_dir+"/"+dbFile.GetValue("filename").(string), dest_dir+"/"+new_filename)
			if err != nil {
				log.Print("DBFile.beforeUpdate: error renaming file: ", err)
				return err
			}
		}
		dbFile.SetValue("filename", new_filename)
		// }
	} else if myself_has_a_file && myself.GetValue("path") != dbFile.GetValue("path") {
		// if myself_has_a_file && myself.GetValue("path") != dbFile.GetValue("path") {
		// } else if myself.GetValue("filename") != nil && myself.GetValue("filename").(string) != "" && myself.GetValue("path") != dbFile.GetValue("path") {
		from_path := myself.generateObjectPath(nil)
		from_dir := dbFiles_root_directory + "/" + dbFiles_dest_directory
		if from_path != "" {
			from_dir = from_dir + "/" + from_path
		}
		dest_path := dbFile.generateObjectPath(nil)
		dest_dir := dbFiles_root_directory + "/" + dbFiles_dest_directory
		if dest_path != "" {
			dest_dir = dest_dir + "/" + dest_path
		}
		// Create destination directory if it does not exist
		os.MkdirAll(dest_dir, os.FileMode(0755))
		err := os.Rename(from_dir+"/"+myself.GetValue("filename").(string), dest_dir+"/"+myself.GetValue("filename").(string))
		if err != nil {
			return err
		}
		// TODO check if it works
		dbFile.SetValue("filename", myself.GetValue("filename"))
	} else if myself_has_a_file {
		// TODO check if it works
		dbFile.SetValue("filename", myself.GetValue("filename"))
	}

	// Check if father_id has changed and move file accordingly
	if dbFile.GetValue("father_id") != nil && dbFile.GetValue("father_id") != myself.GetValue("father_id") {
		from_path := myself.generateObjectPath(nil)
		from_dir := dbFiles_root_directory + "/" + dbFiles_dest_directory
		if from_path != "" {
			from_dir = from_dir + "/" + from_path
		}
		dest_path := dbFile.generateObjectPath(nil)
		dest_dir := dbFiles_root_directory + "/" + dbFiles_dest_directory
		if dest_path != "" {
			dest_dir = dest_dir + "/" + dest_path
		}
		// Create destination directory if it does not exist
		os.MkdirAll(dest_dir, os.FileMode(0755))
		err := os.Rename(from_dir+"/"+myself.GetValue("filename").(string), dest_dir+"/"+myself.GetValue("filename").(string))
		if err != nil {
			return err
		}
	}

	// Checksum
	fullpath := dbFile.GetFullpath(nil)
	if _, err := os.Stat(fullpath); err == nil {
		// File exists
		checksum, err := dbFile.computeSHA1(fullpath)
		if err != nil {
			return err
		}
		dbFile.SetValue("checksum", checksum)
	} else {
		dbFile.SetValue("checksum", "File '"+dbFile.GetValue("filename").(string)+"' not found!")
	}
	// Mime type
	if _, err := os.Stat(fullpath); err == nil {
		mimeType, err := dbFile.detectMimeType(fullpath)
		if err != nil {
			return err
		}
		dbFile.SetValue("mime", mimeType)
	} else {
		dbFile.SetValue("mime", "text/plain")
	}
	// Image
	if dbFile.IsImage() {
		dbFile.createThumbnail(fullpath)
	}
	return nil
}

// function _before_update(&$dbmgr) {
// 	parent::_before_update($dbmgr);
// 	// Eredita la 'radice' dal padre
// 	$father_id = $this->getValue('father_id');
// 	if($father_id>0) {
// 		$query="select fk_obj_id from ". $dbmgr->buildTableName($this)." where id='".DBEntity::uuid2hex($this->getValue('father_id'))."'";
// 		$tmp = $dbmgr->select("DBE",$this->getTableName(),$query);
// 		if(count($tmp)==1) {
// 			$this->setValue('fk_obj_id', $tmp[0]->getValue('fk_obj_id'));
// 		}
// 	}
// 	// Controllo se ho già un file salvato
// 	eval("\$cerca = new ".get_class($this)."();");
// 	$cerca->setValue('id', $this->getValue('id'));
// 	$tmp=$dbmgr->search($cerca,$uselike=0);
// 	$myself=$tmp[0];
// 	if($this->getValue('filename')>'' && $myself->getValue('filename')!=$this->getValue('filename')) {
// 		// Filename diversi ==> elimino il vecchio
// 		$dest_path = $myself->generaObjectPath();
// 		$dest_dir=realpath($GLOBALS['root_directory'].'/'.$this->dest_directory);
// 		if($dest_path>'') $dest_dir.="/$dest_path";
// 		$dest_file = $dest_dir."/".$myself->generaFilename();
// 		if(!file_exists($dest_file)) {
// 			// Do nothing
// 		} else {
// 			unlink($dest_file);
// 			// Image
// 			if($this->isImage())
// 				$this->deleteThumbnail($dest_file);
// 		}
// 	}
// 	// Aggiungo il prefisso al nome del file
// 	if($this->getValue('filename')>'') {
// 		$from_dir=realpath($GLOBALS['root_directory'].'/'.$this->dest_directory);
// 		$dest_path = $this->generaObjectPath();
// 		$dest_dir=realpath($GLOBALS['root_directory'].'/'.$this->dest_directory);
// 		if($dest_path>'') $dest_dir.="/$dest_path";
// 		if(!file_exists($dest_dir)) mkdir($dest_dir, 0755);
// 		$nuovo_filename = $this->generaFilename($this->getValue('id'), basename($this->getValue('filename')));
// 		rename("$from_dir/".$this->getValue('filename'),"$dest_dir/$nuovo_filename");
// 		$this->setValue('filename', $nuovo_filename);
// 	} else if($myself->getValue('path')!=$this->getValue('path')) {
// 		$from_path = $myself->generaObjectPath();
// 		$from_dir=realpath($GLOBALS['root_directory'].'/'.$this->dest_directory);
// 		if($from_path>'') $from_dir.="/$from_path";
// 		$dest_path = $this->generaObjectPath();
// 		$dest_dir=realpath($GLOBALS['root_directory'].'/'.$this->dest_directory);
// 		if($dest_path>'') $dest_dir.="/$dest_path";
// 		if(!file_exists($dest_dir)) mkdir($dest_dir, 0755);
// 		rename("$from_dir/".$myself->getValue('filename'),"$dest_dir/".$myself->getValue('filename'));
// 		// TODO controllare se funziona
// 		$this->setValue('filename', $myself->getValue('filename'));
// 	} else {
// 		// TODO controllare se funziona
// 		$this->setValue('filename', $myself->getValue('filename'));
// 	}
// 	// Checksum
// 	$_fullpath = $this->getFullpath();
// 	if(file_exists($_fullpath)) {
// 		$newchecksum = sha1_file($_fullpath);
// 		$this->setValue('checksum',$newchecksum);
// 	} else {
// 		$this->setValue('checksum',"File '".$this->getValue('filename')."' not found!");
// 	}
// 	// Mime type
// 	if(file_exists($_fullpath)) {
// 		if(function_exists('finfo_open')) {
// 			$finfo = finfo_open(FILEINFO_MIME);
// 			if(!$finfo) {
// 				if(function_exists('mime_content_type'))
// 					$this->setValue('mime',mime_content_type($_fullpath));
// 				else
// 					$this->setValue('mime','text/plain');
// 				return;
// 			}
// 			$this->setValue('mime',finfo_file($finfo,$_fullpath));
// 			finfo_close($finfo);
// 		} elseif(function_exists('mime_content_type'))
// 			$this->setValue('mime',mime_content_type($_fullpath));
// 		else
// 			$this->setValue('mime','text/plain');
// 	} else {
// 		$this->setValue('mime','text/plain');
// 	}
// 	// Image
// 	if($this->isImage())
// 		$this->createThumbnail($_fullpath);
// }

func (dbFile *DBFile) beforeDelete(dbr *DBRepository, tx *sql.Tx) error {
	if dbr.Verbose {
		log.Print("DBFile.beforeDelete called")
	}
	is_deleted := dbFile.HasDeletedDate()

	// If it has been marked deleted, then now is a REAL delete, so remove the file BEFORE calling parent
	if is_deleted {
		// Use current dbFile values (not from DB, as it will be deleted)
		if dbFile.GetValue("filename") != nil && dbFile.GetValue("filename").(string) != "" {
			// ==> delete the file
			dest_path := dbFile.generateObjectPath(nil)
			dest_dir := dbFiles_root_directory + "/" + dbFiles_dest_directory
			if dest_path != "" {
				dest_dir = dest_dir + "/" + dest_path
			}
			fullpath := dest_dir + "/" + dbFile.GetValue("filename").(string)
			err := os.Remove(fullpath)
			if err != nil && !os.IsNotExist(err) {
				log.Printf("Error removing file %s: %v", fullpath, err)
			} else {
				log.Printf("Deleted file at %s", fullpath)
			}
			// Image
			if dbFile.IsImage() {
				dbFile.deleteThumbnail(fullpath)
			}
		}
	}

	err := dbFile.DBObject.beforeDelete(dbr, tx)
	if err != nil {
		return err
	}
	return nil
}

// function _before_delete(&$dbmgr) {
// 	// Has it been marked deleted before?
// 	$is_deleted = $this->isDeleted();
// 	parent::_before_delete($dbmgr);
// 	// If it has been marked deleted, then now is a REAL delete, so remove the file
// 	if($is_deleted) {
// 		// Controllo se ho già un file salvato
// 		$cerca = new DBEFile();
// 		$cerca->setValue('id', $this->getValue('id'));
// 		// BUGFIX 2012.04.04: start.
// 		$tmp=$dbmgr->search($cerca,0,false,null,false);
// // 			$tmp=$dbmgr->search($cerca,$uselike=0);
// 		// BUGFIX 2012.04.04: end.
// 		if(count($tmp)>0) {
// 			$myself=$tmp[0];
// 			if($myself->getValue('filename')>'') {
// 				// ==> elimino il file
// 				$dest_path = $myself->generaObjectPath();
// 				$dest_dir=realpath($GLOBALS['root_directory'].'/'.$this->dest_directory);
// 				if($dest_path>'') $dest_dir.="/$dest_path";
// 				unlink($dest_dir."/".$myself->generaFilename());
// 				// Image
// 				if($this->isImage())
// 					$this->deleteThumbnail($dest_dir."/".$myself->generaFilename());
// 			}
// 		}
// 	}
// }

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
	DBObject
}

func NewDBFolder() *DBFolder {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "varchar(9)", Constraints: []string{"NOT NULL"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{}},
		{Name: "fk_obj_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "childs_sort_order", Type: "text", Constraints: []string{}},
	}
	keys := []string{"id"}
	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "objects", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "companies", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "people", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "projects", RefColumn: "id"},
	}
	return &DBFolder{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBFolder",
				"folders",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbFile *DBFolder) NewInstance() DBEntityInterface {
	return NewDBFolder()
}
func (dbFile *DBFolder) IsDBObject() bool {
	return true
}

func (dbFile *DBFolder) GetChildsSortOrder() []string {
	sorted_childs := []string{}
	if dbFile.HasValue("childs_sort_order") && dbFile.GetValue("childs_sort_order") != nil && dbFile.GetValue("childs_sort_order").(string) != "" {
		for _, child := range strings.Split(dbFile.GetValue("childs_sort_order").(string), ",") {
			sorted_childs = append(sorted_childs, hex2uuid(child))
		}
	}
	return sorted_childs
}

func (dbFile *DBFolder) SetDefaultValues(repo *DBRepository) {
	if repo.Verbose {
		log.Print("DBFolder.SetDefaultValues called")
	}
	dbFile.DBObject.SetDefaultValues(repo)

	if !dbFile.HasValue("father_id") || dbFile.GetValue("father_id") == "" || dbFile.GetValue("father_id") == "0" {
		return
	}
	father := repo.GetEntityByID("folders", dbFile.GetValue("father_id").(string))
	if father != nil {
		if fatherFolder, ok := father.(*DBFolder); ok {
			if fatherFolder.HasValue("fk_obj_id") && fatherFolder.GetValue("fk_obj_id") != "" && fatherFolder.GetValue("fk_obj_id") != "0" {
				dbFile.SetValue("fk_obj_id", fatherFolder.GetValue("fk_obj_id"))
			}
		}
	}
}

func (dbFolder *DBFolder) beforeInsert(dbr *DBRepository, tx *sql.Tx) error {
	if dbr.Verbose {
		log.Print("DBFolder.beforeInsert called")
	}
	err := dbFolder.DBObject.beforeInsert(dbr, tx)
	if err != nil {
		return err
	}
	// This seems redundant with SetDefaultValues, but keeping it for compatibility
	// I don't know, it seems to hide the effects of SetDefaultValues... maybe it should be removed?
	if !dbFolder.HasValue("father_id") || dbFolder.GetValue("father_id") == "" || dbFolder.GetValue("father_id") == "0" {
		return nil
	}
	father_id, ok := dbFolder.GetValue("father_id").(string)
	if ok {
		father := dbr.GetEntityByIDWithTx("folders", father_id, tx)
		if father != nil {
			if fatherFolder, ok := father.(*DBFolder); ok {
				if fatherFolder.HasValue("fk_obj_id") && fatherFolder.GetValue("fk_obj_id") != "" && fatherFolder.GetValue("fk_obj_id") != "0" {
					dbFolder.SetValue("fk_obj_id", fatherFolder.GetValue("fk_obj_id"))
				}
			}
		}
	}
	return nil
}

func (dbFile *DBFolder) beforeUpdate(dbr *DBRepository, tx *sql.Tx) error {
	if dbr.Verbose {
		log.Print("DBFolder.beforeUpdate called")
	}
	err := dbFile.DBObject.beforeUpdate(dbr, tx)
	if err != nil {
		return err
	}
	// This seems redundant with SetDefaultValues, but keeping it for compatibility
	// I don't know, it seems to hide the effects of SetDefaultValues... maybe it should be removed?
	if !dbFile.HasValue("father_id") || dbFile.GetValue("father_id") == "" || dbFile.GetValue("father_id") == "0" {
		return nil
	}
	father := dbr.GetEntityByIDWithTx("folders", dbFile.GetValue("father_id").(string), tx)
	if father != nil {
		if fatherFolder, ok := father.(*DBFolder); ok {
			if fatherFolder.HasValue("fk_obj_id") && fatherFolder.GetValue("fk_obj_id") != "" && fatherFolder.GetValue("fk_obj_id") != "0" {
				dbFile.SetValue("fk_obj_id", fatherFolder.GetValue("fk_obj_id"))
			}
		}
	}
	return nil
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
	DBObject
}

func NewDBLink() *DBLink {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "varchar(9)", Constraints: []string{"NOT NULL"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{}},
		{Name: "href", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "target", Type: "varchar(255)", Constraints: []string{}},
		{Name: "fk_obj_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{}},
	}
	keys := []string{"id"}

	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "folders", RefColumn: "id"},
		{Column: "father_id", RefTable: "pages", RefColumn: "id"},
		{Column: "father_id", RefTable: "news", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "companies", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "folders", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "people", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "projects", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "pages", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "news", RefColumn: "id"},
	}
	return &DBLink{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBLink",
				"links",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbLink *DBLink) NewInstance() DBEntityInterface {
	return NewDBLink()
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
	DBObject
}

func NewDBNote() *DBNote {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "varchar(9)", Constraints: []string{"NOT NULL"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{}},
		{Name: "fk_obj_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{}},
	}
	keys := []string{"id"}

	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "objects", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "companies", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "folders", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "people", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "projects", RefColumn: "id"},
	}
	return &DBNote{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBNote",
				"notes",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbNote *DBNote) NewInstance() DBEntityInterface {
	return NewDBNote()
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
	DBObject
}

func NewDBPage() *DBPage {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "varchar(9)", Constraints: []string{"NOT NULL"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{}},
		{Name: "html", Type: "text", Constraints: []string{}},
		{Name: "fk_obj_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "language", Type: "varchar(5)", Constraints: []string{}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{}},
	}
	keys := []string{"id"}

	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "objects", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "companies", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "folders", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "people", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "projects", RefColumn: "id"},
	}
	return &DBPage{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBPage",
				"pages",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbPage *DBPage) NewInstance() DBEntityInterface {
	return NewDBPage()
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
	DBObject
}

func NewDBNews() *DBNews {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "owner", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "permissions", Type: "varchar(9)", Constraints: []string{"NOT NULL"}},
		{Name: "creator", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "creation_date", Type: "datetime", Constraints: []string{}},
		{Name: "last_modify", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "last_modify_date", Type: "datetime", Constraints: []string{}},
		{Name: "deleted_by", Type: "varchar(16)", Constraints: []string{}},
		{Name: "deleted_date", Type: "datetime", Constraints: []string{}},
		{Name: "father_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "description", Type: "text", Constraints: []string{}},
		{Name: "html", Type: "text", Constraints: []string{}},
		{Name: "fk_obj_id", Type: "varchar(16)", Constraints: []string{}},
		{Name: "language", Type: "varchar(5)", Constraints: []string{}},
	}
	keys := []string{"id"}

	foreignKeys := []ForeignKey{
		{Column: "owner", RefTable: "users", RefColumn: "id"},
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
		{Column: "creator", RefTable: "users", RefColumn: "id"},
		{Column: "last_modify", RefTable: "users", RefColumn: "id"},
		{Column: "deleted_by", RefTable: "users", RefColumn: "id"},
		{Column: "father_id", RefTable: "objects", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "companies", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "folders", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "people", RefColumn: "id"},
		{Column: "fk_obj_id", RefTable: "projects", RefColumn: "id"},
	}
	return &DBNews{
		DBObject: DBObject{
			DBEntity: *NewDBEntity(
				"DBNews",
				"news",
				columns,
				keys,
				foreignKeys,
				make(map[string]any),
			),
		},
	}
}
func (dbNews *DBNews) NewInstance() DBEntityInterface {
	return NewDBNews()
}
