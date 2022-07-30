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

// Options represents possbile command line options
type Options struct {
	TokenFile   string
	Scopes      []string
	LabelName   string
	ProjectName string
	Priorities  []int64
	Html        bool
}

// DisplayLabel represents a Label for display
type DisplayLabel struct {
	Name     string
	HexColor string
}

// DisplayTask represents a Task for display
type DisplayTask struct {
	Project          string
	ProjectHexColor  string
	Content          string
	Description      string
	Priority         int64
	Order            int
	PriorityHexColor string
	Due              string
	Labels           []DisplayLabel
}

// ParseCommandLine parses the command line returning the options provided or default value.
func ParseCommandLine() Options {
	var opt Options
	var scopeString string
	var prioritiesString string

	// define usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "options:")
		flag.PrintDefaults()
	}

	// define name, default value, and usage for flag values and bind to variable
	flag.StringVar(&opt.TokenFile, "token", ".token.todoist", "path to `file` to use for token")
	flag.StringVar(&opt.LabelName, "label", "", "label to filter tasks")
	flag.StringVar(&opt.ProjectName, "project", "", "project to filter tasks")
	flag.StringVar(&scopeString, "scopes", "data:read", "comma-seperated `scopes` to use for request")
	flag.StringVar(&prioritiesString, "priorities", "1,2,3,4", "comma-seperated `priorities` to use for request")
	flag.BoolVar(&opt.Html, "html", false, "display in html format")

	flag.Parse()

	// convert string to slice for scopes
	opt.Scopes = strings.Split(scopeString, ",")

	// convert string to slice for priorities
	for _, priority := range strings.Split(prioritiesString, ",") {
		priorityInt, err := strconv.ParseInt(priority, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		opt.Priorities = append(opt.Priorities, priorityInt)
	}

	return opt
}

// ContainsInt64 reports whether v is present in s.
func ContainsInt64(s []int64, v int64) bool {
	for _, value := range s {
		if value == v {
			return true
		}
	}
	return false
}

func main() {
	// task priority is stored as an intger from 1 (normal, default value) to 4 (urgent).
	// priorityValue maps priority to a value, with urgent as 1.
	var priorityValue = [...]int64{0, 4, 3, 2, 1}

	// get Todoist Client ID from env to avoid storing in source code
	clientID, present := os.LookupEnv("TDCLIENTID")
	if !present {
		log.Fatal("Must set TDCLIENTID")
	}

	// get Todoist Client Secret from env to avoid storing in source code
	clientSecret, present := os.LookupEnv("TDCLIENTSECRET")
	if !present {
		log.Fatal("Must set TDCLIENTSECRET")
	}

	// parse command line to get program options
	opt := ParseCommandLine()

	// print usage if invalid command line
	if len(flag.Args()) != 0 {
		flag.Usage()
		return
	}

	// create todoist client
	todoistClient := tdapi.New(opt.TokenFile, clientID, clientSecret, opt.Scopes)

	// get all projects
	projects, err := todoistClient.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}

	// projectID contains the id of the project to filter, or default of 0 if not provided
	var projectID int64

	// projectByID is a map to allow the lookup of a project by ID
	projectByID := make(map[int64]tdapi.Project)
	for _, project := range projects {
		projectByID[project.ID] = project
		if project.Name == opt.ProjectName {
			projectID = project.ID
		}
	}

	// if project name was supplied, check if it exists
	if opt.ProjectName != "" && projectID == 0 {
		fmt.Printf("Project %q not found.\n", opt.ProjectName)
		return
	}

	// get all labels
	labels, err := todoistClient.GetAllLabels()
	if err != nil {
		log.Fatal(err)
	}

	// labelID contains the id of the label to filter, or default of 0 if not provided
	var labelID int64

	// labelByID is a map to allow the lookup of a label by ID
	labelByID := make(map[int64]tdapi.Label)
	for _, label := range labels {
		labelByID[label.ID] = label
		if label.Name == opt.LabelName {
			labelID = label.ID
		}
	}

	// if label name was supplied, check if it exists
	if opt.LabelName != "" && labelID == 0 {
		fmt.Printf("Label %q not found.\n", opt.LabelName)
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

	// sort tasks by Project, Priority, Due Date, Task Order
	sort.Slice(tasks, func(i, j int) bool {
		// sort by Project order
		if projectByID[tasks[i].ProjectID].Order < projectByID[tasks[j].ProjectID].Order {
			return true
		}
		if projectByID[tasks[i].ProjectID].Order > projectByID[tasks[j].ProjectID].Order {
			return false
		}

		// sort by Priority, reverse order since p4=1 and p1=4
		if tasks[i].Priority > tasks[j].Priority {
			return true
		}
		if tasks[i].Priority < tasks[j].Priority {
			return false
		}

		// sort by Due Date
		iDate := tasks[i].Due.Date
		jDate := tasks[j].Due.Date
		if iDate == "" {
			iDate = "9999-99-99"
		}
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

	// build display structures for use in output template
	var displayTasks []DisplayTask
	for _, task := range tasks {
		// filter tasks
		if (labelID == 0 || ContainsInt64(task.LabelIds, labelID)) &&
			(projectID == 0 || projectID == task.ProjectID) &&
			ContainsInt64(opt.Priorities, priorityValue[task.Priority]) {

			var displayTask DisplayTask

			displayTask.Project = projectByID[task.ProjectID].Name
			displayTask.ProjectHexColor = tdapi.ColorToHex[projectByID[task.ProjectID].Color]

			displayTask.Content = task.Content
			displayTask.Description = task.Description
			displayTask.Order = task.Order

			displayTask.Priority = priorityValue[task.Priority]
			displayTask.PriorityHexColor = tdapi.PriorityToHexColor[task.Priority]

			if task.Due.String != "" {
				displayTask.Due = task.Due.String
			}

			// loop thru labels, which are sorted, and display matching names
			for _, label := range labels {
				if ContainsInt64(task.LabelIds, label.ID) {
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

{{ printf "  %s" .Content }} p{{ .Priority }} {{- if .Due }} <{{ .Due }}>{{ end }} {{- range $n, $v := .Labels }} @{{ .Name }}{{ end }} {{ if .Description }} - {{ .Description }}{{ end }}
{{ end -}}
`

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
{{ range $n, $v := .Labels -}}<span style="color:{{ .HexColor }};">@{{ .Name }}</span> {{ end -}}
{{ if .Description }}<div style="font-size: 90%;">{{ .Description }}</div>{{ end }}
</li>

{{ end -}}
</ul>

</body>
</html>
`

	funcMap := template.FuncMap{"Join": strings.Join}

	var t *template.Template

	if opt.Html {
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
