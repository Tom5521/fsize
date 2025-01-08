package settings

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Tom5521/fsize/locales"
	"github.com/spf13/viper"
)

const (
	AlwaysPrintOnWalk  = "Always-Print-On-Walk"
	AlwaysSkipWalk     = "Always-Skip-Walk"
	AlwaysShowProgress = "Always-Show-Progress"
	HideWarnings       = "Hide-Warnings"
)

var po = locales.Po

func Load() error {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("error getting user config path: %v", err)
	}
	configPath += "/fsize"

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	viper.SetDefault(AlwaysPrintOnWalk, false)
	viper.SetDefault(AlwaysSkipWalk, false)
	viper.SetDefault(AlwaysShowProgress, true)
	viper.SetDefault(HideWarnings, false)

read:
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err = os.MkdirAll(configPath, os.ModePerm); err != nil {
				return fmt.Errorf("error creating configuration directory: %v", err)
			}
			if err = viper.SafeWriteConfigAs(configPath + "/config.json"); err != nil {
				return fmt.Errorf("error writing to the default configuration file: %v", err)
			}
			goto read
		} else {
			return fmt.Errorf("error reading configuration: %v", err)
		}
	}

	return nil
}

func Parse(optionsArgs []string) error {
	for _, option := range optionsArgs {
		data := strings.SplitN(option, "=", 2)
		if len(data) != 2 {
			return errors.New(po.Get("syntax error"))
		}
		key, value := data[0], data[1]

		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf(po.Get("unrecognized value type \"%s\"", value))
		}
		if !viper.IsSet(key) {
			return fmt.Errorf(po.Get("unrecognized key \"%s\"", key))
		}
		viper.Set(key, v)
	}
	return nil
}
