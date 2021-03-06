package storage

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"

	"github.com/TradeWars/warehouse/types"
)

func (mgr *Manager) ensureAdminCollection() (err error) {
	mgr.admins = mgr.db.C("admins")

	err = mgr.admins.EnsureIndex(mgo.Index{
		Key:    []string{"player_id"},
		Unique: true,
	})
	if err != nil {
		return
	}

	return
}

// AdminSetLevel creates, updates or removes an admin record based on level
func (mgr *Manager) AdminSetLevel(id bson.ObjectId, level int32) (err error) {
	n, err := mgr.admins.Find(bson.M{"player_id": id}).Count()
	if err != nil {
		return errors.Wrap(err, "failed to check if admin already exists")
	}

	if n == 0 {
		if level == 0 {
			return
		}

		err = mgr.admins.Insert(types.Admin{
			ID:       bson.NewObjectId(),
			PlayerID: id,
			Level:    &level,
			Date:     time.Now(),
		})
	} else {
		if level == 0 {
			err = mgr.admins.Remove(bson.M{"player_id": id})
		} else {
			err = mgr.admins.Update(bson.M{"player_id": id}, bson.M{"$set": bson.M{"level": level}})
		}
	}
	return
}

// AdminGetList returns a list of all admins
func (mgr *Manager) AdminGetList() (result []types.Admin, err error) {
	err = mgr.admins.Find(bson.M{}).All(&result)
	return
}
