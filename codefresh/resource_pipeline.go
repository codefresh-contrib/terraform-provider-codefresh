package codefresh

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	ghodss "github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"gopkg.in/yaml.v2"
)

var terminationPolicyOnCreateBranchAttributes = []string{"branchName", "ignoreTrigger", "ignoreBranch"}

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineCreate,
		Read:   resourcePipelineRead,
		Update: resourcePipelineUpdate,
		Delete: resourcePipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"original_yaml_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"spec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"concurrency": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0, // zero is unlimited
						},
						"branch_concurrency": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0, // zero is unlimited
						},
						"trigger_concurrency": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0, // zero is unlimited
						},
						"spec_template": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"location": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "git",
									},
									"repo": {
										Type:     schema.TypeString,
										Required: true,
									},
									"path": {
										Type:     schema.TypeString,
										Required: true,
									},
									"revision": {
										Type:     schema.TypeString,
										Required: true,
									},
									"context": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "github",
									},
								},
							},
						},
						"variables": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"trigger": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "git",
									},
									"repo": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"branch_regex": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "/.*/gi",
										ValidateFunc: stringIsValidRe2RegExp,
									},
									"branch_regex_input": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "regex",
										ValidateFunc: validation.StringInSlice([]string{"multiselect-exclude", "multiselect", "regex"}, false),
									},
									"pull_request_target_branch_regex": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: stringIsValidRe2RegExp,
									},
									"comment_regex": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "/.*/gi",
										ValidateFunc: stringIsValidRe2RegExp,
									},
									"modified_files_glob": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"events": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"provider": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "github",
									},
									"disabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"options": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"no_cache": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"no_cf_cache": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"reset_volume": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"enable_notifications": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
											},
										},
									},
									"pull_request_allow_fork_events": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"commit_status_title": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"context": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "github",
									},
									"contexts": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"runtime_environment": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"memory": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"cpu": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"dind_storage": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"variables": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"contexts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"termination_policy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"on_create_branch": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"branch_name": {
													Type:          schema.TypeString,
													Optional:      true,
													ValidateFunc:  stringIsValidRe2RegExp,
													ConflictsWith: []string{"spec.0.termination_policy.0.on_create_branch.0.ignore_branch"},
												},
												"ignore_trigger": {
													Optional: true,
													Type:     schema.TypeBool,
												},
												"ignore_branch": {
													Optional: true,
													Type:     schema.TypeBool,
												},
											},
										},
									},
									"on_terminate_annotation": {
										Optional: true,
										Type:     schema.TypeBool,
										Default:  false,
									},
								},
							},
						},
						"runtime_environment": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"memory": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cpu": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dind_storage": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"options": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"keep_pvcs_for_pending_approval": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"pending_approval_concurrency_applied": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourcePipelineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	pipeline := *mapResourceToPipeline(d)

	resp, err := client.CreatePipeline(&pipeline)
	if err != nil {
		return err
	}

	d.SetId(resp.Metadata.ID)

	return resourcePipelineRead(d, meta)
}

func resourcePipelineRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	pipelineID := d.Id()

	if pipelineID == "" {
		d.SetId("")
		return nil
	}

	pipeline, err := client.GetPipeline(pipelineID)
	if err != nil {
		return err
	}

	err = mapPipelineToResource(*pipeline, d)
	if err != nil {
		return err
	}

	return nil
}

func resourcePipelineUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	pipeline := *mapResourceToPipeline(d)
	pipeline.Metadata.ID = d.Id()

	_, err := client.UpdatePipeline(&pipeline)
	if err != nil {
		return err
	}

	return resourcePipelineRead(d, meta)
}

func resourcePipelineDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*cfClient.Client)

	err := client.DeletePipeline(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func mapPipelineToResource(pipeline cfClient.Pipeline, d *schema.ResourceData) error {

	err := d.Set("name", pipeline.Metadata.Name)
	if err != nil {
		return err
	}

	err = d.Set("revision", pipeline.Metadata.Revision)
	if err != nil {
		return err
	}

	err = d.Set("project_id", pipeline.Metadata.ProjectId)
	if err != nil {
		return err
	}

	err = d.Set("is_public", pipeline.Metadata.IsPublic)
	if err != nil {
		return err
	}

	err = d.Set("spec", flattenSpec(pipeline.Spec))
	if err != nil {
		return err
	}

	err = d.Set("tags", pipeline.Metadata.Labels.Tags)
	if err != nil {
		return err
	}

	err = d.Set("original_yaml_string", pipeline.Metadata.OriginalYamlString)
	if err != nil {
		return err
	}

	return nil
}

func flattenSpec(spec cfClient.Spec) []interface{} {

	var res = make([]interface{}, 0)
	m := make(map[string]interface{})

	if len(spec.Triggers) > 0 {
		m["trigger"] = flattenTriggers(spec.Triggers)
	}

	if spec.SpecTemplate != nil {
		m["spec_template"] = flattenSpecTemplate(*spec.SpecTemplate)
	}

	if len(spec.Variables) != 0 {
		m["variables"] = convertVariables(spec.Variables)
	}

	if spec.RuntimeEnvironment != (cfClient.RuntimeEnvironment{}) {
		m["runtime_environment"] = flattenSpecRuntimeEnvironment(spec.RuntimeEnvironment)
	}

	if len(spec.TerminationPolicy) > 0 {
		m["termination_policy"] = flattenSpecTerminationPolicy(spec.TerminationPolicy)
	}

	if len(spec.Options) > 0 {
		var resOptions []map[string]bool
		options := map[string]bool{}
		for keyOption, valueOption := range spec.Options {
			switch {
			case keyOption == "keepPVCsForPendingApproval":
				options["keep_pvcs_for_pending_approval"] = valueOption
			case keyOption == "pendingApprovalConcurrencyApplied":
				options["pending_approval_concurrency_applied"] = valueOption
			}
		}
		resOptions = append(resOptions, options)
		m["options"] = resOptions
	}

	m["concurrency"] = spec.Concurrency
	m["branch_concurrency"] = spec.BranchConcurrency
	m["trigger_concurrency"] = spec.TriggerConcurrency

	m["priority"] = spec.Priority

	m["contexts"] = spec.Contexts

	res = append(res, m)
	return res
}

func flattenSpecTerminationPolicy(terminationPolicy []map[string]interface{}) []map[string]interface{} {
	var res []map[string]interface{}
	attribute := make(map[string]interface{})
	for _, policy := range terminationPolicy {
		eventName, _ := policy["event"]
		typeName, _ := policy["type"]
		attributeName := convertOnCreateBranchAttributeToPipelineFormat(eventName.(string) + "_" + typeName.(string))
		switch attributeName {
		case "on_create_branch":
			var valueList []map[string]interface{}
			attributeValues := make(map[string]interface{})
			for _, eventAttribute := range terminationPolicyOnCreateBranchAttributes {
				if item, ok := policy[eventAttribute]; ok {
					attributeValues[convertOnCreateBranchAttributeToPipelineFormat(eventAttribute)] = item
				}
			}
			attribute[attributeName] = append(valueList, attributeValues)
		case "on_terminate_annotation":
			if value, ok := policy["key"]; ok && value == "cf_predecessor" {
				attribute[attributeName] = true
			}
		default:
			log.Fatal("Unsupported event found in TerminationPolicy")
		}
	}
	res = append(res, attribute)
	return res
}

func flattenSpecTemplate(spec cfClient.SpecTemplate) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"location": spec.Location,
			"repo":     spec.Repo,
			"context":  spec.Context,
			"revision": spec.Revision,
			"path":     spec.Path,
		},
	}
}

func flattenSpecRuntimeEnvironment(spec cfClient.RuntimeEnvironment) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":         spec.Name,
			"memory":       spec.Memory,
			"cpu":          spec.CPU,
			"dind_storage": spec.DindStorage,
		},
	}
}

func flattenTriggerOptions(options cfClient.TriggerOptions) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"no_cache":             options.NoCache,
			"no_cf_cache":          options.NoCfCache,
			"reset_volume":         options.ResetVolume,
			"enable_notifications": options.EnableNotifications,
		},
	}
}

func flattenTriggers(triggers []cfClient.Trigger) []map[string]interface{} {
	var res = make([]map[string]interface{}, len(triggers))
	for i, trigger := range triggers {
		m := make(map[string]interface{})
		m["name"] = trigger.Name
		m["description"] = trigger.Description
		m["context"] = trigger.Context
		m["contexts"] = trigger.Contexts
		m["repo"] = trigger.Repo
		m["branch_regex"] = trigger.BranchRegex
		m["branch_regex_input"] = trigger.BranchRegexInput
		m["pull_request_target_branch_regex"] = trigger.PullRequestTargetBranchRegex
		m["comment_regex"] = trigger.CommentRegex
		m["modified_files_glob"] = trigger.ModifiedFilesGlob
		m["disabled"] = trigger.Disabled
		if trigger.Options != nil {
			m["options"] = flattenTriggerOptions(*trigger.Options)
		}
		m["pull_request_allow_fork_events"] = trigger.PullRequestAllowForkEvents
		m["commit_status_title"] = trigger.CommitStatusTitle
		m["provider"] = trigger.Provider
		m["type"] = trigger.Type
		m["events"] = trigger.Events
		m["variables"] = convertVariables(trigger.Variables)
		if trigger.RuntimeEnvironment != nil {
			m["runtime_environment"] = flattenSpecRuntimeEnvironment(*trigger.RuntimeEnvironment)
		}
		res[i] = m
	}
	return res
}

func mapResourceToPipeline(d *schema.ResourceData) *cfClient.Pipeline {

	tags := d.Get("tags").(*schema.Set).List()

	originalYamlString := strings.Replace(
		d.Get("original_yaml_string").(string),
		"\n",
		"\n",
		-1)
	pipeline := &cfClient.Pipeline{
		Metadata: cfClient.Metadata{
			Name:      d.Get("name").(string),
			Revision:  d.Get("revision").(int),
			ProjectId: d.Get("project_id").(string),
			IsPublic:  d.Get("is_public").(bool),
			Labels: cfClient.Labels{
				Tags: convertStringArr(tags),
			},
			OriginalYamlString: originalYamlString,
		},
		Spec: cfClient.Spec{
			Priority:           d.Get("spec.0.priority").(int),
			Concurrency:        d.Get("spec.0.concurrency").(int),
			BranchConcurrency:  d.Get("spec.0.branch_concurrency").(int),
			TriggerConcurrency: d.Get("spec.0.trigger_concurrency").(int),
		},
	}

	if _, ok := d.GetOk("spec.0.spec_template"); ok {
		pipeline.Spec.SpecTemplate = &cfClient.SpecTemplate{
			Location: d.Get("spec.0.spec_template.0.location").(string),
			Repo:     d.Get("spec.0.spec_template.0.repo").(string),
			Path:     d.Get("spec.0.spec_template.0.path").(string),
			Revision: d.Get("spec.0.spec_template.0.revision").(string),
			Context:  d.Get("spec.0.spec_template.0.context").(string),
		}
	} else {
		extractSpecAttributesFromOriginalYamlString(originalYamlString, pipeline)
	}

	if _, ok := d.GetOk("spec.0.runtime_environment"); ok {
		pipeline.Spec.RuntimeEnvironment = cfClient.RuntimeEnvironment{
			Name:        d.Get("spec.0.runtime_environment.0.name").(string),
			Memory:      d.Get("spec.0.runtime_environment.0.memory").(string),
			CPU:         d.Get("spec.0.runtime_environment.0.cpu").(string),
			DindStorage: d.Get("spec.0.runtime_environment.0.dind_storage").(string),
		}
	}

	contexts := d.Get("spec.0.contexts").([]interface{})
	pipeline.Spec.Contexts = contexts

	variables := d.Get("spec.0.variables").(map[string]interface{})
	pipeline.SetVariables(variables)

	triggers := d.Get("spec.0.trigger").([]interface{})
	for idx := range triggers {
		events := d.Get(fmt.Sprintf("spec.0.trigger.%v.events", idx)).([]interface{})
		contexts := d.Get(fmt.Sprintf("spec.0.trigger.%v.contexts", idx)).([]interface{})
		codefreshTrigger := cfClient.Trigger{
			Name:                         d.Get(fmt.Sprintf("spec.0.trigger.%v.name", idx)).(string),
			Description:                  d.Get(fmt.Sprintf("spec.0.trigger.%v.description", idx)).(string),
			Type:                         d.Get(fmt.Sprintf("spec.0.trigger.%v.type", idx)).(string),
			Repo:                         d.Get(fmt.Sprintf("spec.0.trigger.%v.repo", idx)).(string),
			BranchRegex:                  d.Get(fmt.Sprintf("spec.0.trigger.%v.branch_regex", idx)).(string),
			BranchRegexInput:             d.Get(fmt.Sprintf("spec.0.trigger.%v.branch_regex_input", idx)).(string),
			PullRequestTargetBranchRegex: d.Get(fmt.Sprintf("spec.0.trigger.%v.pull_request_target_branch_regex", idx)).(string),
			CommentRegex:                 d.Get(fmt.Sprintf("spec.0.trigger.%v.comment_regex", idx)).(string),
			ModifiedFilesGlob:            d.Get(fmt.Sprintf("spec.0.trigger.%v.modified_files_glob", idx)).(string),
			Provider:                     d.Get(fmt.Sprintf("spec.0.trigger.%v.provider", idx)).(string),
			Disabled:                     d.Get(fmt.Sprintf("spec.0.trigger.%v.disabled", idx)).(bool),
			PullRequestAllowForkEvents:   d.Get(fmt.Sprintf("spec.0.trigger.%v.pull_request_allow_fork_events", idx)).(bool),
			CommitStatusTitle:            d.Get(fmt.Sprintf("spec.0.trigger.%v.commit_status_title", idx)).(string),
			Context:                      d.Get(fmt.Sprintf("spec.0.trigger.%v.context", idx)).(string),
			Contexts:                     convertStringArr(contexts),
			Events:                       convertStringArr(events),
		}
		variables := d.Get(fmt.Sprintf("spec.0.trigger.%v.variables", idx)).(map[string]interface{})
		codefreshTrigger.SetVariables(variables)
		if _, ok := d.GetOk(fmt.Sprintf("spec.0.trigger.%v.options", idx)); ok {
			options := cfClient.TriggerOptions{
				NoCache:             d.Get(fmt.Sprintf("spec.0.trigger.%v.options.0.no_cache", idx)).(string),
				NoCfCache:           d.Get(fmt.Sprintf("spec.0.trigger.%v.options.0.no_cf_cache", idx)).(string),
				ResetVolume:         d.Get(fmt.Sprintf("spec.0.trigger.%v.options.0.reset_volume", idx)).(string),
				EnableNotifications: d.Get(fmt.Sprintf("spec.0.trigger.%v.options.0.enable_notifications", idx)).(string),
			}
			codefreshTrigger.Options = &options
		}
		if _, ok := d.GetOk(fmt.Sprintf("spec.0.trigger.%v.runtime_environment", idx)); ok {
			triggerRuntime := cfClient.RuntimeEnvironment{
				Name:        d.Get(fmt.Sprintf("spec.0.trigger.%v.runtime_environment.0.name", idx)).(string),
				Memory:      d.Get(fmt.Sprintf("spec.0.trigger.%v.runtime_environment.0.memory", idx)).(string),
				CPU:         d.Get(fmt.Sprintf("spec.0.trigger.%v.runtime_environment.0.cpu", idx)).(string),
				DindStorage: d.Get(fmt.Sprintf("spec.0.trigger.%v.runtime_environment.0.dind_storage", idx)).(string),
			}
			codefreshTrigger.RuntimeEnvironment = &triggerRuntime
		}
		pipeline.Spec.Triggers = append(pipeline.Spec.Triggers, codefreshTrigger)
	}

	var codefreshTerminationPolicy []map[string]interface{}

	if _, ok := d.GetOk("spec.0.termination_policy.0.on_create_branch"); ok {
		var onCreatBranchPolicy = make(map[string]interface{})
		onCreatBranchPolicy = getSupportedTerminationPolicyAttributes("on_create_branch")
		for _, attribute := range terminationPolicyOnCreateBranchAttributes {
			if attributeValue, ok := d.GetOk(fmt.Sprintf("spec.0.termination_policy.0.on_create_branch.0.%s", convertOnCreateBranchAttributeToPipelineFormat(attribute))); ok {
				onCreatBranchPolicy[attribute] = attributeValue
			}
		}
		codefreshTerminationPolicy = append(codefreshTerminationPolicy, onCreatBranchPolicy)
	}
	if _, ok := d.GetOk("spec.0.termination_policy.0.on_terminate_annotation"); ok {
		var onTerminateAnnotationPolicy = make(map[string]interface{})
		onTerminateAnnotationPolicy = getSupportedTerminationPolicyAttributes("on_terminate_annotation")
		onTerminateAnnotationPolicy["key"] = "cf_predecessor"
		codefreshTerminationPolicy = append(codefreshTerminationPolicy, onTerminateAnnotationPolicy)
	}
	if _, ok := d.GetOk("spec.0.options"); ok {
		pipelineSpecOption := make(map[string]bool)
		if keepPVCs, ok := d.GetOkExists("spec.0.options.0.keep_pvcs_for_pending_approval"); ok {
			pipelineSpecOption["keepPVCsForPendingApproval"] = keepPVCs.(bool)
		}
		if pendingApprovalConcurrencyApplied, ok := d.GetOkExists("spec.0.options.0.pending_approval_concurrency_applied"); ok {
			pipelineSpecOption["pendingApprovalConcurrencyApplied"] = pendingApprovalConcurrencyApplied.(bool)
		}
		pipeline.Spec.Options = pipelineSpecOption
	} else {
		pipeline.Spec.Options = nil
	}

	pipeline.Spec.TerminationPolicy = codefreshTerminationPolicy

	return pipeline
}

// extractSpecAttributesFromOriginalYamlString extracts the steps and stages from the original yaml string to enable propagation in the `Spec` attribute of the pipeline
// We cannot leverage on the standard marshal/unmarshal because the steps attribute needs to maintain the order of elements
// while by default the standard function doesn't do it because in JSON maps are unordered
func extractSpecAttributesFromOriginalYamlString(originalYamlString string, pipeline *cfClient.Pipeline) {
	// Use mapSlice to preserve order of items from the YAML string
	m := yaml.MapSlice{}
	err := yaml.Unmarshal([]byte(originalYamlString), &m)
	if err != nil {
		log.Fatalf("Unable to unmarshall original_yaml_string. Error: %v", err)
	}

	stages := "[]"
	// Dynamically build JSON object for steps using String builder
	stepsBuilder := strings.Builder{}
	stepsBuilder.WriteString("{")
	// Dynamically build JSON object for steps using String builder
	hooksBuilder := strings.Builder{}
	hooksBuilder.WriteString("{")

	// Parse elements of the YAML string to extract Steps and Stages if defined
	for _, item := range m {
		key := item.Key.(string)
		switch key {
		case "steps":
			switch x := item.Value.(type) {
			default:
				log.Fatalf("unsupported value type: %T", item.Value)

			case yaml.MapSlice:
				numberOfSteps := len(x)
				for index, item := range x {
					// We only need to preserve order at the first level to guarantee order of the steps, hence the child nodes can be marshalled
					// with the standard library
					y, _ := yaml.Marshal(item.Value)
					j2, _ := ghodss.YAMLToJSON(y)
					stepsBuilder.WriteString("\"" + item.Key.(string) + "\" : " + string(j2))
					if index < numberOfSteps-1 {
						stepsBuilder.WriteString(",")
					}
				}
			}
		case "stages":
			// For Stages we don't have ordering issue because it's a list
			y, _ := yaml.Marshal(item.Value)
			j2, _ := ghodss.YAMLToJSON(y)
			stages = string(j2)
		case "hooks":
			switch hooks := item.Value.(type) {
			default:
				log.Fatalf("unsupported value type: %T", item.Value)

			case yaml.MapSlice:
				numberOfHooks := len(hooks)
				for indexHook, hook := range hooks {
					// E.g. on_finish
					hooksBuilder.WriteString("\"" + hook.Key.(string) + "\" : {")
					numberOfAttributes := len(hook.Value.(yaml.MapSlice))
					for indexAttribute, hookAttribute := range hook.Value.(yaml.MapSlice) {
						attribute := hookAttribute.Key.(string)
						switch attribute {
						case "steps":
							hooksBuilder.WriteString("\"steps\" : {")
							numberOfSteps := len(hookAttribute.Value.(yaml.MapSlice))
							for indexStep, step := range hookAttribute.Value.(yaml.MapSlice) {
								// We only need to preserve order at the first level to guarantee order of the steps, hence the child nodes can be marshalled
								// with the standard library
								y, _ := yaml.Marshal(step.Value)
								j2, _ := ghodss.YAMLToJSON(y)
								hooksBuilder.WriteString("\"" + step.Key.(string) + "\" : " + string(j2))
								if indexStep < numberOfSteps-1 {
									hooksBuilder.WriteString(",")
								}
							}
							hooksBuilder.WriteString("}")
						default:
							// For Other elements we don't need to preserve order
							y, _ := yaml.Marshal(hookAttribute.Value)
							j2, _ := ghodss.YAMLToJSON(y)
							hooksBuilder.WriteString("\"" + hookAttribute.Key.(string) + "\" : " + string(j2))
						}

						if indexAttribute < numberOfAttributes-1 {
							hooksBuilder.WriteString(",")
						}
					}
					hooksBuilder.WriteString("}")
					if indexHook < numberOfHooks-1 {
						hooksBuilder.WriteString(",")
					}
				}
			}
		case "mode":
			pipeline.Spec.Mode = item.Value.(string)
		case "fail_fast":
			ff, ok := item.Value.(bool)
			if ok {
				pipeline.Spec.FailFast = &ff
			}
		default:
			log.Printf("Unsupported entry %s", key)
		}
	}
	stepsBuilder.WriteString("}")
	hooksBuilder.WriteString("}")
	steps := stepsBuilder.String()
	hooks := hooksBuilder.String()
	pipeline.Spec.Steps = &cfClient.Steps{
		Steps: steps,
	}
	pipeline.Spec.Stages = &cfClient.Stages{
		Stages: stages,
	}
	pipeline.Spec.Hooks = &cfClient.Hooks{
		Hooks: hooks,
	}

}

func getSupportedTerminationPolicyAttributes(policy string) map[string]interface{} {
	switch policy {
	case "on_create_branch":
		return map[string]interface{}{"event": "onCreate", "type": "branch"}
	case "on_terminate_annotation":
		return map[string]interface{}{"event": "onTerminate", "type": "annotation"}
	default:
		log.Fatal("Invalid termination policy selected: ", policy)
	}
	return nil
}

func convertOnCreateBranchAttributeToResourceFormat(src string) string {
	var re = regexp.MustCompile(`_[a-z]`)
	return re.ReplaceAllStringFunc(src, func(w string) string {
		return strings.ToUpper(strings.ReplaceAll(w, "_", ""))
	})
}

func convertOnCreateBranchAttributeToPipelineFormat(src string) string {
	var re = regexp.MustCompile(`[A-Z]`)
	return re.ReplaceAllStringFunc(src, func(w string) string {
		return "_" + strings.ToLower(w)
	})
}
