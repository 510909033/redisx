package e3w

import (
	"flag"
	"fmt"
	"github.com/510909033/redisx/e3w/internal/conf"
	"github.com/510909033/redisx/e3w/internal/e3ch"
	"github.com/510909033/redisx/e3w/internal/routers"
	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/version"
	"os"
)

const (
	PROGRAM_NAME    = "e3w"
	PROGRAM_VERSION = "0.1.0"
)

var configFilepath string

func init() {
	flag.StringVar(&configFilepath, "conf", "conf/config.default.ini", "config file path")
	rev := flag.Bool("rev", false, "print rev")
	flag.Parse()

	if *rev {
		fmt.Printf("[%s v%s]\n[etcd %s]\n",
			PROGRAM_NAME, PROGRAM_VERSION,
			version.Version,
		)
		os.Exit(0)
	}
}

func StartE3w() {
	config, err := conf.Init(configFilepath)
	if err != nil {
		panic(err)
	}

	client, err := e3ch.NewE3chClient(config)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.UseRawPath = true
	routers.InitRouters(router, config, client)
	router.Run(":" + config.Port)
}
