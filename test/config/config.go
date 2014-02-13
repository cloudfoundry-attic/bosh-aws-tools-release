package config

import (
	"encoding/json"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	DirectorUUID string `json:"director_uuid"`
	DirectorHost string `json:"director_host"`

	AwsAccessId         string `json:"aws_access_id"`
	AwsSecretAcccessKey string `json:"aws_secret_access_key"`

	Route53ZoneNames []string `json:"route53_zone_names"`
}

func LoadAndValidate() Config {
	path := os.Getenv("CONFIG")
	if path == "" {
		panic("Must set $CONFIG to point to an integration config .json file.")
	}

	config := loadPath(path)

	inferDirectorAttributes(&config)

	config.validate()

	return config
}

func loadPath(path string) Config {
	configFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	config := Config{}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func (c Config) validate() {
	if c.DirectorHost == "" {
		panic("Must set director_host")
	}
	if c.DirectorUUID == "" {
		panic("Must set director_uuid")
	}

	if c.AwsAccessId == "" {
		panic("Must set aws_access_id")
	}
	if c.AwsSecretAcccessKey == "" {
		panic("Must set aws_secret_access_key")
	}

	if len(c.Route53ZoneNames) == 0 {
		panic("Must set route53_zone_names")
	}
}

func inferDirectorAttributes(config *Config) {
	if config.DirectorUUID == "" {
		config.DirectorUUID = inferDirectorUUID()
	}
	if config.DirectorHost == "" {
		config.DirectorHost = inferDirectorHost()
	}
}

// assume that director is pre-targeted
func inferDirectorUUID() string {
	output, err := exec.Command("bash", "-c", "bosh status | grep UUID | cut -d' ' -f 10").Output()
	if err != nil {
		panic(err)
	}

	return strings.Trim(string(output), "\n\r\t ")
}

func inferDirectorHost() string {
	output, err := exec.Command("bash", "-c", "bosh status | grep URL | cut -d' ' -f 11").Output()
	if err != nil {
		panic(err)
	}

	url, err := url.Parse(string(output))
	if err != nil {
		panic(err)
	}

	return strings.Split(url.Host, ":")[0]
}
