package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ArjenSchwarz/journal-collectors/config"
	collectors "github.com/ArjenSchwarz/journal-collectors/shared"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type configuration struct {
	CheckDuration    string
	GitHubUsername   string
	GitHubToken      string
	Dateformat       string
	IFTTTTriggerName string
	IFTTTKey         string
	KMS              bool
}

var settings configuration

// entry struct for holding
type entry struct {
	commits []*github.PushEvent
}

type commitlist struct {
	entries map[string]entry
}

func main() {
	lambda.Start(handler)
}

func handler() error {
	err := config.ParseConfig(&settings)
	if err != nil {
		panic(err)
	}
	if settings.KMS {
		settings.GitHubToken = collectors.DecryptKMS(settings.GitHubToken)
		settings.IFTTTKey = collectors.DecryptKMS(settings.IFTTTKey)
	}
	checkDuration, err := time.ParseDuration(settings.CheckDuration)
	if err != nil {
		panic(err)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: settings.GitHubToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	events, _, err := client.Activity.ListEventsPerformedByUser(context.Background(), settings.GitHubUsername, false, nil)
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
		collectors.SendToIFTTT(settings.IFTTTKey, settings.IFTTTTriggerName, []string{output})
	}
	return nil
}

func formatOutput(commits commitlist) string {
	output := ""
	if len(commits.entries) == 0 {
		return output
	}
	output += fmt.Sprintf("# GitHub commits for %s", time.Now().Format(settings.Dateformat))
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
