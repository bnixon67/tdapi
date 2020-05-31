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
	"strings"

	"github.com/bnixon67/tdapi"
)

func ParseCommandLine() (tokenFile string, scopes []string, label string, project string) {
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

	flag.Parse()

	scopes = strings.Split(scopeString, ",")

	return tokenFile, scopes, label, project
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

	var priority_lookup = [...]int{0, 4, 3, 2, 1}

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
	tokenFile, scopes, labelName, projectName := ParseCommandLine()

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

	var lastProject int64 = 0

	// loop thru and display approproiate tasks
	for _, task := range tasks {
		if (labelID == 0 || ContainsInt(task.LabelIds, labelID)) &&
			(projectID == 0 || projectID == task.ProjectID) {

			if lastProject != task.ProjectID {
				fmt.Printf("----- #%s\n\n",
					mapByProjectID[task.ProjectID].Name)
				lastProject = task.ProjectID
			}

			fmt.Printf("%s\n", task.Content)

			fmt.Printf("  ")

			fmt.Printf("P%d ", priority_lookup[task.Priority])

			if task.Due.String != "" {
				fmt.Printf("<%s> ", task.Due.String)
			}

			// loop thru labels, which are sorted, and display matching names
			for _, label := range labels {
				if ContainsInt(task.LabelIds, label.ID) {
					fmt.Printf("@%s ", label.Name)
				}
			}

			fmt.Printf("\n\n")
		}
	}
}
