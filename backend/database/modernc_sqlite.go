package database

// Import the pure-Go modernc sqlite driver to avoid cgo dependency (replaces mattn/go-sqlite3).
// The driver registers itself with database/sql on import.
import _ "modernc.org/sqlite"
