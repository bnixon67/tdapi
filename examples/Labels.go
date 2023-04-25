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
	"strings"

	"github.com/bnixon67/tdapi"
)

func ParseCommandLine() (tokenFile string, scopes []string) {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [options] request\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "options:")
		flag.PrintDefaults()
	}

	flag.StringVar(&tokenFile, "token", ".token.todoist", "path to `file` to use for token")

	var scopeString string
	flag.StringVar(&scopeString,
		"scopes", "data:read", "comma-seperated `scopes` to use for request")

	flag.Parse()

	scopes = strings.Split(scopeString, ",")

	return tokenFile, scopes
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
	tokenFile, scopes := ParseCommandLine()

	// need one remaining arg for request
	if len(flag.Args()) != 0 {
		flag.Usage()
		return
	}

	todoistClient := tdapi.New(tokenFile, clientID, clientSecret, scopes)

	labels, err := todoistClient.GetAllLabels()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("count(labels) = %d\n", len(labels))
	for n, label := range labels {
		fmt.Printf("%d: %+v\n", n, label)
	}

	/*
	label, err := todoistClient.GetLabel(labels[len(labels)-1].ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tdapi.VarToJsonString(label))
	*/

	sharedLabels, err := todoistClient.GetAllSharedLabels()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("count(sharedLabels) = %d\n", len(sharedLabels))
	for n, label := range sharedLabels {
		fmt.Printf("%d: name=%s\n", n, label)
	}

}
