package piwigo

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

func (p *Piwigo) ConfigPath() (configPath string, err error) {
	configDir, err := os.UserConfigDir()
	if err == nil {
		configPath = strings.Join([]string{configDir, "piwigo-cli"}, "/")
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

	configFile := strings.Join([]string{configPath, "config.json"}, "/")

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

	configFile := strings.Join([]string{configPath, "config.json"}, "/")
	b, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			err = errors.New("missing configuration file")
		}
		return
	}

	err = json.Unmarshal(b, &p)
	if p.Url == "" || p.Token == "" {
		err = errors.New("missing configuration url or token")
	}

	return
}
