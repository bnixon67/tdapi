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
func (c *TodoistClient) GetLabel(id int64) (response Label, err error) {
	var body []byte

	body, err = c.Get("/labels/"+strconv.FormatInt(id, 10), nil)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}
