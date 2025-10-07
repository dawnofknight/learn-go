package main

// GetDSN returns the database DSN from environment or a sensible default.
// Override by setting `DB_DSN` in your environment.
const defaultDSN = "root:root@tcp(127.0.0.1:3306)/testdb?parseTime=true&charset=utf8mb4&loc=Local"

func GetDSN() string {
	return env("DB_DSN", defaultDSN)
}
