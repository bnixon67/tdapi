/*
Copyright 2021 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/bnixon67/tdapi"
)

func ParseCommandLine() (tokenFile string, scopes []string, label string,
	project string, priorities []int64, html bool) {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [options] request\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "options:")
		flag.PrintDefaults()
	}

	flag.StringVar(&tokenFile, "token", ".token.todoist", "path to `file` to use for token")

	flag.StringVar(&label, "label", "", "label to filter tasks")

	flag.StringVar(&project, "project", "", "project to filter tasks")

	var scopeString string
	flag.StringVar(&scopeString,
		"scopes", "data:read", "comma-seperated `scopes` to use for request")

	var prioritiesString string
	flag.StringVar(&prioritiesString,
		"priorities", "1,2,3,4", "comma-seperated `priorities` to use for request")

	flag.BoolVar(&html, "html", false, "display in html format")

	flag.Parse()

	scopes = strings.Split(scopeString, ",")

	for _, priority := range strings.Split(prioritiesString, ",") {
		priorityInt, err := strconv.ParseInt(priority, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		priorities = append(priorities, priorityInt)
	}

	return
}

func ContainsInt(slice []int64, want int64) bool {
	for _, value := range slice {
		if value == want {
			return true
		}
	}
	return false
}

func main() {

	var priority_lookup = [...]int64{0, 4, 3, 2, 1}

	// Get Todoist Client ID
	// The ID is not in the source code to avoid someone reusing the ID
	clientID, present := os.LookupEnv("TDCLIENTID")
	if !present {
		log.Fatal("Must set TDCLIENTID")
	}

	// Get Todoist Client Secret
	// The Secret is not in the source code to avoid someone reusing the ID
	clientSecret, present := os.LookupEnv("TDCLIENTSECRET")
	if !present {
		log.Fatal("Must set TDCLIENTSECRET")
	}

	// parse command line to get path to the token file and scopes to use in request
	tokenFile, scopes, labelName, projectName, priorities, html := ParseCommandLine()

	// print usage if invalid command line
	if len(flag.Args()) != 0 {
		flag.Usage()
		return
	}

	// create todoist client
	todoistClient := tdapi.New(tokenFile, clientID, clientSecret, scopes)

	// get all projects
	resp, err := todoistClient.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}

	var projectID int64

	// save projects mapped by project ID
	mapByProjectID := make(map[int64]tdapi.Project)
	for _, project := range resp {
		mapByProjectID[project.ID] = project
		if project.Name == projectName {
			projectID = project.ID
		}
	}

	// Project not found
	if projectName != "" && projectID == 0 {
		fmt.Printf("Project %q not found.\n", projectName)
		return
	}

	// get all labels
	labels, err := todoistClient.GetAllLabels()
	if err != nil {
		log.Fatal(err)
	}

	// store labelID for given label or default to 0 for no label
	var labelID int64

	// save labels mapped by label ID
	mapByLabelID := make(map[int64]tdapi.Label)

	// loop thru all labels
	for _, label := range labels {
		mapByLabelID[label.ID] = label
		if label.Name == labelName {
			labelID = label.ID
		}
	}

	// Label not found
	if labelName != "" && labelID == 0 {
		fmt.Printf("Label %q not found.\n", labelName)
		return
	}

	// sort labels by Order
	sort.Slice(labels, func(i, j int) bool {
		return labels[i].Order < labels[j].Order
	})

	// get all tasks
	tasks, err := todoistClient.GetActiveTasks()
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(tasks, func(i, j int) bool {
		// sort by Project order
		if mapByProjectID[tasks[i].ProjectID].Order < mapByProjectID[tasks[j].ProjectID].Order {
			return true
		}
		if mapByProjectID[tasks[i].ProjectID].Order > mapByProjectID[tasks[j].ProjectID].Order {
			return false
		}

		// sort by Priority, reverse order since p4=1 and p1=4
		if tasks[i].Priority > tasks[j].Priority {
			return true
		}
		if tasks[i].Priority < tasks[j].Priority {
			return false
		}

		// sort by Task Order
		return tasks[i].Order < tasks[j].Order
	})

	/*
	// sort tasks by priority, project order, date, task order
	sort.Slice(tasks, func(i, j int) bool {
		// sort by Priority, reverse order since p4=1 and p1=4
		if tasks[i].Priority > tasks[j].Priority {
			return true
		}
		if tasks[i].Priority < tasks[j].Priority {
			return false
		}

		// sort by Project order
		if mapByProjectID[tasks[i].ProjectID].Order < mapByProjectID[tasks[j].ProjectID].Order {
			return true
		}
		if mapByProjectID[tasks[i].ProjectID].Order > mapByProjectID[tasks[j].ProjectID].Order {
			return false
		}

		// sort by Date
		// copy dates to new variable to default empty date to max value for sorting
		iDate := tasks[i].Due.Date
		if iDate == "" {
			iDate = "9999-99-99"
		}

		jDate := tasks[j].Due.Date
		if jDate == "" {
			jDate = "9999-99-99"
		}

		if iDate < jDate {
			return true
		}
		if iDate > jDate {
			return false
		}

		// sort by Task Order
		return tasks[i].Order < tasks[j].Order
	})
	*/

	type DisplayLabel struct {
		Name     string
		HexColor string
	}

	type DisplayTask struct {
		Project          string
		ProjectHexColor  string
		Content          string
		Priority         int64
		Order            int
		PriorityHexColor string
		Due              string
		Labels           []DisplayLabel
	}

	var displayTasks []DisplayTask

	// loop thru and build display structures for use in template
	for _, task := range tasks {
		if (labelID == 0 || ContainsInt(task.LabelIds, labelID)) && // filter on label (if provided)
			(projectID == 0 || projectID == task.ProjectID) && // filter on project (if provided)
			ContainsInt(priorities, priority_lookup[task.Priority]) { // filter on priorities (if provided)

			var displayTask DisplayTask

			displayTask.Project = mapByProjectID[task.ProjectID].Name
			displayTask.ProjectHexColor = tdapi.ColorToHex[mapByProjectID[task.ProjectID].Color]

			displayTask.Content = task.Content
			displayTask.Order = task.Order

			displayTask.Priority = priority_lookup[task.Priority]
			displayTask.PriorityHexColor = tdapi.PriorityToHexColor[task.Priority]

			if task.Due.String != "" {
				displayTask.Due = task.Due.String
			}

			// loop thru labels, which are sorted, and display matching names
			for _, label := range labels {
				if ContainsInt(task.LabelIds, label.ID) {
					displayTask.Labels = append(displayTask.Labels,
						DisplayLabel{label.Name,
							tdapi.ColorToHex[label.Color]})
				}
			}

			displayTasks = append(displayTasks, displayTask)
		}
	}

	const txtTemplate = `
{{- $lastPriority := 0  -}}
{{- $lastProject  := "" -}}
{{- range . -}}

{{ if (ne $lastProject .Project) -}}
{{ printf "%s" .Project }}{{ $lastProject = .Project }}
{{ end -}}

{{ printf "  %s" .Content }} p{{ .Priority }} {{- if .Due }} <{{ .Due }}>{{ end }} {{- range $n, $v := .Labels }} @{{ .Name }}{{ end }}
{{ end -}}
`
/*
	const txtTemplate = `
{{- $lastPriority := 0  -}}
{{- $lastProject  := "" -}}
{{- range . -}}

{{ if (ne $lastPriority .Priority) -}}
{{ if (ne $lastPriority 0) }}{{ println }}{{ end -}}
Priority {{ .Priority }}{{ $lastPriority = .Priority }}{{ $lastProject = "" }}
{{ end -}}

{{ if (ne $lastProject .Project) -}}
{{ printf "  %s" .Project }}{{ $lastProject = .Project }}
{{ end -}}

{{ printf "    %s" .Content }} {{- if .Due }} <{{ .Due }}>{{ end }} ({{- range $n, $v := .Labels }}{{ if (ne $n 0) }}, {{ end }}{{ .Name }}{{ end }})
{{ end -}}
`
*/

	const htmlTemplate = `
{{- $lastProject := "" -}}
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Task List</title>
</head>
<body style="font-family:sans-serif;">
{{ range . -}}

{{ if (ne $lastProject .Project) -}}
{{ if (ne $lastProject "") }}</ul>{{ end }}
<h2 style="color: {{ .ProjectHexColor }};">{{ .Project }}{{ $lastProject = .Project }}</h2>
<ul>
{{ end -}}

<li>
{{ .Content }}
{{ if .Due }}&lt;<em>{{ .Due }}</em>&gt;{{ end }}
<span style="color:{{.PriorityHexColor}};">P{{ .Priority }}</span>
{{ range $n, $v := .Labels -}}
<span style="color:{{ .HexColor }};">@{{ .Name }}</span>
{{ end -}}
</li>

{{ end -}}
</ul>

</body>
</html>
`
/*
	const htmlTemplate = `
{{- $lastPriority := 0  -}}
{{- $lastProject := "" -}}
<!DOCTYPE html>
<html lang="en">
<head>
<title>Task List</title>
</head>
<body>
{{ range . -}}

{{ if (ne $lastPriority .Priority) -}}
{{ if (ne $lastPriority 0) }}{{ println }}{{ end -}}
{{ if (ne $lastProject "") }}</ul>{{ end }}
<h1 style="color:{{.PriorityHexColor}};">Priority {{ .Priority }}{{ $lastPriority = .Priority }}{{ $lastProject = "" }}</h1>
{{ end -}}

{{ if (ne $lastProject .Project) -}}
{{ if (ne $lastProject "") }}</ul>{{ end }}
<h2 style="color: {{ .ProjectHexColor }};">{{ .Project }}{{ $lastProject = .Project }}</h2>
<ul>
{{ end -}}

<li>
{{ .Content }} {{ if .Due }} &lt;<em>{{ .Due }}</em>&gt;{{ end }}
(<span style="color:{{.PriorityHexColor}};">P{{ .Priority }}</span>,
<span style="color: {{ .ProjectHexColor }};">#{{ .Project }}</span>
{{- range $n, $v := .Labels }} , <span style="color:{{ .HexColor }};">@{{ .Name }}</span>{{ end -}})
</li>
{{ end -}}
</ul>

</body>
</html>
`
*/

	funcMap := template.FuncMap{"Join": strings.Join}

	var t *template.Template

	if html {
		t, err = template.New("output").Funcs(funcMap).Parse(htmlTemplate)
	} else {
		t, err = template.New("output").Funcs(funcMap).Parse(txtTemplate)
	}
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(os.Stdout, displayTasks)
	if err != nil {
		log.Fatal(err)
	}
}
