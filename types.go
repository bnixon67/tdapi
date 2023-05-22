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

import "time"

type Task struct {
	// Task id.
	ID string `json:"id"`

	// Project ID of the task
	ProjectID string `json:"project_id"`

	// ID of section task belongs to.
	SectionID int64 `json:"section_id"`

	// Task content.
	Content string `json:"content"`

	// Description for the task.
	Description string `json:"description,omitempty"`

	// Flag to mark completed tasks.
	Completed bool `json:"completed"`

	// Array of label ids, associated with a task.
	LabelIds []int64 `json:"label_ids"`

	// ID of parent task (read-only, absent for top-level tasks).
	Parent int64 `json:"parent,omitempty"`

	// Position under the same parent or project for top-level tasks (read-only).
	Order int `json:"order"`

	// Task priority from 1 (normal, default value) to 4 (urgent).
	Priority int `json:"priority"`

	// Task due date/time.
	Due struct {
		// Human defined date in arbitrary format.
		String string `json:"string"`

		// Date in format YYYY-MM-DD corrected to userÃ¢ÂÂs timezone.
		Date string `json:"date"`

		// Only returned if exact due time set (i.e. itÃ¢ÂÂs not a whole-day task),
		// date and time in RFC3339 format in UTC.
		Datetime *time.Time `json:"datetime,omitempty"`

		// Only returned if exact due time set, userÃ¢ÂÂs timezone definition either in
		// tzdata-compatible format (Ã¢ÂÂEurope/BerlinÃ¢ÂÂ) or as a string specifying east
		// of UTC offset as Ã¢ÂÂUTCÃÂ±HH:MMÃ¢ÂÂ (i.e. Ã¢ÂÂUTC-01:00Ã¢ÂÂ).
		Timezone string `json:"timezone,omitempty"`

		// Flag to indicate recurring.
		Recurring bool `json:"recurring"`
	} `json:"due"`

	// URL to access this task in Todoist web interface.
	URL string `json:"url,omitempty"`

	CommentCount int `json:"comment_count"`

	// DateTime when Task was created.
	Created time.Time `json:"created"`
}

// ColorToHex maps color ids for a project, label, or filter, to a hex value.
var ColorToHex = map[string]string{
	"berry_red":   "#b8256f",
	"red":         "#db4035",
	"orange":      "#ff9933",
	"yellow":      "#fad000",
	"olive_green": "#afb83b",
	"lime_green":  "#7ecc49",
	"green":       "#299438",
	"mint_green":  "#6accbc",
	"teal":        "#158fad",
	"sky_blue":    "#14aaf5",
	"light_blue":  "#96c3eb",
	"blue":        "#4073ff",
	"grape":       "#884dff",
	"violet":      "#af38eb",
	"lavendar":    "#eb96eb",
	"magenta":     "#e05194",
	"salmon":      "#ff8d85",
	"charcoal":    "#808080",
	"grey":        "#b8b8b8",
	"taupe":       "#ccac93",
}

// PriorityToHexColor maps task priority (1=p4, ..., 4=p1) to a hex value.
var PriorityToHexColor = [...]string{
	1: "#000000",
	2: "#246fe0",
	3: "#eb8909",
	4: "#d1453b",
}

type Comment struct {
	// Comment id.
	ID int64 `json:"id"`

	// CommentÃ¢ÂÂs task id (for task comments).
	TaskID int64 `json:"task_id,omitempty"`

	// CommentÃ¢ÂÂs project id (for project comments).
	ProjectID int64 `json:"project_id,omitempty"`

	// Date and time when comment was added, RFC3339 format in UTC.
	Posted time.Time `json:"posted"`

	// Comment content.
	Content string `json:"content"`

	// Attachment file (optional).
	Attachment struct {
		// The name of the file.
		FileName string `json:"file_name,omitempty"`

		// MIME type (i.e. text/plain, image/png).
		FileType string `json:"file_type,omitempty"`

		// The URL where the file is located (a string value
		// representing an HTTP URL). Note that we donÃ¢ÂÂt cache
		// the remote content on our servers and stream or expose
		// files directly from third party resources. In particular
		// this means that you should avoid providing links to
		// non-encrypted (plain HTTP) resources, as exposing this
		// files in Todoist may issue a browser warning.
		FileURL string `json:"file_url,omitempty"`

		FileSize    int64  `json:"file_size,omitempty"`
		UploadState string `json:"upload_state,omitempty"`

		ResourceType string `json:"resource_type"`

		SiteName    string `json:"site_name,omitempty"`
		Description string `json:"description,omitempty"`
		Title       string `json:"title,omitempty"`
		URL         string `json:"url,omitempty"`
		FavIcon     string `json:"favicon,omitempty"`
	} `json:"attachment"`
}
