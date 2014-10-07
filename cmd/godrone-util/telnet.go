package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	gotelnet "github.com/ziutek/telnet"
)

// DialTelnet creates a new telnet connection.
func DialTelnet(addr string) (*Telnet, error) {
	conn, err := gotelnet.Dial("tcp", addr)
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

// Exec executes the given command and returns the output. Exit codes codes are
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
	if err := t.discardUntil(cmd + "\r\n"); err != nil {
		return err
	}
	if err := t.copyUntil(output, t.prompt); err != nil {
		return err
	}
	return nil
}

// discardUntil discards the output until deli, including the delim itself.
func (t *Telnet) discardUntil(delim string) error {
	return t.copyUntil(ioutil.Discard, delim)
}

// discardUntil copies the output until delim, excluding the delim itself.
// @TODO Could be made more efficient by reading/writing more than one byte
// at a time.
func (t *Telnet) copyUntil(dst io.Writer, delim string) error {
	buf := ""
	for {
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
	return t.discardUntil(t.prompt)
}

// Close closes the connection.
func (t *Telnet) Close() error {
	t.Exec("exit")
	return t.conn.Close()
}
