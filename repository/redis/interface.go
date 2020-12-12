package redis

import "github.com/gfes980615/Diana/models/dto"

type MapleRedisRepository interface {
	GetBulletinData() ([]*dto.MapleBulletin, error)
}
