package main

import (
	"github.com/ghodss/yaml"
	"github.com/harsha-aqfer/todo/internal/service_echo"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing parameter, provide file name!")

	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	conf := service_echo.NewConfig()
	if err = yaml.Unmarshal(data, conf); err != nil {
		log.Fatal(err)
	}

	app, err := service_echo.NewService(conf)
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}
