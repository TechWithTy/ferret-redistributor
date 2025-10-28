package db

import (
    "database/sql"
    "errors"
    "os"
    "strconv"
    "time"
)

// Open creates a *sql.DB with sane pool defaults. Caller must import a driver.
func Open(dsn string) (*sql.DB, error) {
    if dsn == "" {
        return nil, errors.New("db: empty DSN")
    }
    d, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    // Pool tuning via env, with sensible defaults
    d.SetMaxOpenConns(intFromEnv("DB_MAX_OPEN_CONNS", 10))
    d.SetMaxIdleConns(intFromEnv("DB_MAX_IDLE_CONNS", 5))
    d.SetConnMaxLifetime(durationFromEnv("DB_CONN_MAX_LIFETIME", 30*time.Minute))
    // Validate connection
    if err := withTimeoutPing(d, 5*time.Second); err != nil {
        _ = d.Close()
        return nil, err
    }
    return d, nil
}

// OpenFromEnv reads DATABASE_URL and opens a DB.
func OpenFromEnv() (*sql.DB, error) {
    return Open(os.Getenv("DATABASE_URL"))
}

func withTimeoutPing(d *sql.DB, timeout time.Duration) error {
    if d == nil { return errors.New("db: nil DB") }
    c := make(chan error, 1)
    go func() { c <- d.Ping() }()
    select {
    case err := <-c:
        return err
    case <-time.After(timeout):
        return errors.New("db: ping timeout")
    }
}

func intFromEnv(key string, def int) int {
    if v := os.Getenv(key); v != "" {
        if n, err := strconv.Atoi(v); err == nil {
            return n
        }
    }
    return def
}

func durationFromEnv(key string, def time.Duration) time.Duration {
    if v := os.Getenv(key); v != "" {
        if d, err := time.ParseDuration(v); err == nil {
            return d
        }
    }
    return def
}

