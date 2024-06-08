package main

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

func main() {
	projects, err := initProjects()
	if err != nil {
		slog.Error("Failed: unable to read the projects Metadata", "Reason", err)
		os.Exit(1)
	}

	slog.Info("Got projects metadata", "Projects", projects)

	storage, err := NewGHRepoVarStorage()
	if err != nil {
		slog.Error("Failed: initialize storage", "Reason", err)
		os.Exit(1)
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(projects.Projects))

	errChan := make(chan error, len(projects.Projects))

	for _, project := range projects.Projects {
		go func(proj Project) {
			defer wg.Done()

			latestRelease, err := fetchLatestRelease(
				proj.Organization,
				proj.Name,
			)
			if err != nil {
				errChan <- err
				return
			}

			storedRelease, err := storage.ReadRepoVariable(proj.Name)
			if err != nil {
				errChan <- err
				return
			}

			if strings.Compare(*storedRelease, *latestRelease) != 0 {
				slog.Info("Release Versions", "Proj", proj.Name, "Org", proj.Organization, "Latest", *latestRelease, "Current", *storedRelease)
				if err := storage.UpdateRepoVariable(proj.Name, *latestRelease); err != nil {
					errChan <- err
					return
				}
				slog.Info("Updated to Latest version", "Proj", proj.Name, "Org", proj.Organization, "Ver", *latestRelease)

				if err := storage.DispatchCall(
					proj,
					*latestRelease,
				); err != nil {
					errChan <- err
					return
				}

			} else {
				slog.Info("Already in Latest version", "Proj", proj.Name, "Org", proj.Organization, "Ver", *latestRelease)
			}
		}(project)
	}

	wg.Wait()
	close(errChan)

	hadFailures := false
	for err := range errChan {
		if err != nil {
			hadFailures = true
			slog.Error("Failed", "Reason", err)
		}
	}
	if hadFailures {
		os.Exit(1)
	}
}
