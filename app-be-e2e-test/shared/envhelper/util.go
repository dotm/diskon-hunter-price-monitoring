package envhelper

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func GetEnvVar(key string) string {
	envFromLambda := os.Getenv(key)
	if envFromLambda != "" {
		return envFromLambda
	}
	return viper.GetString(key)
}

// for use in playground only
func SetLocalEnvVar() (err error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config file: %v\n", err)
		os.Exit(1)
	}

	return nil
}
