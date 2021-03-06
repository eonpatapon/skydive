package api

import (
	"encoding/json"
	"net/http"

	auth "github.com/abbot/go-http-auth"
	"github.com/gorilla/mux"

	"github.com/skydive-project/skydive/config"
	shttp "github.com/skydive-project/skydive/http"
	"github.com/spf13/viper"
)

type ConfigAPI struct {
	cfg *viper.Viper
}

func (c *ConfigAPI) configGet(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(&r.Request)
	// lookup into ConfigAPI
	Value := c.cfg.GetString(vars["key"])
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Value); err != nil {
		panic(err)
	}
}

func (c *ConfigAPI) registerEndpoints(r *shttp.Server) {
	routes := []shttp.Route{
		{
			Name:        "ConfigGet",
			Method:      "GET",
			Path:        "/api/config/{key}",
			HandlerFunc: c.configGet,
		},
	}

	r.RegisterRoutes(routes)
}

func RegisterConfigAPI(r *shttp.Server) {
	c := &ConfigAPI{
		cfg: config.GetConfig(),
	}

	c.registerEndpoints(r)
}
