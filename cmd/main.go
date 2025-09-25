package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/xuhe2/olsync/internal"
)

func main() {
	config, err := internal.ParseConfigFromFile("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	client := internal.NewOLClient().
		WithCookies(config.GetCookies()).
		WithProjectPageUrl(config.GetBaseURL())

	projects := client.GetProjects()

	// try to download the first project

	log.Printf("Found %d projects", len(projects))
	log.Printf("First project: %s", projects[0].Name)
	p1 := projects[0]

	reader, err := client.DownloadProjectZip(p1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer reader.Close()

	// TODO: save to file to disk
	filePath := fmt.Sprintf("%s.zip", p1.Name)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	copied, err := io.Copy(file, reader)
	fmt.Println("Downloaded project to", filePath)
	fmt.Println(copied)
}
