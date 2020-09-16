package jira

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Jira struct {
	client   http.Client
	url      string
	username string
	password string
}
type issuetype struct {
	Self             string `json:"self"`
	ID               string `json:"id"`
	Description      string `json:"description"`
	IconURL          string `json:"iconUrl"`
	Name             string `json:"name"`
	UntranslatedName string `json:"untranslatedName"`
	SubTask          bool   `json:"subtask"`
}
type project struct {
	Self       string            `json:"self"`
	ID         string            `json:"id"`
	Key        string            `json:"key"`
	Name       string            `json:"name"`
	AvatarURLs map[string]string `json:"avatarUrls"`
	IssueTypes []issuetype       `json:"issuetypes"`
}
type GetIssueMetadataOutput struct {
	Expand   string    `json:"expand"`
	Projects []project `json:"projects"`
}

func (jira *Jira) newrequest(method, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", jira.username, jira.password)))))
	return req
}

func NewJira(url, username, password string) *Jira {
	return &Jira{client: http.Client{}, username: username, password: password, url: url}
}

type GetIssueMetadataInput struct {
	ProjectIds     []string `json:"projectIds"`
	ProjectKeys    []string `json:"projectKeys"`
	IssueTypeIds   []string `json:"issuetypeIds"`
	IssueTypeNames []string `json:"issuetypeNames"`
	Expand         string   `json:"expand"`
}
type CreateIssueType struct {
	ID string `json:"id"`
}
type CreateIssueProject struct {
	ID string `json:"id"`
}
type CreateIssueInput struct {
	IssueType CreateIssueType    `json:"issuetype"`
	Summary   string             `json:"summary"`
	Project   CreateIssueProject `json:"project"`
}

type CreateIssueOutput struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}
type document struct {
	Update interface{}       `json:"update,omitempty"`
	Fields *CreateIssueInput `json:"fields"`
}

func (jira *Jira) CreateIssue(input *CreateIssueInput) (*CreateIssueOutput, error) {
	body, err := json.Marshal(document{Fields: input})
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	req := jira.newrequest("POST", fmt.Sprintf("%s/rest/api/latest/issue/", jira.url), bytes.NewReader(body))
	resp, err := jira.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Printf("%d %s", resp.StatusCode, resp.Status)
	var output CreateIssueOutput

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}
	return &output, nil
}

func (jira *Jira) GetIssueMetadata(input *GetIssueMetadataInput) (*GetIssueMetadataOutput, error) {
	// TODO add logic to handle the input
	req := jira.newrequest("GET", fmt.Sprintf("%s/rest/api/latest/issue/createmeta", jira.url), nil)
	resp, err := jira.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var output GetIssueMetadataOutput
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}
	return &output, nil
}

type GetIssuesInput struct{}

func (jira *Jira) GetIssues(input *GetIssuesInput) {
	req := jira.newrequest("GET", fmt.Sprintf("%s/rest/api/latest/issues", jira.url), nil)
	resp, err := jira.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))

}
