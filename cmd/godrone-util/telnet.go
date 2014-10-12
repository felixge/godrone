package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	gotelnet "github.com/ziutek/telnet"
)

// DialTelnet creates a new telnet connection.
func DialTelnet(addr string) (*Telnet, error) {
	conn, err := gotelnet.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	telnet := &Telnet{conn: conn, prompt: "# "}
	if err := telnet.waitPrompt(); err != nil {
		return nil, err
	}
	return telnet, err
}

// Telnet is a telnet connection.
type Telnet struct {
	conn   *gotelnet.Conn
	prompt string
}

// Exec executes the given command and returns the output. Exit codes > 0 are
// considered errors.
func (t *Telnet) Exec(cmd string) ([]byte, error) {
	out, err := t.ExecRaw(cmd)
	if err != nil {
		return out, err
	}
	code, err := t.ExecRaw("echo $?")
	if err != nil {
		return code, err
	}
	var intCode int
	if _, err := fmt.Sscanf(string(code), "%d", &intCode); err != nil {
		return code, err
	}
	if intCode != 0 {
		return out, TelnetExitError{code: intCode}
	}
	return out, nil
}

type TelnetExitError struct {
	code int
}

func (e TelnetExitError) Error() string {
	return fmt.Sprintf("exit code %d", e.code)
}

func (e TelnetExitError) Code() int {
	return e.code
}

// ExecRaw executes the given command and returns the output. Exit codes are
// not considered.
func (t *Telnet) ExecRaw(cmd string) ([]byte, error) {
	var buf bytes.Buffer
	err := t.ExecRawWriter(cmd, &buf)
	return buf.Bytes(), err
}

// ExecRawWriter executes the given command and writes the output into the
// writer. Exit codes are not considered.
func (t *Telnet) ExecRawWriter(cmd string, output io.Writer) error {
	if _, err := fmt.Fprintf(t.conn, cmd+"\n"); err != nil {
		return err
	}
	if err := t.discardUntil(cmd+"\r\n", 0); err != nil {
		return err
	}
	if err := t.copyUntil(output, t.prompt, 0); err != nil {
		return err
	}
	return nil
}

// A telnetError wraps up the error code and what arrived instead
// of what was expected.
type telnetError struct {
	error error
	read  []byte
}

func (te telnetError) Error() string {
	if len(te.read) > 0 {
		return fmt.Sprintf("%v (unexpected input: %q)", te.error, string(te.read))
	} else {
		return te.error.Error()
	}
}

// discardUntil discards the output until delim, including the delim itself.
func (t *Telnet) discardUntil(delim string, to time.Duration) error {
	// save what arrived anyway; if there's a timeout we'll give it
	// to them in an error.
	var buf bytes.Buffer

	err := t.copyUntil(&buf, delim, to)
	if err != nil {
		return telnetError{error: err, read: buf.Bytes()}
	}
	return err
}

var noDeadline = time.Time{}

// copyUntil copies the output until delim, excluding the delim itself.
// If timeout is not 0, it is used to prevent copyUntil from hanging
// if no input is arriving.
// @TODO Could be made more efficient by reading/writing more than one byte
// at a time.
func (t *Telnet) copyUntil(dst io.Writer, delim string, to time.Duration) error {
	buf := ""
	for {
		if to != time.Duration(0) {
			t.conn.SetDeadline(time.Now().Add(to))
			defer t.conn.SetDeadline(noDeadline)
		}
		b, err := t.conn.ReadByte()
		if err != nil {
			return err
		}
		buf += string(b)
		if strings.HasSuffix(buf, delim) {
			return nil
		}
		if len(buf) >= len(delim) {
			char := buf[len(buf)-len(delim)]
			if _, err := dst.Write([]byte{char}); err != nil {
				return err
			}
		}
	}
}

// waitPrompt waits for the prompt to show.
func (t *Telnet) waitPrompt() error {
	return t.discardUntil(t.prompt, 1*time.Second)
}

// Close closes the connection.
func (t *Telnet) Close() error {
	t.Exec("exit")
	return t.conn.Close()
}
