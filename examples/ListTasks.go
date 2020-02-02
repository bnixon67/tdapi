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

func ParseCommandLine() (tokenFile string, scopes []string, label string) {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [options] request\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "options:")
		flag.PrintDefaults()
	}

	flag.StringVar(&tokenFile, "token", ".token.todoist", "path to `file` to use for token")

	flag.StringVar(&label, "label", "", "label to filter tasks")

	var scopeString string
	flag.StringVar(&scopeString,
		"scopes", "data:read", "comma-seperated `scopes` to use for request")

	flag.Parse()

	scopes = strings.Split(scopeString, ",")

	return tokenFile, scopes, label
}

func ContainsInt(a []int, x int) bool {
	ret := false

	for _, n := range a {
		if x == n {
			ret = true
		}
	}

	return (ret)
}

func main() {
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
	tokenFile, scopes, labelStr := ParseCommandLine()

	if len(flag.Args()) != 0 {
		flag.Usage()
		return
	}

	todoistClient := tdapi.New(tokenFile, clientID, clientSecret, scopes)

	projectIntToString := make(map[int]string)

	resp, err := todoistClient.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}

	for _, project := range resp {
		projectIntToString[project.ID] = project.Name
	}

	var labelID int

	labelIntToString := make(map[int]string)
	labelStringToInt := make(map[string]int)

	labels, err := todoistClient.GetAllLabels()
	if err != nil {
		log.Fatal(err)
	}

	for _, label := range labels {
		labelIntToString[label.ID] = label.Name
		labelStringToInt[label.Name] = label.ID

		if label.Name == labelStr {
			labelID = label.ID
		}
	}

	if labelStr != "" && labelID == 0 {
		// Label not found
		fmt.Printf("Label %s not found.\n", labelStr)
	}

	tasks, err := todoistClient.GetActiveTasks()
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Priority > tasks[j].Priority {
			return true
		}
		if tasks[i].Priority < tasks[j].Priority {
			return false
		}

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

		return tasks[i].Order < tasks[j].Order
	})

	for _, task := range tasks {
		if ContainsInt(task.LabelIds, labelID) || (labelID == 0) {
			fmt.Printf("%s\n", task.Content)
			for _, x := range task.LabelIds {
				fmt.Printf("@%s ", labelIntToString[x])
			}
			fmt.Printf("#%s %s\n", projectIntToString[task.ProjectID], task.Due.String)
			fmt.Println()
		}
	}
}
