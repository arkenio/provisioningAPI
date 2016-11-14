package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	router.PathPrefix("/doc").Handler(http.FileServer(FS(false)))
	router.PathPrefix("/swagger.yaml").HandlerFunc(serveSwaggerYaml)
	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! go to /doc")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"ProvisionS3Post",
		"POST",
		"/provision/s3",
		ProvisionS3Post,
	},

	Route{
		"ProvisionAtlasPost",
		"POST",
		"/provision/atlas",
		ProvisionAtlasPost,
	},
	
	Route{
		"ProvisionAtlasGetCluster",
		"GET",
		"/provision/atlas/{clusterName}",
		ProvisionAtlasGetCluster,
	},

	Route{
		"ProvisionersGet",
		"GET",
		"/provisioners",
		ProvisionersGet,
	},
}

func serveSwaggerYaml(w http.ResponseWriter, r *http.Request) {
	type TemplateVars struct {
		Host string
	}
	swaggerTpl := FSMustString(false, "/swagger.tpl")
	t := template.Must(template.New("swagger").Parse(swaggerTpl))
	t.Execute(w, &TemplateVars{r.RemoteAddr})
}
