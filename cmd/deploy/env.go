package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	b, err := ioutil.ReadFile("./app.deploy.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lines := string(b)
	projectID := os.Args[1]
	githubToken := os.Args[2]
	webhookSecret := os.Args[3]
	lines = strings.Replace(lines, "##GCP_PROJECT", projectID, 1)
	lines = strings.Replace(lines, "##GH_SECRET_TOKEN", githubToken, 1)
	lines = strings.Replace(lines, "##GH_WEBHOOK_SECRET", webhookSecret, 1)
	err = ioutil.WriteFile("./app.yaml", []byte(lines), 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
