package testing

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_wrapper"
	"log"
	"net/http"
	"strconv"
)

func TestingController(httpResponseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	switch params["param1"] {
	case "home": testing_home(httpResponseWriter, params)
	case "readstars": testing_readstars(httpResponseWriter, params)
	case "readstarbyname": testing_readstarbyfilter(httpResponseWriter, params)
	case "insertupdatestar": testing_insertupdatestar(httpResponseWriter, params)
	case "readcons": testing_readconstellations(httpResponseWriter, params)
	case "insertupdatecons": testing_insertupdateconstellation(httpResponseWriter, params)
	default:
		fmt.Fprintln(httpResponseWriter, "No matching method found")
		fmt.Fprintln(httpResponseWriter, "Method: " + params["param1"])
	}

}

func testing_home(httpResponseWriter http.ResponseWriter, params map[string]string) {
	fmt.Fprintln(httpResponseWriter, "App Name: " + config.GetAppName())
}

func testing_readstars(httpResponseWriter http.ResponseWriter, params map[string]string) {
	starIdString, ok := params["param2"]

	if !ok || starIdString == "" {
		// If no star id was passed, print all stars

		var star_collection StarCollection
		mongo_wrapper.GetDataItemCollection(&star_collection)
		for _, star := range star_collection.Stars() {
			fmt.Fprintln(httpResponseWriter, star)
		}
		return
	} else {
		// If specific star id was passed, print that

		starId, _ := strconv.Atoi(starIdString)

		var star Star
		mongo_wrapper.GetDataItemByPrimaryKeys([]interface{}{starId}, &star)

		_, err := fmt.Fprintln(httpResponseWriter, star)
		if err != nil {
			log.Fatal("Something went wrong: " + err.Error())
		}
	}
}

func testing_insertupdatestar(httpResponseWriter http.ResponseWriter, params map[string]string) {
	starName, _ := params["param3"]
	starIdString, ok := params["param2"]
	if !ok {
		fmt.Fprintln(httpResponseWriter, "Please specify star id")
		return
	}
	starId, err := strconv.ParseInt(starIdString, 10, 32)
	if err != nil {
		log.Fatal("Error parsing star id: " + err.Error())
	}

	var star_collection StarCollection
	mongo_wrapper.GetDataItemCollection(&star_collection)

	isNewStar := true
	for _, star := range star_collection.Stars() {
		if star.Id == int32(starId) {
			fmt.Fprintln(httpResponseWriter, "Found matching star: " + star.String())
			isNewStar = false
			star.Name = starName
			star.SetDirty(true)
		}
	}

	if isNewStar {
		fmt.Fprintln(httpResponseWriter, "Inserting new star")
		newStar := Star{Name:starName, Id:int32(starId)}
		newStar.SetDirty(true)
		star_collection.AddDataItem(&newStar)
	}

	star_collection.SetDirty(true)
	mongo_wrapper.ApplyDataItemCollection(&star_collection)
}

func testing_readconstellations(httpResponseWriter http.ResponseWriter, params map[string]string) {
	var cc ConstellationCollection
	mongo_wrapper.GetDataItemCollection(&cc)

	for _, constellation := range cc.Constellations() {
		fmt.Fprintln(httpResponseWriter, constellation)
	}
}

func testing_insertupdateconstellation(httpResponseWriter http.ResponseWriter, params map[string]string) {
	consName, ok := params["param2"]
	if !ok {
		fmt.Fprintln(httpResponseWriter, "Please specify constellation name")
		return
	}
	starIdString, ok := params["param3"]
	if !ok {
		fmt.Fprintln(httpResponseWriter, "Please specify star id to add")
		return
	}
	starId, err := strconv.ParseInt(starIdString, 10, 32)
	if err != nil {
		log.Fatal("Error parsing star id: " + err.Error())
	}

	var cc ConstellationCollection
	mongo_wrapper.GetDataItemCollection(&cc)

	var sc StarCollection
	mongo_wrapper.GetDataItemCollection(&sc)

	var foundConstellation *Constellation = nil
	for _, constellation := range cc.Constellations() {
		if consName == constellation.Name {
			foundConstellation = constellation
			break
		}
	}

	if foundConstellation == nil {
		fmt.Fprintln(httpResponseWriter, "No such constellation: " + consName)
		return
	}

	for _, star := range sc.Stars() {
		if int32(starId) == star.Id {
			newStar := *star
			newStar.SetDirty(true)
			foundConstellation.Stars = append(foundConstellation.Stars, newStar)
			foundConstellation.SetDirty(true)
			break
		}
	}

	mongo_wrapper.ApplyDataItem(foundConstellation)
}

func testing_readstarbyfilter(httpResponseWriter http.ResponseWriter, params map[string]string) {
	starName, ok := params["param2"]
	if !ok {
		fmt.Fprintln(httpResponseWriter, "Please specify star name")
		return
	}

	var sc StarCollection
	mongo_wrapper.GetDataItemCollectionByFilter([]string{"name"}, []interface{}{starName}, &sc)

	for _, star := range sc.Stars() {
		fmt.Fprintln(httpResponseWriter, star)
	}
}
