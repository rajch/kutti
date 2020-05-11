package defaults

import (
	"encoding/json"

	"github.com/rajch/kutti/internal/pkg/configfilemanager"
)

var (
	defaultconfigmanager configfilemanager.ConfigManager
	defaultdata          *defaultconfigdata
)

type defaultconfigdata struct {
	defaults map[string]string
}

func (dc *defaultconfigdata) Serialize() ([]byte, error) {
	return json.Marshal(dc.defaults)
}

func (dc *defaultconfigdata) Deserialize(data []byte) error {
	loadeddefaults := make(map[string]string)
	err := json.Unmarshal(data, &loadeddefaults)
	if err == nil {
		dc.defaults = loadeddefaults
	}

	return err
}

func (dc *defaultconfigdata) Setdefaults() {
	dc.defaults = map[string]string{
		"cluster": "",
		"driver":  "vbox",
		"version": "1.18",
	}
}

// Getdefault gets a default value
func Getdefault(defaultname string) string {
	return defaultdata.defaults[defaultname]
}

// Setdefault sets a default value
func Setdefault(name string, value string) {
	defaultdata.defaults[name] = value
	defaultconfigmanager.Save()
}

func init() {
	defaultdata = &defaultconfigdata{}
	defaultconfigmanager = configfilemanager.New("kuttidefaults.json", defaultdata)
}
