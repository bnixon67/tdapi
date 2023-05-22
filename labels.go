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
)

// A PersonalLabel represents a Todoist personal label.
// See https://developer.todoist.com/rest/v2/?shell#labels for more details.
type PersonalLabel struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Order      int    `json:"order"`
	IsFavorite bool   `json:"is_favorite"`
}

// GetAllLabels returns a JSON-encoded array containing all user labels.
func (c *TodoistClient) GetAllPersonalLabels() (response []PersonalLabel, err error) {
	var body []byte

	body, err = c.Get("/labels", nil)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}

// GetAllSharedLabels returns a JSON-encoded array containing all shared labels.
func (c *TodoistClient) GetAllSharedLabels() (response []string, err error) {
	var body []byte

	body, err = c.Get("/labels/shared", nil)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}

// GetPersonalLabel returns a label by ID.
func (c *TodoistClient) GetPersonalLabel(id string) (response PersonalLabel, err error) {
	var body []byte

	body, err = c.Get("/labels/"+id, nil)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}
