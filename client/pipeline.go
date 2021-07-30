package client

import (
	"errors"
	"fmt"
	"strings"
)

type ErrorResponse struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Labels struct {
	Tags []string `json:"tags,omitempty"`
}

type Metadata struct {
	Name               string `json:"name,omitempty"`
	ID                 string `json:"id,omitempty"`
	IsPublic           bool   `json:"isPublic,omitempty"`
	Labels             Labels `json:"labels,omitempty"`
	OriginalYamlString string `json:"originalYamlString,omitempty"`
	Project            string `json:"project,omitempty"`
	ProjectId          string `json:"projectId,omitempty"`
	Revision           int    `json:"revision,omitempty"`
}

type SpecTemplate struct {
	Location string `json:"location,omitempty"`
	Repo     string `json:"repo,omitempty"`
	Path     string `json:"path,omitempty"`
	Revision string `json:"revision,omitempty"`
	Context  string `json:"context,omitempty"`
}

type Trigger struct {
	Name                         string              `json:"name,omitempty"`
	Description                  string              `json:"description,omitempty"`
	Type                         string              `json:"type,omitempty"`
	Repo                         string              `json:"repo,omitempty"`
	Events                       []string            `json:"events,omitempty"`
	BranchRegex                  string              `json:"branchRegex,omitempty"`
	BranchRegexInput             string              `json:"branchRegexInput,omitempty"`
	PullRequestTargetBranchRegex string              `json:"pullRequestTargetBranchRegex,omitempty"`
	CommentRegex                 string              `json:"commentRegex,omitempty"`
	ModifiedFilesGlob            string              `json:"modifiedFilesGlob,omitempty"`
	Provider                     string              `json:"provider,omitempty"`
	Disabled                     bool                `json:"disabled,omitempty"`
	Options                      *TriggerOptions     `json:"options,omitempty"`
	PullRequestAllowForkEvents   bool                `json:"pullRequestAllowForkEvents,omitempty"`
	CommitStatusTitle            string              `json:"commitStatusTitle,omitempty"`
	Context                      string              `json:"context,omitempty"`
	Contexts                     []string            `json:"contexts,omitempty"`
	RuntimeEnvironment           *RuntimeEnvironment `json:"runtimeEnvironment,omitempty"`
	Variables                    []Variable          `json:"variables,omitempty"`
}

type TriggerOptions struct {
	NoCache             string `json:"noCache,omitempty"`
	NoCfCache           string `json:"noCfCache,omitempty"`
	ResetVolume         string `json:"resetVolume,omitempty"`
	EnableNotifications string `json:"enableNotifications,omitempty"`
}

type RuntimeEnvironment struct {
	Name        string `json:"name,omitempty"`
	Memory      string `json:"memory,omitempty"`
	CPU         string `json:"cpu,omitempty"`
	DindStorage string `json:"dindStorage,omitempty"`
}

func (t *Trigger) SetVariables(variables map[string]interface{}) {
	for key, value := range variables {
		t.Variables = append(t.Variables, Variable{Key: key, Value: value.(string)})
	}
}

type Spec struct {
	Variables          []Variable               `json:"variables,omitempty"`
	SpecTemplate       *SpecTemplate            `json:"specTemplate,omitempty"`
	Triggers           []Trigger                `json:"triggers,omitempty"`
	Priority           int                      `json:"priority,omitempty"`
	Concurrency        int                      `json:"concurrency,omitempty"`
	BranchConcurrency  int                      `json:"branchConcurrency,omitempty"`
	TriggerConcurrency int                      `json:"triggerConcurrency,omitempty"`
	Contexts           []interface{}            `json:"contexts,omitempty"`
	Steps              *Steps                   `json:"steps,omitempty"`
	Stages             *Stages                  `json:"stages,omitempty"`
	Mode               string                   `json:"mode,omitempty"`
	FailFast           *bool                    `json:"fail_fast,omitempty"`
	RuntimeEnvironment RuntimeEnvironment       `json:"runtimeEnvironment,omitempty"`
	TerminationPolicy  []map[string]interface{} `json:"terminationPolicy,omitempty"`
	Hooks              *Hooks                   `json:"hooks,omitempty"`
	Options            map[string]bool          `json:"options,omitempty"`
}

type Steps struct {
	Steps string
}
type Stages struct {
	Stages string
}

type Hooks struct {
	Hooks string
}

func (d Steps) MarshalJSON() ([]byte, error) {
	bytes := []byte(d.Steps)
	return bytes, nil
}
func (d Stages) MarshalJSON() ([]byte, error) {
	bytes := []byte(d.Stages)
	return bytes, nil
}

func (d Hooks) MarshalJSON() ([]byte, error) {
	bytes := []byte(d.Hooks)
	return bytes, nil
}

func (d *Steps) UnmarshalJSON(data []byte) error {
	d.Steps = string(data)
	return nil
}
func (d *Stages) UnmarshalJSON(data []byte) error {
	d.Stages = string(data)
	return nil
}
func (d *Hooks) UnmarshalJSON(data []byte) error {
	d.Hooks = string(data)
	return nil
}

type Pipeline struct {
	Metadata Metadata `json:"metadata,omitempty"`
	Spec     Spec     `json:"spec,omitempty"`
	Version  string   `json:"version,omitempty"`
}

func (p *Pipeline) SetVariables(variables map[string]interface{}) {
	for key, value := range variables {
		p.Spec.Variables = append(p.Spec.Variables, Variable{Key: key, Value: value.(string)})
	}
}

func (pipeline *Pipeline) GetID() string {
	if pipeline.Metadata.ID != "" {
		return pipeline.Metadata.ID
	} else {
		return pipeline.Metadata.Name
	}
}

func (client *Client) GetPipeline(name string) (*Pipeline, error) {
	fullPath := fmt.Sprintf("/pipelines/%s", strings.Replace(name, "/", "%2F", 1))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var pipeline Pipeline

	err = DecodeResponseInto(resp, &pipeline)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (client *Client) CreatePipeline(pipeline *Pipeline) (*Pipeline, error) {

	body, err := EncodeToJSON(pipeline)

	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   "/pipelines",
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var respPipeline Pipeline
	err = DecodeResponseInto(resp, &respPipeline)
	if err != nil {
		return nil, err
	}

	return &respPipeline, nil

}

func (client *Client) UpdatePipeline(pipeline *Pipeline) (*Pipeline, error) {

	body, err := EncodeToJSON(pipeline)

	if err != nil {
		return nil, err
	}

	id := pipeline.GetID()
	if id == "" {
		return nil, errors.New("[ERROR] Both Pipeline ID and Name are empty")
	}

	fullPath := fmt.Sprintf("/pipelines/%s", strings.Replace(id, "/", "%2F", 1))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "PUT",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var respPipeline Pipeline
	err = DecodeResponseInto(resp, &respPipeline)
	if err != nil {
		return nil, err
	}

	return &respPipeline, nil
}

func (client *Client) DeletePipeline(name string) error {

	fullPath := fmt.Sprintf("/pipelines/%s", strings.Replace(name, "/", "%2F", 1))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "DELETE",
	}

	_, err := client.RequestAPI(&opts)

	if err != nil {
		return err
	}

	return nil
}
