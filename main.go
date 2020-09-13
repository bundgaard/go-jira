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

func main() {
	// https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-get
	fmt.Printf("%s %s\n", Version, Project)
	url := os.Getenv("JIRA_URL")
	username := os.Getenv("JIRA_USERNAME")
	password := os.Getenv("JIRA_PASSWORD")

	client := jira.NewJira(url, username, password)

	input := jira.GetIssueMetadataInput{}
	output, err := client.GetIssueMetadata(&input)
	if err != nil {
		log.Fatal(err)
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
