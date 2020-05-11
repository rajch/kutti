package configfilemanager

// Configdata provides methods for serializing and deserializing
// config data, and setting default values. These methods will be called
// by a ConfigManager as appropriate.
type Configdata interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	Setdefaults()
}

// ConfigManager saves and loads configuration data from a file in the kutti
// config directory
type ConfigManager interface {
	Load() error
	Save() error
}

type fileConfigManager struct {
	configfilename string
	configdata     Configdata
}

// Load loads a saved config, or initializes default values
func (cm *fileConfigManager) Load() error {
	data, notexist, err := LoadConfigfile(cm.configfilename)
	if notexist {
		cm.configdata.Setdefaults()
		return cm.Save()
	}

	if err != nil {
		return err
	}

	err = cm.configdata.Deserialize(data)
	if err != nil {
		cm.configdata.Setdefaults()
		cm.Save()
		return err
	}

	return nil
}

// Save saves a config
func (cm *fileConfigManager) Save() error {
	data, err := cm.configdata.Serialize()
	if err != nil {
		return err
	}

	return SaveConfigfile(cm.configfilename, data)
}

// Reset resets a config to default values
func (cm *fileConfigManager) Reset() {
	cm.configdata.Setdefaults()
}

// New returns a new Configmanager
func New(filename string, s Configdata) ConfigManager {
	if filename == "" || s == nil {
		panic("Must provide configuration file name and serializer.")
	}
	result := &fileConfigManager{
		configfilename: filename,
		configdata:     s,
	}
	err := result.Load()
	if err != nil {
		panic(err)
	}

	return result
}
