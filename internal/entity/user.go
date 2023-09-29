package entity

import "time"

type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

type File struct {
	Name string
	Size int64
	URL  string
}
