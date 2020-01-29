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

// APIErrorResponse contains a single property named error.
type APIErrorResponse struct {
	Err string
}

// Error return a string representation of the error
func (e *APIErrorResponse) Error() string {
	return e.Err
}

func codeIsError(code int) bool {
	if code >= 400 && code <= 599 {
		return true
	}
	return false
}
