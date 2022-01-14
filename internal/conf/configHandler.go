package conf

import (
	"github.com/spf13/viper"
)

const (
	//Mock files directory configuration key name
	MOCKS_DIR = "mocks-dir"
)

//Initiate and return a ref to Alfred configuration.
func InitConfig() (viper.Viper, error) {

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("configs")

	err := viper.ReadInConfig() // Find and read the config file

	return *viper.GetViper(), err
}
