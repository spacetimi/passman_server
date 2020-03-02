package testing

import (
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_wrapper"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"reflect"
)

/**
 * Implements IDataItemCollection
 */
type ConstellationCollection struct {
	mongo_wrapper.DataItemList
	mongo_wrapper.Dirtyable
}

// Returns a copy each time it is called,
// do not use lightly for large collections.
func (cc *ConstellationCollection) Constellations() []*Constellation {
	constellationsCopy := make([]*Constellation, 0)
	for _, item := range cc.GetDataItems() {
		constellation, castOk := item.(*Constellation)
		if castOk {
			constellationsCopy = append(constellationsCopy, constellation)
		} else {
			logger.LogWarning("Object with incorrect data-type in collection: " +  reflect.TypeOf(item).String())
		}
	}
	return constellationsCopy
}

func (cc *ConstellationCollection) GetDescriptor() mongo_wrapper.IDataItemDescriptor {
	return ConstellationDataDescriptor
}

func (cc *ConstellationCollection) GetDataItemFactory() mongo_wrapper.IDataItemFactory {
	return constellationFactory
}

/**
 * Implements IDataItemFactory
 */
type ConstellationFactory struct {
}
var constellationFactory ConstellationFactory

func (cf ConstellationFactory) CreateDataItem() mongo_wrapper.IDataItem {
	return &Constellation{}
}

