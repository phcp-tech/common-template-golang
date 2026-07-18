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

package dao

import (
	"testing"

	"template/domain/model"
	"template/pkg/dto"

	dbsqlite "github.com/phcp-tech/common-library-golang/dbsqlx/sqlite"
	"github.com/vinovest/sqlx"
)

// newTestDB creates an in-memory SQLite database with a stand-in schema for
// the users table. Real production schema lives in config/schema_sqlite.sql,
// managed independently of this application. — this CREATE TABLE statement exists purely to give the DAO test
// something to run raw SQL against, matching model.User's db tags.
func newTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := dbsqlite.NewSQLite(&dbsqlite.Config{Path: ":memory:"})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`CREATE TABLE users (
		id BIGINT PRIMARY KEY,
		username TEXT,
		nickname TEXT,
		email TEXT,
		kind TEXT,
		status TEXT
	)`)
	if err != nil {
		t.Fatalf("create users table: %v", err)
	}

	return db
}

func seedUser(t *testing.T, db *sqlx.DB, u model.User) {
	t.Helper()
	_, err := db.Exec(db.Rebind(`INSERT INTO users (id, username, nickname, email, kind, status) VALUES (?, ?, ?, ?, ?, ?)`),
		u.Id, u.Username, u.Nickname, u.Email, u.Kind, u.Status)
	if err != nil {
		t.Fatalf("seed user %+v: %v", u, err)
	}
}

func TestUserDao_GetList(t *testing.T) {
	db := newTestDB(t)
	d := NewUserDao(db)

	seedUser(t, db, model.User{Id: 1, Username: "gordon", Nickname: "Gordon", Email: "gordon@example.com", Kind: "Admin", Status: "Active"})
	seedUser(t, db, model.User{Id: 2, Username: "alice", Nickname: "Alice", Email: "alice@example.com", Kind: "Member", Status: "Active"})

	resp, err := d.GetList(&dto.UserListPara{})
	if err != nil {
		t.Fatalf("GetList: %v", err)
	}
	if resp.Total != 2 {
		t.Fatalf("expected total=2, got %d", resp.Total)
	}
	users, ok := resp.List.([]model.User)
	if !ok || len(users) != 2 {
		t.Fatalf("expected 2 users in list, got %+v", resp.List)
	}
	if users[0].Username != "gordon" || users[1].Username != "alice" {
		t.Fatalf("unexpected result: %+v", users)
	}
}

func TestUserDao_GetList_NoRows(t *testing.T) {
	db := newTestDB(t)
	d := NewUserDao(db)

	resp, err := d.GetList(&dto.UserListPara{})
	if err != nil {
		t.Fatalf("GetList: %v", err)
	}
	if resp.Total != 0 {
		t.Fatalf("expected total=0, got %d", resp.Total)
	}
}
