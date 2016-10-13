package main

import (
	"fmt"
	"strings"
	"time"

	collectors "github.com/ArjenSchwarz/journal_collectors/shared"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// checkDurationHours is the time you want to look back
var checkDuration = "24h"

// username is the name of the user whose contributions you want to check
var username = "ArjenSchwarz"

// encryptedToken is the KMS encrypted GitHub token
var encryptedToken = "AQECAHgL87diS+gtOws91E0gf121Z+yeMlGmU2l8DBRhbwyFAAAAAIcwgYQGCSqGSIb3DQEHBqB3MHUCAQAwcAYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAybQK6Y6tkyJ7Yz638CARCAQ7nmYToCzd+8ZiFAWHOJH0EXfzBa1ZRrZEBwpPk+3YeIHqAL4n1NZ/JRSZ7L7L8efPS02JRdfH4l55SUNY4zUjoH9iE="

// dateformat is the format for the date output in the header
var dateformat = "01-02-2006"

// ifttTriggerName is the name of the trigger you set up in your IFTTT
var iftttTriggerName = "dayone_github_collector"

// encryptedIFTTKey is the KMS encrypted IFTTT maker channel key
var encryptedIFTTTKey = "AQECAHgL87diS+gtOws91E0gf121Z+yeMlGmU2l8DBRhbwyFAAAAAIowgYcGCSqGSIb3DQEHBqB6MHgCAQAwcwYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAw1Z0cZWDFjNpC/73wCARCARoZnSB73hWARMeHzJ7v166MysCYzgRxWsoRZZ5JbmoW00NvoGUXBWyknd59O/E6kPZlcBGyXH+s9qeooUwLluDDiunMxvn8="

// entry struct for holding
type entry struct {
	commits []*github.PushEvent
}

type commitlist struct {
	entries map[string]entry
}

func main() {
	checkDuration, err := time.ParseDuration(checkDuration)
	if err != nil {
		panic(err)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: collectors.DecryptKMS(encryptedToken)},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	events, _, err := client.Activity.ListEventsPerformedByUser(username, false, nil)
	if err != nil {
		panic(err)
	}
	commits := commitlist{entries: make(map[string]entry)}
	compareTime := time.Now().Add(-checkDuration)
	for _, event := range events {
		if *event.Type == "PushEvent" && compareTime.Before(*event.CreatedAt) {
			reponame := *event.Repo.Name
			entries := entry{}
			if val, ok := commits.entries[reponame]; ok {
				entries = val
			}
			payload := event.Payload()
			if pushEvent, ok := payload.(*github.PushEvent); ok {
				entries.commits = append(entries.commits, pushEvent)
			}
			commits.entries[reponame] = entries
		}
	}
	output := formatOutput(commits)
	if output != "" {
		collectors.SendToIFTTT(encryptedIFTTTKey, iftttTriggerName, []string{output})
	}
}

func formatOutput(commits commitlist) string {
	output := ""
	if len(commits.entries) == 0 {
		return output
	}
	output += fmt.Sprintf("# GitHub commits for %s", time.Now().Format(dateformat))
	for reponame, entry := range commits.entries {
		output += fmt.Sprintf("<br><br>## %s <br>", reponame)
		for _, event := range entry.commits {
			for _, commit := range event.Commits {
				output += fmt.Sprintf("<br> * %v (%v)",
					strings.Replace(*commit.Message, "\n", "<br>", -1),
					strings.Replace(*event.Ref, "refs/heads/", "", -1))
			}
		}
	}
	return output
}
