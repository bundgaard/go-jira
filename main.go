package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bundgaard/go-jira/pkg/jira"
)

var (
	Version = "1.0.0"
	Project = "go-jira"
)

type config struct {
	url      string
	username string
	password string
}

func main() {
	// https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-get
	fmt.Printf("%s %s\n", Version, Project)
	c := config{url: os.Getenv("JIRA_URL"), username: os.Getenv("JIRA_USERNAME"), password: os.Getenv("JIRA_PASSWORD")}

	os.Exit(0)
	if len(c.url) < 1 {
		fmt.Printf("error: url is not defined")
		os.Exit(1)
	}

	client := jira.NewJira(c.url, c.username, c.password)

	input := jira.GetIssueMetadataInput{}
	output, err := client.GetIssueMetadata(&input)
	if err != nil {
		log.Fatal("metadata", err)
	}
	fmt.Println(output)

	client.CreateIssue(
		&jira.CreateIssueInput{
			Project:   jira.CreateIssueProject{ID: "10000"},
			IssueType: jira.CreateIssueType{ID: "10004"},
			Summary:   "A short summary of the user story",
		})
	client.GetIssues(nil)

}
