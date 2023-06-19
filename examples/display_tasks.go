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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/bnixon67/tdapi"
)

// getCredentials retrieves the Todoist ID and Secret from
// the corresponding environment variables TDCLIENTID and TDCLIENTSECRET.
// If the required environment variables are not present, an error is
// returned.
func getCredentials() (id string, secret string, err error) {
	id, ok := os.LookupEnv("TDCLIENTID")
	if !ok {
		return "", "", fmt.Errorf("TDCLIENTID not set.")
	}

	secret, ok = os.LookupEnv("TDCLIENTSECRET")
	if !ok {
		return "", "", fmt.Errorf("TDCLIENTSECRET not set.")
	}

	return id, secret, nil
}

// setupFlags configures the command-line flags for the application.
// tokenFile and scopesString will store the corresponding flag values.
func setupFlags(tokenFile *string, scopesString *string) {
	flag.Usage = func() {
		output := flag.CommandLine.Output()

		fmt.Fprintf(output, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(output, "Options:")
		flag.PrintDefaults()
	}

	flag.StringVar(tokenFile,
		"token", ".token.todoist", "Path to the token file")
	flag.StringVar(scopesString,
		"scopes", "data:read", "Comma-separated scopes for the request")
}

func main() {
	var (
		tokenFile    string
		scopesString string
	)

	clientID, clientSecret, err := getCredentials()
	if err != nil {
		fmt.Println("Cannot get credentials:", err)
		os.Exit(1)
	}

	setupFlags(&tokenFile, &scopesString)
	flag.Parse()

	// usage error is there are remaining arguments
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}

	scopes := strings.Split(scopesString, ",")
	todoistClient := tdapi.New(tokenFile, clientID, clientSecret, scopes)

	projects, err := todoistClient.GetAllProjects()
	if err != nil {
		fmt.Println("Error retrieving projects:", err)
		os.Exit(3)
	}
	fmt.Printf("Projects count: %d\n", len(projects))

	/*
		// projectByID is a map to allow the lookup of a project by ID
		projectByID := make(map[string]tdapi.Project, len(projects))
		for _, project := range projects {
			projectByID[project.ID] = project
		}
	*/
	projectByID := tdapi.ProjectByID(projects)

	var p tdapi.TaskParameters

	// p.ProjectID = "2306104483"
	// p.IDs = []int64{6952898337, 6952626760}
	// p.Label = "CASA"
	p.Filter = "@Johnette_Shepek | assigned to: johnette.shepek"

	tasks, err := todoistClient.GetActiveTasks(&p)
	if err != nil {
		fmt.Println("Error retrieving tasks:", err)
		os.Exit(3)
	}

	fmt.Printf("Tasks count: %d\n", len(tasks))

	groupByProjectID := tdapi.GroupTasksByProjectID(tasks)
	for projectID, tasks := range groupByProjectID {
		fmt.Println(projectByID[projectID].Name, len(tasks))
	}

	childProjectIDs := tdapi.ChildProjectIDs(projects)

	var printRecursive func(projectID string, indent string)
	printRecursive = func(projectID string, indent string) {
		tasks := groupByProjectID[projectID]
		if len(tasks) > 0 {
			fmt.Println(
				indent,
				projectByID[projectID].Name,
				len(tasks),
			)

			sort.Sort(TasksByPriorityThenOrder(tasks))

			for _, task := range tasks {
				var dueDate string
				if task.Due != nil {
					dueDate = task.Due.Date
				}
				fmt.Printf("%s %.*s %d %d %s\n",
					indent+"  ",
					20,
					task.Content,
					task.Priority,
					task.Order,
					dueDate,
				)
			}
		}

		for _, childID := range childProjectIDs[projectID] {
			printRecursive(childID, indent+"-")
		}
	}

	fmt.Println("-----")
	printRecursive("", "")
}

func prettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

type TasksByPriorityThenOrder []tdapi.Task

func (tasks TasksByPriorityThenOrder) Len() int {
	return len(tasks)
}

func (tasks TasksByPriorityThenOrder) Swap(i, j int) {
	tasks[i], tasks[j] = tasks[j], tasks[i]
}

func (tasks TasksByPriorityThenOrder) Less(i, j int) bool {
	if tasks[i].Priority != tasks[j].Priority {
		// descending Priority since p1 = 4 and p4 = 1
		return tasks[i].Priority > tasks[j].Priority
	}
	// ascending Order if Priority is the same
	return tasks[i].Order < tasks[j].Order
}
