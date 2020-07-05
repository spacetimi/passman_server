package app_init

import (
	"github.com/spacetimi/passman_server/app_src/app_utils/app_emailer"
	"github.com/spacetimi/passman_server/app_src/metadata/faq"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_factory"
	"github.com/spacetimi/timi_shared_server/code/core/shared_init"
)

func GetAppInitializer() shared_init.IAppInitializer {
	return &appInitializer
}

type AppInitializer struct { // Implements IAppInit
}

var appInitializer AppInitializer

/********** Begin IAppInitializer implementation **********/
func (appInitializer *AppInitializer) AppName() string {
	return "passman_server"
}

func (appInitializer *AppInitializer) AppInit() error {

	registerMetadataFactories()

	app_emailer.Initialize()

	return nil
}

/********** End IAppInitializer implementation **********/

func registerMetadataFactories() {
	metadata_factory.RegisterFactory(faq.MetadataKey, faq.MetadataFactory{})
}
