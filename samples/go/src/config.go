package src

import (
	"fmt"
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	SEED string
	HOST string
	PORT string
}

func GetConfig(params ...string) Configuration {
	configuration := Configuration{}
	env := "dev"
	if len(params) > 0 {
		env = params[0]
	}
	fileName := fmt.Sprintf("./src/%s_config.json", env)
	gonfig.GetConf(fileName, &configuration)
	return configuration
}

