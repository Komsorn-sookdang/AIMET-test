package main

import (
	"aimet-test/configs"
	"aimet-test/internal/databases"
	"aimet-test/internal/router"
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	configs.InitConfig()

	databases.InitMongoClient()
	defer databases.CloseMongoClient()

	databases.InitRedisClient()
	defer databases.CloseRedisClient()

	r := router.SetupRouter()

	r.Run(fmt.Sprintf(":%d", viper.GetInt("port")))
}
