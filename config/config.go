package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/iij/legs-client/util"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config wraps viper.Viper
type Config struct{ *viper.Viper }

var home = setHomeDir()

func setHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return usr.HomeDir
}

// InitConfig load config file.
// It try loading config in this order:
//   * path (specify `-c` option)
//	 * $XDG_CONFIG_HOME/legsc/conf.toml
//   * xdg_config_dir/legsc/conf.toml (where xdg_config_dir is in $XDG_CONFIG_DIRS)
//   * ~/.config/legsc/conf.toml
//   * ~/.legsc/conf.toml
func InitConfig(path string) (*Config, error) {
	config := viper.New()
	if path != "" {
		config.SetConfigFile(path)
	} else {
		config.SetConfigName("conf")
		if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
			config.AddConfigPath(filepath.Join(dir, "legsc"))
		}
		if dirs := os.Getenv("XDG_CONFIG_DIRS"); dirs != "" {
			for _, dir := range strings.Split(dirs, fmt.Sprintf("%c", filepath.ListSeparator)) {
				config.AddConfigPath(filepath.Join(dir, "legsc"))
			}
		}
		config.AddConfigPath(filepath.Join(home, ".config", "legsc"))
		config.AddConfigPath(filepath.Join(home, ".legsc"))
	}

	config.SetEnvPrefix("legs")
	config.AutomaticEnv()

	setDefaultConfig(config)

	if err := config.ReadInConfig(); err != nil {
		switch err := errors.Cause(err).(type) {
		case viper.ConfigFileNotFoundError:
			config.SetConfigFile(filepath.Join(home, ".config", "legsc", "conf.toml"))
			return &Config{config}, nil
		default:
			return nil, err
		}
	}

	return &Config{config}, nil
}

// Save set a parameter of config data by key and value.
// And this method save configuration to file.
func (c *Config) Save(key string, value interface{}) error {
	c.Set(key, value)

	if err := c.Write(); err != nil {
		return err
	}

	fmt.Println("save to", c.ConfigFileUsed())
	return nil
}

// Write save config to file, and change permission.
func (c *Config) Write() (err error) {
	err = c.WriteConfig()
	if err != nil {
		return
	}

	err = os.Chmod(c.ConfigFileUsed(), 0600)
	if err != nil {
		return
	}

	return nil
}

func setDefaultConfig(config *viper.Viper) {
	macAddrs := util.GetMacAddr()
	var mac string
	if len(macAddrs) == 0 {
		mac = ""
	} else {
		mac = macAddrs[0]
	}

	tmpDir := filepath.Join("/", "var", "tmp", "legsc")

	config.SetDefault("server", "legs-api.pms.iij.jp")
	config.SetDefault("proxy", "")
	config.SetDefault("secret", "<please set to your secret key>")
	config.SetDefault("device_name", mac)
	config.SetDefault("pid_file", filepath.Join(tmpDir, "legsc.pid"))
	config.SetDefault("sock_file", filepath.Join(tmpDir, "legsc.sock"))
	config.SetDefault("status_file", filepath.Join(tmpDir, "legsc.status"))
	config.SetDefault("use_wss", true)
}
