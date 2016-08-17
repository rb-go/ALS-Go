package log

import "gitlab.com/ergoz/ALS-Go/models/mongo"

type Get struct {
	LogList []mongo.CustomLog `json:"logList"`
}
