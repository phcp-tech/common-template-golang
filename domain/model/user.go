// Copyright(C) 2020-2026 PHCP Technologies. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

// User domain model. Column sizes/defaults/indexes live directly in the real
// schema (config/schema_sqlite.sql), managed independently of this application.
type User struct {
	Id       int64  `db:"id" json:"id,omitempty"`
	Username string `db:"username" json:"username,omitempty"`
	Nickname string `db:"nickname" json:"nickname,omitempty"`
	Email    string `db:"email" json:"email,omitempty"`
	Kind     string `db:"kind" json:"kind,omitempty"`
	Status   string `db:"status" json:"status,omitempty"`
}
