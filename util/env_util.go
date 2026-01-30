package util

import (
	"log"
	"os"
)

type Config interface {
	GetConfig(configName string) string
}

type ConfigImpl struct {
	config map[string]string
}

var EnvConfig *ConfigImpl

func init() {

	var initConfig = make(map[string]string)

	// grpc server port
	GrpcServerPort, ok := os.LookupEnv("GrpcServerPort")

	if !ok {
		GrpcServerPort = "8081"
		log.Println("GRPC Server port not provided in environment. Defaulting to ", GrpcServerPort)
	}

	initConfig["GrpcServerPort"] = GrpcServerPort

	// gateway server port
	GatewayServerPort, ok := os.LookupEnv("GatewayServerPort")

	if !ok {
		GatewayServerPort = "8080"
		log.Println("Gateway Server port not provided in environment. Defaulting to ", GatewayServerPort)
	}

	initConfig["GatewayServerPort"] = GatewayServerPort

	// client http settings
	HttpClientMaxIdleCon, ok := os.LookupEnv("HttpClientMaxIdleCon")

	if !ok {
		HttpClientMaxIdleCon = "10"
		log.Println("HttpClientMaxIdleCon not set in environment. Defaulting to ", HttpClientMaxIdleCon)

	}

	initConfig["HttpClientMaxIdleCon"] = HttpClientMaxIdleCon

	HttpClientIdleConTimeout, ok := os.LookupEnv("HttpClientIdleConTimeout")

	if !ok {
		HttpClientIdleConTimeout = "30"
		log.Println("HttpClientIdleConTimeout not set in environment. Defaulting to ", HttpClientIdleConTimeout)

	}

	initConfig["HttpClientIdleConTimeout"] = HttpClientIdleConTimeout

	HttpClientDisableComp, ok := os.LookupEnv("HttpClientDisableComp")

	if !ok {
		HttpClientDisableComp = "false"
		log.Println("HttpClientDisableComp not set in environment. Defaulting to ", HttpClientDisableComp)

	}

	initConfig["HttpClientDisableComp"] = HttpClientDisableComp

	// START AUTHENTICATION DB INFO
	BptResourceNewsMaxLifetime, ok := os.LookupEnv("BptResourceNewsRepoMaxLifetime")

	if !ok {
		BptResourceNewsMaxLifetime = "30m"
		log.Println("BptResourceNewsMaxLifetime not set in environment. Defaulting to ", BptResourceNewsMaxLifetime)

	}

	initConfig["BptResourceNewsMaxLifetime"] = BptResourceNewsMaxLifetime

	BptResourceNewsMaxAlive, ok := os.LookupEnv("BptResourceNewsRepoMaxAlive")

	if !ok {
		BptResourceNewsMaxAlive = "400"
		log.Println("BptResourceNewsMaxAlive not set in environment. Defaulting to ", BptResourceNewsMaxAlive)

	}

	initConfig["BptResourceNewsMaxAlive"] = BptResourceNewsMaxAlive

	BptResourceNewsMaxIdle, ok := os.LookupEnv("BptResourceNewsRepoMaxIdle")

	if !ok {
		BptResourceNewsMaxIdle = "5"
		log.Println("BptResourceNewsMaxIdle not set in environment. Defaulting to ", BptResourceNewsMaxIdle)

	}

	initConfig["BptResourceNewsMaxIdle"] = BptResourceNewsMaxIdle

	BptResourceNewsAddress, ok := os.LookupEnv("BptResourceNewsRepoAddress")

	if !ok {
		BptResourceNewsAddress = "127.0.0.1"
		log.Println("BptResourceNewsAddress not set in environment. Defaulting to ", BptResourceNewsAddress)
	}

	initConfig["BptResourceNewsAddress"] = BptResourceNewsAddress

	BptResourceNewsPort, ok := os.LookupEnv("BptResourceNewsRepoPort")

	if !ok {
		BptResourceNewsPort = "7856"
		log.Println("BptResourceNewsPort not set in environment. Defaulting to ", BptResourceNewsPort)
	}

	initConfig["BptResourceNewsPort"] = BptResourceNewsPort

	BptResourceNewsUser, ok := os.LookupEnv("BptResourceNewsRepoUser")

	if !ok {
		BptResourceNewsUser = "aaaaaaaaaa"
		log.Println("BptResourceNewsUser not set in environment. Defaulting to ", BptResourceNewsUser)
	}

	initConfig["BptResourceNewsUser"] = BptResourceNewsUser
	BptResourceNewsPassword, ok := os.LookupEnv("BptResourceNewsRepoPassword")

	if !ok {
		BptResourceNewsPassword = "111111111111"
		log.Println("BptResourceNewsPassword not set in environment. Defaulting to ", BptResourceNewsPassword)

	}

	initConfig["BptResourceNewsPassword"] = BptResourceNewsPassword

	BptResourceNewsDBName, ok := os.LookupEnv("BptResourceNewsRepoDBName")

	if !ok {
		BptResourceNewsDBName = "BPTRESOURCENEWS"
		log.Println("BptResourceNewsDBName not set in environment. Defaulting to ", BptResourceNewsDBName)

	}

	initConfig["BptResourceNewsDBName"] = BptResourceNewsDBName
	// END BPT RESOURCE NEWS DB INFO

	ForwardForHeader, ok := os.LookupEnv("ForwardForHeader")

	if !ok {
		ForwardForHeader = "grpcgateway-x-forwarded-for"
		log.Println("ForwardForHeader not set in environment. Defaulting to ", ForwardForHeader)

	}

	initConfig["ForwardForHeader"] = ForwardForHeader

	EnvConfig = &ConfigImpl{initConfig}
}

func (c *ConfigImpl) GetConfig(configName string) string {

	strConfig, ok := c.config[configName]

	if ok {
		return strConfig
	} else {
		log.Println("Config " + configName + "does not exist.")
		return ""
	}

}
