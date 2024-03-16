package profile

import (
	"github.com/spf13/viper"
)

// Profile is the configuration to start main server.
type Profile struct {
	// Mode can be "prod" or "dev" or "demo"
	Mode string `json:"mode"`
	// Addr is the binding address for server
	Addr string `json:"-"`
	// Port is the binding port for server
	Port int `json:"-"`
	// Db address
	DbUri string `json:"-"`
}

// GetProfile will return a profile for dev or prod.
func GetProfile() (*Profile, error) {
	profile := Profile{}
	err := viper.Unmarshal(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
