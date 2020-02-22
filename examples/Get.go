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
	if len(flag.Args()) != 1 {
		flag.Usage()
		return
	}

	todoistClient := tdapi.New(tokenFile, clientID, clientSecret, scopes)

	resp, err := todoistClient.Get(flag.Arg(0), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(resp))
}
