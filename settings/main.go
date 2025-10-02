package settings

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/xdg"
	po "github.com/leonelquinteros/gotext"
	"github.com/spf13/viper"
)

const (
	AlwaysPrintOnWalk  = "Always-Print-On-Walk"
	AlwaysSkipWalk     = "Always-Skip-Walk"
	AlwaysShowProgress = "Always-Show-Progress"
	HideWarnings       = "Hide-Warnings"
	HideProgress       = "Hide-Progress"
	NoColor            = "No-Color"
	Language           = "Language"
	ProgressDelay      = "Progress-Delay"
)

func InitSettings() error {
	configPath := filepath.Join(xdg.ConfigHome, "fsize")

	viper.SetConfigName("fsize")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)

	viper.SetDefault(AlwaysPrintOnWalk, false)
	viper.SetDefault(AlwaysSkipWalk, false)
	viper.SetDefault(AlwaysShowProgress, true)
	viper.SetDefault(HideWarnings, false)
	viper.SetDefault(Language, "default")
	viper.SetDefault(HideProgress, false)
	viper.SetDefault(NoColor, false)
	viper.SetDefault(ProgressDelay, "1s")

read:
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err = os.MkdirAll(configPath, os.ModePerm); err != nil {
				return fmt.Errorf("error creating configuration directory: %w", err)
			}
			if err = viper.SafeWriteConfigAs(filepath.Join(configPath, "fsize.json")); err != nil {
				return fmt.Errorf("error writing to the default configuration file: %w", err)
			}
			goto read
		} else {
			return fmt.Errorf("error reading configuration: %w", err)
		}
	}

	return nil
}

func Parse(optionsArgs []string) error {
	for _, option := range optionsArgs {
		data := strings.Split(option, "=")
		if len(data) != 2 {
			return errors.New(po.Get("syntax error"))
		}
		key, value := data[0], data[1]

		if !viper.IsSet(key) {
			return errors.New(po.Get("unrecognized key \"%s\"", key))
		}

		var (
			newValue any
			err      error
		)

		switch viper.Get(key).(type) {
		case string:
			if key == ProgressDelay {
				_, err = time.ParseDuration(value)
				if err != nil {
					return errors.New(po.Get("invalid time duration(%s): %s", value, err.Error()))
				}
			}
			newValue = value
		case bool:
			newValue, err = strconv.ParseBool(value)
			if err != nil {
				return errors.New(po.Get("unrecognized value type \"%s\"", value))
			}
		}

		viper.Set(key, newValue)
	}
	return nil
}
