package dummy_controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func DummyController(httpResponseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	switch params["param1"] {
	case "home":
		dummy_home(httpResponseWriter, params)
	default:
		fmt.Fprintln(httpResponseWriter, "No matching method found")
	}
}

func dummy_home(httpResponseWriter http.ResponseWriter, params map[string]string) {
	fmt.Fprintln(httpResponseWriter, "Dummy Home")
}
