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

package todoist

import "time"

type Project struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Parent       int    `json:"parent"`
	Order        int    `json:"order"`
	CommentCount int    `json:"comment_count"`
}

type Task struct {
	// Task id.
	ID int `json:"id"`

	// Task’s project id (read-only).
	ProjectID int `json:"project_id"`

	// ID of section task belongs to.
	SectionID int `json:"section_id"`

	// Task content.
	Content string `json:"content"`

	// Flag to mark completed tasks.
	Completed bool `json:"completed"`

	// Array of label ids, associated with a task.
	LabelIds []int `json:"label_ids"`

	// ID of parent task (read-only, absent for top-level tasks).
	Parent int `json:"parent,omitempty"`

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
