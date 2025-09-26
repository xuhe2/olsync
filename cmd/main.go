package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xuhe2/olsync/internal"
)

func main() {
	// get config yaml path from command line argument
	// if not provided, use default config.yaml
	ConfigPath := "config.yaml"
	if len(os.Args) > 1 {
		ConfigPath = os.Args[1]
	}

	config, err := internal.ParseConfigFromFile(ConfigPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := internal.NewOLClient().
		WithCookies(config.GetCookies()).
		WithProjectPageUrl(config.GetBaseURL())

	projects := client.GetProjects()
	backupProjects := make([]internal.Project, 0)
	for _, project := range projects {
		if config.ShouldBackupProject(project.Name) {
			backupProjects = append(backupProjects, project)
		}
	}

	var backupManager internal.BackupManager = internal.NewZipBackupManager(client,
		config.Backup)

	if err := backupManager.RunBackup(backupProjects); err != nil {
		log.Fatalf("Backup failed: %v", err)
	}

	log.Println("Backup finished successfully")
}
