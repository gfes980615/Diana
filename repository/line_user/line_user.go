package line_user

import (
	"fmt"
	"log"

	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/model"
)

type LineUserRepository struct {
}

// GetAllUser ...
func (lr LineUserRepository) GetAllUser() []model.LineUser {
	mysql, err := db.NewMySQL(glob.DataBase)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer mysql.Close()

	sql := fmt.Sprintf("SELECT `user_id` FROM `line_user`")

	user := []model.LineUser{}
	userResult := mysql.DB.Raw(sql).Scan(&user)
	if userResult.Error != nil {
		log.Print(userResult.Error)
		return nil
	}

	return user
}
