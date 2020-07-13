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

import "time"

type Project struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Parent       int64  `json:"parent"`
	Order        int    `json:"order"`
	Color        int    `json:"color"`
	CommentCount int    `json:"comment_count"`
}

type Task struct {
	// Task id.
	ID int64 `json:"id"`

	// Task’s project id (read-only).
	ProjectID int64 `json:"project_id"`

	// ID of section task belongs to.
	SectionID int64 `json:"section_id"`

	// Task content.
	Content string `json:"content"`

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

		// Date in format YYYY-MM-DD corrected to user’s timezone.
		Date string `json:"date"`

		// Only returned if exact due time set (i.e. it’s not a whole-day task),
		// date and time in RFC3339 format in UTC.
		Datetime *time.Time `json:"datetime,omitempty"`

		// Only returned if exact due time set, user’s timezone definition either in
		// tzdata-compatible format (“Europe/Berlin”) or as a string specifying east
		// of UTC offset as “UTC±HH:MM” (i.e. “UTC-01:00”).
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

type Label struct {
	// The id of the label.
	ID int64 `json:"id"`

	// The name of the label.
	Name string `json:"name"`

	// Number used by clients to sort list of labels.
	Order int `json:"order"`

	// Color id. It’s a value between 30 and 49, refer to Colors for more info.
	Color int `json:"color"`
}

// ColorToHex maps color ids for a project, label, or filter, to a hex value.
var ColorToHex = [...]string{
	30: "#b8256f",
	31: "#db4035",
	32: "#ff9933",
	33: "#fad000",
	34: "#afb83b",
	35: "#7ecc49",
	36: "#299438",
	37: "#6accbc",
	38: "#158fad",
	39: "#14aaf5",
	40: "#96c3eb",
	41: "#4073ff",
	42: "#884dff",
	43: "#af38eb",
	44: "#eb96eb",
	45: "#e05194",
	46: "#ff8d85",
	47: "#808080",
	48: "#b8b8b8",
	49: "#ccac93",
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

	// Comment’s task id (for task comments).
	TaskID int64 `json:"task_id,omitempty"`

	// Comment’s project id (for project comments).
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
		// representing an HTTP URL). Note that we don’t cache
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
