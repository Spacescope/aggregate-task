package utils

import (
	"syscall"

	"github.com/jinzhu/configor"
)

// var CNF TomlConfig

type TomlConfig struct {
	Title         string
	AggregateTask AggregateTask `toml:"aggregate-task"`
	Task          Task          `toml:"task"`
}

type AggregateTask struct {
	Addr          string `toml:"listen"`
	ObservatoryDB string `toml:"observatorydb"`
	NotifyDB      string `toml:"notifydb"`
}

type Task struct {
	Name   string   `toml:"name"`
	Depend []string `toml:"depend"`
}

func InitConfFile(file string, cf *TomlConfig) error {
	err := syscall.Access(file, syscall.O_RDONLY)
	if err != nil {
		return err
	}
	err = configor.Load(cf, file)
	if err != nil {
		return err
	}

	return nil
}
