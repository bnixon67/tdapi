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
	"fmt"
)

// A Project represents a Todoist Project.
// See https://developer.todoist.com/rest/v2/?shell#projects for more details.
type Project struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Color          string  `json:"color"`
	ParentID       *string `json:"parent_id"`
	Order          int     `json:"order"`
	CommentCount   int     `json:"comment_count"`
	IsShared       bool    `json:"is_shared"`
	IsFavorite     bool    `json:"is_favorite"`
	IsInboxProject bool    `json:"is_inbox_project"`
	IsTeamInbox    bool    `json:"is_team_inbox"`
	ViewStyle      string  `json:"view_style"`
	URL            string  `json:"url"`
}

// GetProjects returns all user projects.
func (c *TodoistClient) GetAllProjects() ([]Project, error) {
	body, err := c.Get("/projects", nil)
	if err != nil {
		return nil, err
	}

	var projects []Project
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// GetProject returns the project for the given project_id.
func (c *TodoistClient) GetProject(project_id string) (Project, error) {
	if project_id == "" {
		return Project{}, fmt.Errorf("empty project ID")
	}

	body, err := c.Get("/projects/"+project_id, nil)
	if err != nil {
		return Project{}, err
	}

	var project Project
	err = json.Unmarshal(body, &project)

	return project, err
}
