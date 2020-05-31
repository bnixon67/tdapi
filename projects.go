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
