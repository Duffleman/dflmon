package app

import (
	"net/http"

	"dflmon/cachet"
	"dflmon/config"
)

// App is a struct for the app methods to attach to
type App struct {
	Client           *http.Client
	ClientNoValidate *http.Client
	Config           *config.Config
	Cachet           *cachet.Client
}
