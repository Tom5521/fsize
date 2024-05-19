package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	conf "github.com/Tom5521/goconf"
)

const (
	AlwaysPrintOnWalk  = "AlwaysPrintOnWalk"
	AlwaysSkipWalk     = "AlwaysSkipWalk"
	AlwaysShowProgress = "AlwaysShowProgress"
	HideWarnings       = "HideWarnings"
)

var Keys = []string{
	AlwaysPrintOnWalk,
	AlwaysSkipWalk,
	AlwaysShowProgress,
	HideWarnings,
}

var Settings conf.Preferences

func LoadSettings() error {
	var err error
	Settings, err = conf.New("fsize")
	if err != nil {
		return err
	}
	return nil
}

func ParseSettings(optionsArgs []string) error {
	//options := strings.Split(optionsArgs, ";")
	for _, option := range optionsArgs {
		data := strings.SplitN(option, "=", 2)
		if len(data) != 2 {
			return errors.New("syntax error")
		}
		key, value := data[0], data[1]

		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("unrecognized value type \"%s\"", value)
		}
		var exists bool
		for _, k := range Keys {
			if k == key {
				exists = true
				break
			}
		}
		if !exists {
			return fmt.Errorf("unrecognized key \"%s\"", key)
		}
		Settings.SetBool(key, v)
	}
	return nil
}
