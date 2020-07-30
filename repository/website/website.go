package website

import (
	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/glob"
	"log"
)

type WebsiteRepository struct {
}

func (wr WebsiteRepository) Insert(url, tag string) error {
	mysql, err := db.NewMySQL(glob.DataBase)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer mysql.Close()

	e := mysql.DB.Exec("INSERT IGNORE INTO `web_site` (`url`,`tag`,`added_time`) VALUES (?,?,NOW())", url, tag)
	if err.Error != nil {
		return e.Error
	}
	return nil
}
