package redis

import (
	"context"
	"encoding/json"
	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/dto"
	"github.com/go-redis/redis/v8"
)

func init() {
	injection.AutoRegister(&mapleRedisRepository{})
}

type mapleRedisRepository struct {
}

func (mrr *mapleRedisRepository) GetBulletinData() ([]*dto.MapleBulletin, error) {
	val, err := db.RDB.LPop(context.Background(), "maple_bulletin").Bytes()
	if err != nil {
		if err == redis.Nil {
			return []*dto.MapleBulletin{}, nil
		}
		return nil, err
	}
	result := []*dto.MapleBulletin{}
	err = json.Unmarshal(val, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
