package cache

import (
    "bufio"
    "errors"
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
    "time"
)

// Valkey is a tiny RESP2 client sufficient for basic caching against Valkey/Redis.
// It supports AUTH, SELECT, PING, GET, SET (with EX), and DEL.
type Valkey struct {
    conn net.Conn
    rw   *bufio.ReadWriter
}

type ValkeyConfig struct {
    Addr     string // host:port
    Password string
    DB       int
    Timeout  time.Duration
}

// NewValkey opens a TCP connection and performs optional AUTH/SELECT.
func NewValkey(cfg ValkeyConfig) (*Valkey, error) {
    if cfg.Addr == "" {
        cfg.Addr = getenv("VALKEY_ADDR", "127.0.0.1:6379")
    }
    if cfg.Timeout <= 0 {
        cfg.Timeout = 5 * time.Second
    }
    c, err := net.DialTimeout("tcp", cfg.Addr, cfg.Timeout)
    if err != nil {
        return nil, err
    }
    v := &Valkey{conn: c, rw: bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))}
    if pwd := firstNonEmpty(cfg.Password, os.Getenv("VALKEY_PASSWORD")); pwd != "" {
        if err := v.send("AUTH", pwd); err != nil { _ = v.Close(); return nil, err }
        if _, err := v.readSimpleOK(); err != nil { _ = v.Close(); return nil, err }
    }
    db := cfg.DB
    if db == 0 {
        if s := os.Getenv("VALKEY_DB"); s != "" {
            if n, err := strconv.Atoi(s); err == nil { db = n }
        }
    }
    if db > 0 {
        if err := v.send("SELECT", strconv.Itoa(db)); err != nil { _ = v.Close(); return nil, err }
        if _, err := v.readSimpleOK(); err != nil { _ = v.Close(); return nil, err }
    }
    return v, nil
}

func (v *Valkey) Close() error {
    if v == nil || v.conn == nil { return nil }
    return v.conn.Close()
}

// Ping issues PING and expects PONG.
func (v *Valkey) Ping() error {
    if err := v.send("PING"); err != nil { return err }
    s, err := v.readSimpleString()
    if err != nil { return err }
    if strings.ToUpper(s) != "PONG" { return fmt.Errorf("valkey: unexpected PING reply: %q", s) }
    return nil
}

// Get returns the value for key, or ok=false on nil.
func (v *Valkey) Get(key string) (string, bool, error) {
    if err := v.send("GET", key); err != nil { return "", false, err }
    b, n, err := v.readBulk()
    if err != nil { return "", false, err }
    if n < 0 { return "", false, nil }
    return string(b), true, nil
}

// Set sets key to value with expiration ttl. If ttl<=0, it's set without expiry.
func (v *Valkey) Set(key, value string, ttl time.Duration) error {
    if ttl > 0 {
        // SET key value EX seconds
        secs := int(ttl.Seconds())
        return v.sendExpectOK("SET", key, value, "EX", strconv.Itoa(secs))
    }
    return v.sendExpectOK("SET", key, value)
}

// Del deletes keys and returns number of removed keys.
func (v *Valkey) Del(keys ...string) (int64, error) {
    if len(keys) == 0 { return 0, nil }
    args := append([]string{"DEL"}, keys...)
    if err := v.send(args...); err != nil { return 0, err }
    n, err := v.readInteger()
    return n, err
}

// send writes a RESP Array of Bulk Strings for the provided args.
func (v *Valkey) send(args ...string) error {
    if v == nil || v.rw == nil { return errors.New("valkey: closed") }
    // *<n>\r\n
    if _, err := fmt.Fprintf(v.rw, "*%d\r\n", len(args)); err != nil { return err }
    for _, a := range args {
        if _, err := fmt.Fprintf(v.rw, "$%d\r\n%s\r\n", len(a), a); err != nil { return err }
    }
    return v.rw.Flush()
}

func (v *Valkey) sendExpectOK(args ...string) error {
    if err := v.send(args...); err != nil { return err }
    _, err := v.readSimpleOK()
    return err
}

// RESP readers (simple subset)
func (v *Valkey) readByte() (byte, error) { return v.rw.ReadByte() }

func (v *Valkey) readLine() (string, error) {
    s, err := v.rw.ReadString('\n')
    if err != nil { return "", err }
    return strings.TrimRight(s, "\r\n"), nil
}

func (v *Valkey) readSimpleString() (string, error) {
    b, err := v.readByte()
    if err != nil { return "", err }
    switch b {
    case '+':
        return v.readLine()
    case '-':
        s, _ := v.readLine()
        return "", errors.New(s)
    default:
        return "", fmt.Errorf("valkey: expected simple string, got %q", b)
    }
}

func (v *Valkey) readSimpleOK() (string, error) {
    s, err := v.readSimpleString()
    if err != nil { return "", err }
    if !strings.HasPrefix(strings.ToUpper(s), "OK") {
        return s, fmt.Errorf("valkey: expected OK, got %q", s)
    }
    return s, nil
}

func (v *Valkey) readInteger() (int64, error) {
    b, err := v.readByte()
    if err != nil { return 0, err }
    switch b {
    case ':':
        s, err := v.readLine()
        if err != nil { return 0, err }
        n, err := strconv.ParseInt(s, 10, 64)
        if err != nil { return 0, err }
        return n, nil
    case '-':
        s, _ := v.readLine()
        return 0, errors.New(s)
    default:
        return 0, fmt.Errorf("valkey: expected integer, got %q", b)
    }
}

func (v *Valkey) readBulk() ([]byte, int, error) {
    b, err := v.readByte()
    if err != nil { return nil, 0, err }
    switch b {
    case '$':
        s, err := v.readLine()
        if err != nil { return nil, 0, err }
        n, err := strconv.Atoi(s)
        if err != nil { return nil, 0, err }
        if n < 0 {
            return nil, -1, nil
        }
        buf := make([]byte, n+2)
        if _, err := v.rw.Read(buf); err != nil { return nil, 0, err }
        // strip CRLF
        return buf[:n], n, nil
    case '-':
        s, _ := v.readLine()
        return nil, 0, errors.New(s)
    default:
        return nil, 0, fmt.Errorf("valkey: expected bulk, got %q", b)
    }
}

func getenv(k, def string) string { if v := os.Getenv(k); v != "" { return v }; return def }
func firstNonEmpty(a, b string) string { if a != "" { return a }; return b }

