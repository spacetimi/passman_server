package app_init

import (
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_factory"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
	"github.com/spacetimi/timi_shared_server/code/core/shared_init"
)

func GetAppInitializer() shared_init.IAppInitializer {
	return &appInitializer
}

type AppInitializer struct {	// Implements IAppInit
}
var appInitializer AppInitializer

/********** Begin IAppInitializer implementation **********/
func (appInitializer *AppInitializer) AppInit() {
	registerMetadataFactories()
}
/********** End IAppInitializer implementation **********/

// TODO: Avi: Move this somewhere else?
func registerMetadataFactories() {
	metadata_factory.RegisterFactory("MetadataTest", MetadataTestFactory{})
	metadata_factory.RegisterFactory("Planets", PlanetsFactory{})
	metadata_factory.RegisterFactory("Places", PlacesFactory{})
}


type MetadataTestFactory struct {	// Implements IMetadataFactory
}

func (f MetadataTestFactory) Instantiate() metadata_typedefs.IMetadataItem {
	return &MetadataTest{}
}

// TODO: Avi: Move this somewhere else
type MetadataTest struct {
	Id int
}
func (m MetadataTest) GetKey() string {
	return "MetadataTest"
}
func (m MetadataTest) GetMetadataSpace() metadata_typedefs.MetadataSpace {
	return metadata_typedefs.METADATA_SPACE_APP
}





type PlanetsFactory struct {	// Implements IMetadataFactory
}

func (f PlanetsFactory) Instantiate() metadata_typedefs.IMetadataItem {
	return &Planets{}
}

// TODO: Avi: Move this somewhere else
type Planets struct {
	Planets []string
}
func (m Planets) GetKey() string {
	return "Planets"
}
func (m Planets) GetMetadataSpace() metadata_typedefs.MetadataSpace {
	return metadata_typedefs.METADATA_SPACE_APP
}



type PlacesFactory struct {	// Implements IMetadataFactory
}

func (f PlacesFactory) Instantiate() metadata_typedefs.IMetadataItem {
	return &Places{}
}

// TODO: Avi: Move this somewhere else
type Places struct {
	Places []string
}
func (m Places) GetKey() string {
	return "Places"
}
func (m Places) GetMetadataSpace() metadata_typedefs.MetadataSpace {
	return metadata_typedefs.METADATA_SPACE_APP
}


