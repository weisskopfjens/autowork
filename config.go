package main

import (
	"encoding/json"
	"log"
	"os"
)

type configuration struct {
	Portname    string
	Speed       string
	IP          string
	Mode        string
	Repeat      int
	StartAtLine int
}

func (c *configuration) LoadConfig() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal([]byte(file), &c)
	if err != nil {
		log.Println(err)
	}
}

func (c *configuration) SaveConfig() {
	file, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile("config.json", file, 0644)
	if err != nil {
		log.Println(err)
	}
	log.Println("Configuration saved.")
}
