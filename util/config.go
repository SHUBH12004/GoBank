package util

import "github.com/spf13/viper"

// Config stores all configuration of the application
// the values are read by viper from config file or environment variables

// mapstructure is used to map the environment variables to the struct fields
// it is used by viper to decode the configuration file into the struct
type Config struct {
	// DBDriver is the database driver to use
	DBDriver string `mapstructure:"DB_DRIVER"`
	// DBSource is the database source to use
	DBSource string `mapstructure:"DB_SOURCE"`
	// ServerAddress is the address to listen on for HTTP requests
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	// TokenSymmetricKey is the symmetric key used to sign tokens
}

// LoadConfig reads configuration from a file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)  // Set the path to look for the config file
	viper.SetConfigName("app") // Set the name of the config file (without extension)
	viper.SetConfigType("env") // Set the type of the config file (e.g., env, json, yaml)

	viper.AutomaticEnv()       // Automatically read environment variables that match the struct fields
	err = viper.ReadInConfig() // Read the config file
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config) // Unmarshal the config file into the Config struct
	return

}
