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

// User domain model
type User struct {
	Id       uint64 `gorm:"primaryKey" json:"id,omitempty"`
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username,omitempty" validate:"omitempty,min=0,max=50"`
	Nickname string `gorm:"size:50" json:"nickname,omitempty" validate:"omitempty,min=0,max=50"`
	Email    string `gorm:"uniqueIndex;size:50;not null" json:"email,omitempty" validate:"omitempty,email,min=0,max=50"`
	Kind     string `gorm:"index;size:40" json:"kind,omitempty" validate:"omitempty,min=0,max=40"`
	Status   int    `gorm:"default:1" json:"status,omitempty" validate:"omitempty,gte=0,lte=3"`
}

// Define User table name
func (h *User) TableName() string {
	return "knw_users.temp_users"
}
