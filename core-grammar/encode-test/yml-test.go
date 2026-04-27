package encode_test

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Database string `yaml:"database"`
	Url      string `yaml:"url"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func YmlTest1() {
	config := Config{
		Database: "1",
		Url:      "2",
		Port:     0,
		Username: "3",
		Password: "4",
	}
	out, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(out))
}

func YmlTest2() {
	bytes, err := os.ReadFile("./encode-test/test.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	err1 := yaml.Unmarshal(bytes, &Config{})
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	fmt.Println(string(bytes))
}
