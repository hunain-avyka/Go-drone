// Copyright 2024 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package jenkinsxml converts Jenkins XML pipelines to Harness pipelines.
package jenkinsxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	jenkinsxml "github.com/hunain-avyka/Go-drone/convert/jenkinsxml/xml"
	"github.com/hunain-avyka/Go-drone/internal/store"
	harness "github.com/hunain-avyka/go-spec/dist/go"

	"github.com/ghodss/yaml"
)

// conversion context
type context struct {
	config *jenkinsxml.Project
}

// Converter converts a Jenkins XML file to a Harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	identifiers   *store.Identifiers
}

func New(options ...Option) *Converter {
	d := new(Converter)

	// create the unique identifier store. this store
	// is used for registering unique identifiers to
	// prevent duplicate names, unique index violations.
	d.identifiers = store.New()

	// loop through and apply the options.
	for _, option := range options {
		option(d)
	}

	// set the default kubernetes namespace.
	if d.kubeNamespace == "" {
		d.kubeNamespace = "default"
	}

	// set the runtime to kubernetes if the kubernetes
	// connector is configured.
	if d.kubeConnector != "" {
		d.kubeEnabled = true
	}

	return d
}

// Convert downgrades a v1 pipeline.
func (d *Converter) Convert(r io.Reader) ([]byte, error) {
	src, err := jenkinsxml.Parse(r)
	if err != nil {
		return nil, err
	}
	return d.convert(&context{
		config: src,
	})
}

// ConvertBytes downgrades a v1 pipeline.
func (d *Converter) ConvertBytes(b []byte) ([]byte, error) {
	return d.Convert(
		bytes.NewBuffer(b),
	)
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertString(s string) ([]byte, error) {
	return d.Convert(
		bytes.NewBufferString(s),
	)
}

// ConvertFile downgrades a v1 pipeline.
func (d *Converter) ConvertFile(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return d.Convert(f)
}

// converts converts a Jenkins XML pipeline to a Harness pipeline.
func (d *Converter) convert(ctx *context) ([]byte, error) {

	// create the harness pipeline spec
	dst := &harness.Pipeline{}

	// create the harness pipeline resource
	config := &harness.Config{
		Version: 1,
		Kind:    "pipeline",
		Spec:    dst,
	}

	// cacheFound := false

	// create the harness stage.
	dstStage := &harness.Stage{
		Type: "ci",
		Name: "ci",
		// When: convertCond(from.Trigger),
		Spec: &harness.StageCI{
			// Delegate: convertNode(from.Node),
			// Envs: convertVariables(ctx.config.Variables),
			// Platform: convertPlatform(from.Platform),
			// Runtime:  convertRuntime(from),
			Steps: make([]*harness.Step, 0), // Initialize the Steps slice
		},
	}
	dst.Stages = append(dst.Stages, dstStage)
	stageSteps := make([]*harness.Step, 0)

	tasks := ctx.config.Builders.Tasks
	for _, task := range tasks {
		step := &harness.Step{}

		switch taskname := task.XMLName.Local; taskname {
		case "hudson.tasks.Shell":
			step = convertShellTaskToStep(&task)
		case "hudson.tasks.Ant":
			step = convertAntTaskToStep(&task)
		case "hudson.tasks.BatchFile":
			step = convertBatchFileTaskToStep(&task)
		case "hudson.plugins.gradle.Gradle":
			step = convertGradleFileTaskToStep(&task)
		case "hudson.tasks.Maven":
			step = convertMavenTaskToStep(&task)
		// case "com.cloudbees.jenkins.GitHubSetCommitStatusBuilder":
		//  	step = convertGitHubSetCommitStatusTaskToStep(&task)
		case "hudson.plugins.build__timeout.BuildStepWithTimeout":
			step = convertTimeoutTaskToStep(&task)

			//hudson.plugins.build__timeout.BuildStepWithTimeout
		default:
			step = unsupportedTaskToStep(taskname)
		}

		stageSteps = append(stageSteps, step)
	}
	dstStage.Spec.(*harness.StageCI).Steps = stageSteps

	// marshal the harness yaml
	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return out, nil
}
func convertMavenTaskToStep(task *jenkinsxml.Task) *harness.Step {
	mavenTask := &jenkinsxml.HudsonMavenTask{}
	// TODO: wrapping task.Content with 'builders' tags is ugly.
	err := xml.Unmarshal([]byte("<builders>"+task.Content+"</builders>"), mavenTask)
	if err != nil {
		return nil
	}

	spec := new(harness.StepPlugin)
	spec.Image = "harnesscommunitytest/maven-plugin"
	spec.Inputs = map[string]interface{}{
		"goals": mavenTask.Targets,
	}
	step := &harness.Step{
		Name: "maven",
		Type: "plugin",
		Spec: spec,
	}

	return step
}

type Buildspec struct {
	Language             string    `json:"language,omitempty"`
	BuildTool            string    `json:"buildTool,omitempty"`
	Args                 string    `json:"args,omitempty"`
	RunOnlySelectedTests bool      `json:"runOnlySelectedTests,omitempty"`
	PreCommand           string    `json:"preCommand,omitempty"`
	PostCommand          string    `json:"postCommand,omitempty"`
	Reports              []*Report `json:"reports,omitempty"`
	EnableTestSplitting  bool      `json:"enableTestSplitting,omitempty"`
}

type Report struct {
	Type string `json:"type,omitempty"`
	Spec *Spec  `json:"spec,omitempty"`
}

type Spec struct {
	Paths []string `json:"paths,omitempty"`
}

func convertGradleFileTaskToStep(task *jenkinsxml.Task) *harness.Step {
	gradleTask := &jenkinsxml.HudsonGradleTask{}
	// TODO: wrapping task.Content with 'builders' tags is ugly.
	err := xml.Unmarshal([]byte("<builders>"+task.Content+"</builders>"), gradleTask)
	if err != nil {
		return nil
	}
	// step:
	//               type: RunTests
	//               name: RunTests_1
	//               identifier: RunTests_1
	//               spec:
	//                 language: Java
	//                 buildTool: Gradle
	//                 args: clean build test
	//                 runOnlySelectedTests: false
	//                 preCommand: chmod +x ./gradlew
	//                 postCommand: buils --info
	//                 reports:
	//                   type: JUnit
	//                   spec:
	//                     paths:
	//                       - "**/*.xml"
	//                 enableTestSplitting: false

	// spec := new(harness.StepPlugin)
	// spec.Image = "harnesscommunitytest/gradle-plugin"
	// spec.Inputs = map[string]interface{}{
	// 	"goals": gradleTask.Tasks,
	// }

	spec := new(Buildspec)
	spec.Language = "Java"
	spec.BuildTool = "Gradle"
	spec.Args = gradleTask.Tasks
	spec.RunOnlySelectedTests = false
	spec.Reports = []*Report{
		{
			Type: "JUnit",
			Spec: &Spec{
				Paths: []string{"**/*.xml"},
			},
		},
	}
	spec.EnableTestSplitting = false

	step := &harness.Step{
		Name: "Gradle",
		Type: "RunTests",
		Spec: spec,
	}

	return step
}

//	func convertGitHubSetCommitStatusTaskToStep(task *jenkinsxml.Task) *harness.Step {
//		//	<com.cloudbees.jenkins.GitHubSetCommitStatusBuilder plugin="github@1.40.0">
//		//
//		// <statusMessage>
//		// <content/>
//		// </statusMessage>
//		// <contextSource class="org.jenkinsci.plugins.github.status.sources.DefaultCommitContextSource"/>
//		// </com.cloudbees.jenkins.GitHubSetCommitStatusBuilder>
//	}

func convertTimeoutTaskToStep(task *jenkinsxml.Task) *harness.Step {
	// Create a new HudsonTimeoutTask instance
	timeoutTask := &jenkinsxml.HudsonTimeoutTask{}

	// Trim the content from the task
	xmlContent := strings.TrimSpace(task.Content)

	// Unmarshal the XML content into the timeoutTask
	err := xml.Unmarshal([]byte(xmlContent), timeoutTask)
	if err != nil {
		log.Println("Error unmarshalling XML:", err)
		return nil
	}
	timeoutValue := fmt.Sprintf("%sm", timeoutTask.TimeoutMinutes)

	// Create StepExec object and set the timeout value to the Run field
	spec := new(harness.StepExec)
	//spec.Run = timeoutTask.TimeoutMinutes
	step := &harness.Step{
		Spec:    spec,
		Name:    "Run0",
		Type:    "script",
		Timeout: timeoutValue,
	}

	return step
}

func convertBatchFileTaskToStep(task *jenkinsxml.Task) *harness.Step {
	batchTask := &jenkinsxml.HudsonBatchTask{}
	err := xml.Unmarshal([]byte("<builders>"+task.Content+"</builders>"), batchTask)
	if err != nil {
		return nil
	}
	spec := new(harness.StepExec)
	spec.Run = batchTask.Command
	step := &harness.Step{
		Name: "shell",
		Type: "script",
		Spec: spec,
	}

	return step

}

// convertAntTaskToStep converts a Jenkins Ant task to a Harness step.
func convertAntTaskToStep(task *jenkinsxml.Task) *harness.Step {
	antTask := &jenkinsxml.HudsonAntTask{}
	// TODO: wrapping task.Content with 'builders' tags is ugly.
	err := xml.Unmarshal([]byte("<builders>"+task.Content+"</builders>"), antTask)
	if err != nil {
		return nil
	}

	spec := new(harness.StepPlugin)
	spec.Image = "harnesscommunitytest/ant-plugin"
	spec.Inputs = map[string]interface{}{
		"goals": antTask.Targets,
	}
	step := &harness.Step{
		Name: "ant",
		Type: "plugin",
		Spec: spec,
	}

	return step
}

func convertShellTaskToStep(task *jenkinsxml.Task) *harness.Step {
	shellTask := &jenkinsxml.HudsonShellTask{}
	// TODO: wrapping task.Content with 'builders' tags is ugly.
	err := xml.Unmarshal([]byte("<builders>"+task.Content+"</builders>"), shellTask)
	if err != nil {
		return nil
	}

	spec := new(harness.StepExec)
	spec.Run = shellTask.Command
	step := &harness.Step{
		Name: "shell",
		Type: "script",
		Spec: spec,
	}

	return step
}

// func convertGitHubSetCommitStatusTaskToStep(task *jenkinsxml.Task) *harness.Step {
// 	shellTask := &jenkinsxml.ConvertGitHubSetCommitStatusTask{}
// 	// TODO: wrapping task.Content with 'builders' tags is ugly.
// 	err := xml.Unmarshal([]byte("<builders>"+task.Content+"</builders>"), shellTask)
// 	if err != nil {
// 		return nil
// 	}

// 	spec := new(harness.StepExec)
// 	spec.Run = shellTask.Command
// 	step := &harness.Step{
// 		Name: "shell",
// 		Type: "script",
// 		Spec: spec,
// 	}

// 	return step
// }

func unsupportedTaskToStep(task string) *harness.Step {
	spec := new(harness.StepExec)
	spec.Run = "echo Unsupported field " + task
	step := &harness.Step{
		Name: "shell",
		Type: "script",
		Spec: spec,
	}

	return step
}
