// Code generated by ent, DO NOT EDIT.

package history

import (
	"time"
)

const (
	// Label holds the string label denoting the history type in the database.
	Label = "history"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldEvent holds the string denoting the event field in the database.
	FieldEvent = "event"
	// FieldIP holds the string denoting the ip field in the database.
	FieldIP = "ip"
	// FieldUserAgent holds the string denoting the user_agent field in the database.
	FieldUserAgent = "user_agent"
	// Table holds the table name of the history in the database.
	Table = "histories"
)

// Columns holds all SQL columns for history fields.
var Columns = []string{
	FieldID,
	FieldUserID,
	FieldCreatedAt,
	FieldEvent,
	FieldIP,
	FieldUserAgent,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
)
