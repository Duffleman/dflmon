package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"dflmon/app"
	"dflmon/cachet"
	"dflmon/config"

	cachetSDK "github.com/andygrunwald/cachet"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initViper() {
	viper.SetEnvPrefix("DFLMON")
	viper.AutomaticEnv()

	viper.SetDefault("CACHET_URL", "https://status.dfl.mn")
	viper.SetDefault("DEBUG", "true")
}

func main() {
	initViper()

	log.Info("app launch")

	cachetClient, err := cachetSDK.NewClient(viper.GetString("CACHET_URL"), nil)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot make cachet client: %w", err))
		return
	}

	log.Info("pinging cachet instance")

	_, _, err = cachetClient.General.Ping()
	if err != nil {
		log.Fatal(fmt.Errorf("cannot ping cachet: %w", err))
		return
	}

	cachetClient.Authentication.SetTokenAuth(viper.GetString("CACHET_TOKEN"))

	log.Info("parsing config")

	config, err := config.Parse()
	if err != nil {
		log.Fatal(fmt.Errorf("cannot load config: %w", err))
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
		},
	}

	clientNoValidate := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	app := &app.App{
		Client:           client,
		ClientNoValidate: clientNoValidate,
		Config:           config,
		Cachet:           &cachet.Client{cachetClient},
	}

	log.Info("syncing jobs with cachet components")

	err = app.SyncWithCachet()
	if err != nil {
		log.Fatal(fmt.Errorf("cannot sync with cachet: %w", err))
		return
	}

	log.Info("starting workers")

	if viper.GetString("DEBUG") == "false" {
		log.Info("setting debug OFF")

		log.SetLevel(log.WarnLevel)
	}

	log.Info("setting debug ON")

	err = app.StartWorkers()
	if err != nil {
		log.Fatal(fmt.Errorf("cannot start workers: %w", err))
		return
	}
}
