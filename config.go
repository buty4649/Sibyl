package main

import (
    "gopkg.in/yaml.v2"
)

type ConfigParameter struct {
    CmdLine string `yaml:commandline`
    Exec    string `yaml:exec`
}
type Config struct {
    Fork     []ConfigParameter
    Exec     []ConfigParameter
    Coredump []ConfigParameter
    Exit     []ConfigParameter
}

func LoadConfig(r []byte) (*Config, error) {
    config := Config{}

    err := yaml.Unmarshal(r, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}
