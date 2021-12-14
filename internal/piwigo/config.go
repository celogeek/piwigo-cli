package piwigo

import (
	"encoding/json"
	"fmt"
	"os"
)

func (p *Piwigo) ConfigPath() (configPath string, err error) {
	configDir, err := os.UserConfigDir()
	if err == nil {
		configPath = fmt.Sprintf("%s/piwigo-cli", configDir)
	}
	return
}

func (p *Piwigo) CreateConfigDir() (configPath string, err error) {
	configPath, err = p.ConfigPath()
	if err != nil {
		return
	}

	err = os.MkdirAll(configPath, os.FileMode(0700))
	return
}

func (p *Piwigo) SaveConfig() (err error) {
	configPath, err := p.CreateConfigDir()
	if err != nil {
		return
	}

	configFile := fmt.Sprintf("%s/config.json", configPath)

	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, b, os.FileMode(0700))
	return err
}

func (p *Piwigo) LoadConfig() (err error) {
	configPath, err := p.ConfigPath()
	if err != nil {
		return
	}

	configFile := fmt.Sprintf("%s/config.json", configPath)
	b, err := os.ReadFile(configFile)
	if os.IsNotExist(err) {
		err = nil
		return
	}
	if err != nil {
		err = json.Unmarshal(b, p)
	}
	return
}
