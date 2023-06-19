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
	"net/url"
	"strconv"
	"strings"
	"time"
)

// See https://developer.todoist.com/rest/v2/?shell#tasks
type TaskDue struct {
	String      string     `json:"string"`
	Date        string     `json:"date"`
	IsRecurring bool       `json:"is_recurring"`
	Datetime    *time.Time `json:"datetime,omitempty"`
	Timezone    string     `json:"timezone,omitempty"`
}

// See https://developer.todoist.com/rest/v2/?shell#tasks
type Task struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"project_id"`
	SectionID    string    `json:"section_id,omitempty"`
	Content      string    `json:"content"`
	Description  string    `json:"description,omitempty"`
	IsCompleted  bool      `json:"is_completed"`
	Labels       []string  `json:"labels,omitempty"`
	ParentID     string    `json:"parent,omitempty"`
	Order        int       `json:"order"`
	Priority     int       `json:"priority"`
	Due          *TaskDue  `json:"due,omitempty"`
	URL          string    `json:"url,omitempty"`
	CommentCount int       `json:"comment_count,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	CreatorID    string    `json:"creator_id"`
	AssigneeID   string    `json:"assignee_id,omitempty"`
	AssignerID   string    `json:"assigner_id,omitempty"`
	Duration     int       `json:"duration"`
	DurationUnit string    `json:"duration_unit"`
}

// See https://developer.todoist.com/rest/v2/?shell#get-active-tasks
type TaskParameters struct {
	ProjectID string
	SectionID string
	Label     string
	Filter    string
	Lang      string
	IDs       []int64
}

func joinInt64Slice(slice []int64, separator string) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(strSlice, separator)
}

// GetActiveTasks returns an array containing all active tasks.
func (c *TodoistClient) GetActiveTasks(p *TaskParameters) (response []Task, err error) {
	var body []byte

	query := url.Values{}

	if p != nil {
		if p.ProjectID != "" {
			query.Set("project_id", p.ProjectID)
		}

		// TODO: section_id

		if len(p.Label) > 0 {
			query.Set("label", p.Label)
		}

		if len(p.Filter) > 0 {
			query.Set("filter", p.Filter)
		}

		// TODO: lang

		if len(p.IDs) > 0 {
			query.Set("ids", joinInt64Slice(p.IDs, ","))
		}
	}

	body, err = c.Get("/tasks", query)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}

// GetActiveTask returns an active (non-completed) task by id.
func (c *TodoistClient) GetActiveTask(id int) (response Task, err error) {
	var body []byte

	body, err = c.Get("/tasks/"+strconv.Itoa(id), nil)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}

// GroupTasksByProjectID groups Tasks by their Project ID.
func GroupTasksByProjectID(tasks []Task) map[string][]Task {
	groups := make(map[string][]Task)
	for _, task := range tasks {
		groups[task.ProjectID] = append(groups[task.ProjectID], task)
	}
	return groups
}
