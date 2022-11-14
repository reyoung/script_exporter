package main

import "time"

type CfgV1 struct {
	Metrics   map[string]CfgV1Metric `yaml:"metrics"`
	Namespace string                 `yaml:"namespace"`
}

func (c *CfgV1) isVersionedConfig() {}

type CfgV1Metric struct {
	Kind     string              `yaml:"kind"`
	Command  string              `yaml:"command"`
	Matrix   map[string][]string `yaml:"matrix"`
	Interval time.Duration       `yaml:"interval"`
	Help     string              `yaml:"help"`
}

func init() {
	gVersionedConfigs["v1"] = func() VersionedConfig {
		return &CfgV1{}
	}
}
