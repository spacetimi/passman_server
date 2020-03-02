package testing

import (
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_wrapper"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"reflect"
)

/**
 * Implements IDataItemCollection
 */
type StarCollection struct {
	mongo_wrapper.DataItemList
	mongo_wrapper.Dirtyable
}

// Returns a copy each time it is called,
// do not use lightly for large collections.
func (sc *StarCollection) Stars() []*Star {
	starsCopy := make([]*Star, 0)
	for _, item := range sc.GetDataItems() {
		star, castOk := item.(*Star)
		if castOk {
			starsCopy = append(starsCopy, star)
		} else {
			logger.LogWarning("Object with incorrect data-type in collection: " +  reflect.TypeOf(item).String())
		}
	}
	return starsCopy
}

func (sc *StarCollection) GetDescriptor() mongo_wrapper.IDataItemDescriptor {
	return StarDataDescriptor
}

func (sc *StarCollection) GetDataItemFactory() mongo_wrapper.IDataItemFactory {
	return starFactory
}

/**
 * Implements IDataItemFactory
 */
type StarFactory struct {
}
var starFactory StarFactory

func (sf StarFactory) CreateDataItem() mongo_wrapper.IDataItem {
	return &Star{}
}

