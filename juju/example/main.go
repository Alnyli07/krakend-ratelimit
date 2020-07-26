package main

import (
	"flag"
	"log"
	"os"

	jujuproxy "github.com/alnyli07/krakend-ratelimit/juju/proxy"
	jujurouter "github.com/alnyli07/krakend-ratelimit/juju/router/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	krakendgin "github.com/devopsfaith/krakend/router/gin"
	"github.com/devopsfaith/krakend/transport/http/client"
	http "github.com/devopsfaith/krakend/transport/http/server"
	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.Int("p", 0, "Port of the service")
	logLevel := flag.String("l", "ERROR", "Logging level")
	debug := flag.Bool("d", false, "Enable the debug")
	configFile := flag.String("c", "/etc/krakend/configuration.json", "Path to the configuration filename")
	flag.Parse()

	parser := config.NewParser()
	serviceConfig, err := parser.Parse(*configFile)
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}
	serviceConfig.Debug = serviceConfig.Debug || *debug
	if *port != 0 {
		serviceConfig.Port = *port
	}

	logger, err := logging.NewLogger(*logLevel, os.Stdout, "[KRAKEND]")
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}

	routerFactory := krakendgin.NewFactory(krakendgin.Config{
		Engine:         gin.Default(),
		ProxyFactory:   proxy.NewDefaultFactory(jujuproxy.BackendFactory(proxy.CustomHTTPProxyFactory(client.NewHTTPClient)), logger),
		Middlewares:    []gin.HandlerFunc{},
		Logger:         logger,
		HandlerFactory: jujurouter.HandlerFactory,
		RunServer:      http.RunServer,
	})

	routerFactory.New().Run(serviceConfig)
}
