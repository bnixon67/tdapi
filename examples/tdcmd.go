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
	"os"
	"strings"

	"github.com/bnixon67/tdapi"
)

// getCredentials retrieves the Todoist client ID and secret from the corresponding environment variables TDCLIENTID and TDCLIENTSECRET.
// If the required environment variables are not present, an error is returned.
func getCredentials() (string, string, error) {
	clientID, present := os.LookupEnv("TDCLIENTID")
	if !present {
		return "", "", fmt.Errorf("TDCLIENTID environment variable is not set")
	}

	clientSecret, present := os.LookupEnv("TDCLIENTSECRET")
	if !present {
		return "", "", fmt.Errorf("TDCLIENTSECRET environment variable is not set")
	}

	return clientID, clientSecret, nil
}

// setupFlags configures the command-line flags for the application.
// It takes two pointer parameters, tokenFile and scopesString, which will be used to store the corresponding flag values.
func setupFlags(tokenFile *string, scopesString *string) {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <command> [arguments]\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "Command:")
		fmt.Fprintln(flag.CommandLine.Output(), "  list")
		fmt.Fprintln(flag.CommandLine.Output(), "    projects")
		fmt.Fprintln(flag.CommandLine.Output(), "    shared_labels")
		fmt.Fprintln(flag.CommandLine.Output(), "    personal_labels")
	}

	flag.StringVar(tokenFile, "token", ".token.todoist", "Path to the token file")
	flag.StringVar(scopesString, "scopes", "data:read", "Comma-separated scopes for the request")
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

	/*
		// usage error is there are remaining arguments
		if flag.NArg() != 0 {
			flag.Usage()
			os.Exit(2)
		}
	*/

	scopes := strings.Split(scopesString, ",")
	todoistClient := tdapi.New(tokenFile, clientID, clientSecret, scopes)

	command := flag.Arg(0)
	subCommand := flag.Arg(1)

	fmt.Printf("command %q subCommand %q\n", command, subCommand)

	switch command {
	case "list":
		switch subCommand {
		case "projects":
			runGetAllProjects(todoistClient)
		case "shared_labels":
			runGetAllSharedLabels(todoistClient)
		case "personal_labels":
			runGetAllPersonalLabels(todoistClient)
		default:
			flag.Usage()
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func runGetAllProjects(todoistClient *tdapi.TodoistClient) {
	projects, err := todoistClient.GetAllProjects()
	if err != nil {
		fmt.Println("Error retrieving projects:", err)
		os.Exit(3)
	}

	fmt.Printf("Projects count: %d\n", len(projects))
	for index, project := range projects {
		projectJSON, err := json.MarshalIndent(project, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling project JSON:", err)
			return
		}
		fmt.Printf("Project %d:\n%s\n", index, projectJSON)
	}
}

func runGetAllSharedLabels(todoistClient *tdapi.TodoistClient) {
	labels, err := todoistClient.GetAllSharedLabels()
	if err != nil {
		fmt.Println("Error retrieving shared labels:", err)
		os.Exit(3)
	}

	fmt.Printf("Shared labels count: %d\n", len(labels))
	for index, label := range labels {
		fmt.Printf("Label %d: %s\n", index, label)
	}
}

func runGetAllPersonalLabels(todoistClient *tdapi.TodoistClient) {
	labels, err := todoistClient.GetAllPersonalLabels()
	if err != nil {
		fmt.Println("Error retrieving personal labels:", err)
		os.Exit(3)
	}

	fmt.Printf("Shared labels count: %d\n", len(labels))
	for index, label := range labels {
		labelJSON, err := json.MarshalIndent(label, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling label JSON:", err)
			return
		}
		fmt.Printf("Label %d:\n%s\n", index, labelJSON)
	}
}
