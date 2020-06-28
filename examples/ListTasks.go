/*
Copyright 2020 Bill Nixon

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published
by the Free Software Foundation, either version 3 of the License,
or (at your option) any later version.

This program is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
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

	// sort tasks by project order, priority, date, task order
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

	type DisplayTask struct {
		Content  string
		Priority int64
		Labels   []string
		Due      string
	}

	type DisplayProject struct {
		Project string
		Tasks   []DisplayTask
	}

	var displayProjects []DisplayProject
	var lastProject int64 = 0

	// loop thru and build display structures for use in template
	for _, task := range tasks {
		if (labelID == 0 || ContainsInt(task.LabelIds, labelID)) &&
			(projectID == 0 || projectID == task.ProjectID) &&
			ContainsInt(priorities, priority_lookup[task.Priority]) {

			var displayTask DisplayTask

			if lastProject != task.ProjectID {
				displayProjects = append(displayProjects,
					DisplayProject{Project: mapByProjectID[task.ProjectID].Name})
				lastProject = task.ProjectID
			}

			displayTask.Content = task.Content

			displayTask.Priority = priority_lookup[task.Priority]

			if task.Due.String != "" {
				displayTask.Due = task.Due.String
			}

			// loop thru labels, which are sorted, and display matching names
			for _, label := range labels {
				if ContainsInt(task.LabelIds, label.ID) {
					displayTask.Labels = append(displayTask.Labels, label.Name)
				}
			}

			displayProjects[len(displayProjects)-1].Tasks = append(displayProjects[len(displayProjects)-1].Tasks, displayTask)
		}
	}

	const txtTemplate = `
{{- range . -}}
#{{.Project}}
  {{- range .Tasks }}
  {{ .Content }}
  P{{ .Priority }} {{ if .Due }}<{{ .Due }}> {{ end }}{{ range .Labels }}@{{ . }} {{ end }}
  {{ end }}
{{ end -}}
`

	const htmlTemplate = `
{{- $lastProject := "" -}}
<!DOCTYPE html>
<html>
<body>
{{- range . }}
<h1>{{.Project}}</h1>
  <ul>
  {{- range .Tasks }}
  <li>
  {{ .Content }} <em>( Priority {{ .Priority }}, {{ if .Due }}Due {{ .Due }},{{ end }} {{ range .Labels }}@{{ . }} {{ end }}</em>)
  </li>
  {{ end }}
  </ul>
{{ end -}}
</body>
</html>
`

	var t *template.Template

	if html {
		t, err = template.New("output").Parse(htmlTemplate)
	} else {
		t, err = template.New("output").Parse(txtTemplate)
	}
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(os.Stdout, displayProjects)
	if err != nil {
		log.Fatal(err)
	}
}
