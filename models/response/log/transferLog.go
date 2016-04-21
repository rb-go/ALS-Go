package log

import "gopkg.in/mgo.v2/bson"

type TransferLog struct {
	TransferedLogId []bson.ObjectId `json`
}
