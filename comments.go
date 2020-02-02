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
	"net/url"
	"strconv"
)

func (c *TodoistClient) GetTaskComments(id int) (response []Comment, err error) {
	var body []byte

	query := url.Values{}
	query.Set("task_id", strconv.Itoa(id))

	body, err = c.Get("/comments", query)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}
