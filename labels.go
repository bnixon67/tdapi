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
	"strconv"
)

// GetAllLabels returns a JSON-encoded array containing all user labels.
func (c *TodoistClient) GetAllLabels() (response []Label, err error) {
	var body []byte

	body, err = c.Get("/labels", nil)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}

// GetLabel returns a label by ID.
func (c *TodoistClient) GetLabel(id int) (response Label, err error) {
	var body []byte

	body, err = c.Get("/labels/"+strconv.Itoa(id), nil)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}
