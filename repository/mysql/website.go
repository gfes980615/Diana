package mysql

type WebsiteRepository struct {
}

//func (wr WebsiteRepository) Insert(url, tag string) error {
//	mysql, err := db.NewMySQL(glob.DataBase)
//	if err != nil {
//		log.Print(err)
//		return nil
//	}
//	defer mysql.Close()
//
//	e := mysql.DB.Exec("INSERT INTO `web_site` (`url`,`tag`,`added_time`) VALUE (?,?,NOW())", url, tag)
//	if e.Error != nil {
//		return e.Error
//	}
//	return nil
//}
