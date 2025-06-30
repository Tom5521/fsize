package stat

import (
	"fmt"
	"io/fs"
	"os/user"
	"strings"
	"text/template"
	"time"

	"github.com/Tom5521/fsize/flags"
	"github.com/gookit/color"
	"github.com/labstack/gommon/bytes"
	"github.com/leonelquinteros/gotext"
)

func _() {
	// I call these gets here so that they are not deleted in the translation files.
	gotext.Get("Name:")
	gotext.Get("Size:")
	gotext.Get("Absolute path:")
	gotext.Get("Physical path:")
	gotext.Get("Modify:")
	gotext.Get("Access:")
	gotext.Get("Birth:")
	gotext.Get("Is directory:")
	gotext.Get("Permissions:")
	gotext.Get("UID/User:")
	gotext.Get("GID/Group:")
}

const formatTemplate = `{{get "Name:"}} {{.Name}}
{{get "Size:"}} {{formatSize .Size}}
{{get "Absolute path:"}} {{.AbsPath}}
{{- if .IsLink}}
{{get "Physical path:"}} {{.PhysicalPath}}
{{- end }}
{{get "Modify:"}} {{formatTime .ModTime}}
{{get "Access:"}} {{formatTime .AccessTime}}
{{- if .SupportCreationDate }}
{{get "Birth:"}} {{formatTime .CreationTime}}
{{- end}}
{{get "Is directory:"}} {{formatBool .IsDir}}
{{get "Permissions:"}} {{formatPerms .Perms}}
{{- if .SupportFileIDs }}
{{get "UID/User:"}} {{formatUser .User}}
{{get "GID/Group:"}} {{formatGroup .Group}}
{{- end }}`

var funcMap = template.FuncMap{
	"get": func(id string) string {
		return color.Green.Render(gotext.Get(id))
	},
	"formatSize": func(size int64) string {
		return bytes.New().Format(size)
	},
	"formatTime": func(t time.Time) string {
		return t.Format(time.DateTime)
	},
	"formatBool": func(v bool) string {
		if v {
			return color.Green.Render(v)
		}
		return color.Red.Render(v)
	},
	"formatPerms": func(perms fs.FileMode) string {
		return fmt.Sprintf("%d/%s", int(perms), perms.String())
	},
	"formatUser": func(usr *user.User) string {
		return fmt.Sprintf("%s/%s", usr.Uid, usr.Username)
	},
	"formatGroup": func(grp *user.Group) string {
		return fmt.Sprintf("%s/%s", grp.Gid, grp.Name)
	},
	"walkEnabled": func() bool { return !flags.NoWalk },
}

func (f File) String() string {
	var b strings.Builder
	tmpl := template.New("format").Funcs(funcMap)
	tmpl = template.Must(tmpl.Parse(formatTemplate))

	tmpl.Execute(&b, f)

	return b.String()
}
