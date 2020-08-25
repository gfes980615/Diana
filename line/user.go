package line
//
//import (
//	"fmt"
//	"log"
//
//	"github.com/gfes980615/Diana/db"
//	"github.com/gfes980615/Diana/glob"
//	"github.com/gfes980615/Diana/models"
//)
//
//// SaveUserID TODO 需要再想想看怎麼寫比較好
//func SaveUserID(id string) {
//	mysql, err := db.NewMySQL(glob.DataBase)
//	if err != nil {
//		log.Print(err)
//	}
//	defer mysql.Close()
//
//	sql := fmt.Sprintf("SELECT `user_id` FROM `line_user` WHERE `user_id` = '%s'", id)
//	user := []models.LineUser{}
//	userResult := mysql.DB.Raw(sql).Scan(&user)
//	if userResult.Error != nil {
//		log.Print(userResult.Error)
//	}
//
//	// 沒存過的會員才儲存
//	if len(user) == 0 {
//		result := mysql.DB.Exec("INSERT INTO `line_user` (`user_id`,`added_time`) VALUES (?,NOW())", id)
//		if result.Error != nil {
//			log.Print(result.Error)
//		}
//	}
//}
