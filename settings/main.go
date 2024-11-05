package settings

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Tom5521/fsize/locales"
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

var (
	Settings *conf.Preferences
	po       = locales.Po
)

func Load() (err error) {
	Settings, err = conf.New("fsize")
	for _, k := range Keys {
		Settings.CreateNewFields(false, conf.Field{Key: k, FieldType: reflect.TypeFor[bool]()})
	}
	return
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
		var exists bool
		for _, k := range Keys {
			if k == key {
				exists = true
				break
			}
		}
		if !exists {
			return fmt.Errorf(po.Get("unrecognized key \"%s\"", key))
		}
		Settings.SetBool(key, v)
	}
	return nil
}
