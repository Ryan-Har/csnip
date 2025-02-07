package database

import (
	"database/sql"
	"time"
)

// Helper function to extract a string from sql.NullString
func getString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// Helper function to extract a time.Time from sql.NullTime
func getTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{} // Returns zero time (0001-01-01 00:00:00 UTC)
}

// Helper function to extract an int64 from sql.NullInt64
func getInt64(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

// Helper to convert string to sql.NullString
func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
