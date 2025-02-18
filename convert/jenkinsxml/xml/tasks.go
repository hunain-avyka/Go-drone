// // Copyright 2024 Harness, Inc.
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// //      http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.

// package xml

// type (
// 	HudsonShellTask struct {
// 		Command              string `xml:"command"`
// 		ConfiguredLocalRules string `xml:"configuredLocalRules,omitempty"`
// 	}

// 	HudsonAntTask struct {
// 		Plugin  string `xml:"plugin,attr"`
// 		Targets string `xml:"targets"`
// 	}

// )
package xml

type (
	HudsonShellTask struct {
		Command              string `xml:"command"`
		ConfiguredLocalRules string `xml:"configuredLocalRules,omitempty"`
	}
	HudsonAntTask struct {
		Plugin  string `xml:"plugin,attr"`
		Targets string `xml:"targets"`
	}
	HudsonBatchTask struct {
		Command              string `xml:"command"`
		ConfiguredLocalRules string `xml:"configuredLocalRules,omitempty"`
	}
	HudsonTimeoutTask struct {
		TimeoutMinutes string `xml:"timeoutMinutes"`
	}
	//targets, name, pom, properties, jvmOptions, usePrivateRepository, settings, globalSettings
	HudsonMavenTask struct {
		Targets string `xml:"targets"`
		Name    string `xml:"name"`
		// Pom                  string `xml:"pom"`
		// Properties           string `xml:"properties"`
		// JvmOptions           string `xml:"jvmOptions"`
		UsePrivateRepository bool   `xml:"usePrivateRepository"`
		Settings             string `xml:"settings"`
		GlobalSettings       string `xml:"globalSettings"`
		InjectBuildVariables bool   `xml:"injectBuildVariables"`
	}
	HudsonGradleTask struct {
		Switches                   string `xml:"switches"`
		RootBuildScriptDir         string `xml:"rootBuildScriptDir"`
		WrapperLocation            string `xml:"wrapperLocation"`
		SystemProperties           string `xml:"systemProperties"`
		ProjectProperties          string `xml:"projectProperties"`
		Tasks                      string `xml:"tasks"`
		GradleName                 string `xml:"gradleName"`
		UseWrapper                 bool   `xml:"useWrapper"`
		MakeExecutable             bool   `xml:"makeExecutable"`
		UseWorkspaceAsHome         bool   `xml:"useWorkspaceAsHome"`
		PassAllAsSystemProperties  bool   `xml:"passAllAsSystemProperties"`
		PassAllAsProjectProperties bool   `xml:"passAllAsProjectProperties"`
	}
	ConvertGitHubSetCommitStatusTask struct {
		StatusMessage string `xml:"statusMessage"`
	}

//	<builders>
//
// <com.cloudbees.jenkins.GitHubSetCommitStatusBuilder plugin="github@1.40.0">
// <statusMessage>
// <content/>
// </statusMessage>
// <contextSource class="org.jenkinsci.plugins.github.status.sources.DefaultCommitContextSource"/>
// </com.cloudbees.jenkins.GitHubSetCommitStatusBuilder>
)
