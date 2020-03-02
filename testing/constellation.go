package testing

import (
	"encoding/json"
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_wrapper"
)

/**
 * Implements IDataItem
 */
type Constellation struct {
	Name string
	Stars [] Star

	mongo_wrapper.Dirtyable
}

/********** Begin IDataItem implementation **********/

var ConstellationDataDescriptor = &mongo_wrapper.DataItemDescriptor {
	DBType: mongo_wrapper.APP_DB,
	CollectionName: "constellations",
	PrimaryKeys: []string{"id"},
}

func (constellation *Constellation) GetDescriptor() mongo_wrapper.IDataItemDescriptor {
	return ConstellationDataDescriptor
}

/********** End IDataItem implementation **********/

func (constellation *Constellation) String() string {
	constellation_json, _ := json.Marshal(constellation)
	return string(constellation_json)
}

