package testing

import (
	"encoding/json"
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_wrapper"
)

/**
 * Implements IDataItem
 */
type Star struct {
	Name string
	Id int32

	mongo_wrapper.Dirtyable
}

/********** Begin IDataItem implementation **********/

var StarDataDescriptor = &mongo_wrapper.DataItemDescriptor {
	DBType: mongo_wrapper.APP_DB,
	CollectionName: "stars",
	PrimaryKeys: []string{"id"},
}

func (star *Star) GetDescriptor() mongo_wrapper.IDataItemDescriptor {
	return StarDataDescriptor
}

/********** End IDataItem implementation **********/

func (star *Star) String() string {
	star_json, _ := json.Marshal(star)
	return string(star_json)
}


