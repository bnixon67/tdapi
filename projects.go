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

package tdapi

import (
	"encoding/json"
	"errors"
	"strconv"
)

// GetAllProjects returns JSON-encoded array containing all user projects.
func (c *TodoistClient) GetAllProjects() (response []Project, err error) {
	var body []byte

	body, err = c.Get("/projects", nil)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}

// GetProject returns a JSON object containing a project object related to the given id.
func (c *TodoistClient) GetProject(id int64) (response Project, err error) {
	var body []byte

	// project id 0 returns array of projects
	if id == 0 {
		return response, errors.New("Invalid project id 0")
	}

	body, err = c.Get("/projects/"+strconv.FormatInt(id, 10), nil)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}
