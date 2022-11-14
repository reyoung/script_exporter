package main

import "emperror.dev/errors"

type VersionedConfig interface {
	isVersionedConfig()
}

type Config struct {
	Version string
	Config  VersionedConfig
}

var gVersionedConfigs = map[string]func() VersionedConfig{}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Version string `yaml:"version"`
	}
	err := unmarshal(&aux)
	if err != nil {
		return err
	}
	c.Version = aux.Version
	cfgCreator, ok := gVersionedConfigs[c.Version]
	if !ok {
		return errors.Errorf("config creator %s is not registered", c.Version)
	}
	c.Config = cfgCreator()
	return unmarshal(c.Config)
}
