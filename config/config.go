package config

import "github.com/numbatx/numbat-proxy/data"

// GeneralSettingsConfig will hold the general settings for a node
type GeneralSettingsConfig struct {
	ServerPort          int
	CfgFileReadInterval int
}

// Config will hold the whole config file's data
type Config struct {
	GeneralSettings GeneralSettingsConfig
	Observers       []*data.Observer
}
