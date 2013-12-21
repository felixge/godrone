// Command deploy allows people to deploy and run GoDrone binaries on their
// ardrone.
package main

import (
	"bitbucket.org/kardianos/osext"
	"fmt"
	ftp "github.com/jlaffaye/goftp"
	ptelnet "github.com/ziutek/telnet"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const cmdName = "deploy"

func Printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

func Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func main() {
	dir, err := osext.ExecutableFolder()
	if err != nil {
		Fatalf("Could not determine ExecutableFolder: %s", err)
	}

	task := DeployTask{
		Host:         "192.168.1.1",
		FtpPort:      "21",
		FtpDir:       "/data/video",
		TelnetPort:   "23",
		TelnetPrompt: "# ",
		Src:          dir,
		Dst:          "godrone",
		NetTimeout:   5 * time.Second,
	}
	if err := task.Run(); err != nil {
		Fatalf("Failed to deploy godrone: %s", err)
	}
}

type DeployTask struct {
	// Host is the IP address of the drone.
	Host string
	// FtpPort is the ftp port.
	FtpPort string
	// FtpDir is the absolute path of the FTP directory on the drone file system.
	FtpDir string
	// TelnetPort is the telnet port.
	TelnetPort string
	// TelnetPrompt is the prompt string used by the drone's shell.
	TelnetPrompt string
	// Src is the absolute path to the godrone folder on the host.
	Src string
	// Dst is path relative to the FtpDir to deploy godrone on the drone.
	Dst string
	// NetTimeout is the timeout for all networking operations.
	NetTimeout time.Duration
}

func (t DeployTask) Run() error {
	dir, err := os.Open(t.Src)
	if err != nil {
		return err
	}
	defer dir.Close()

	ftpAddr := net.JoinHostPort(t.Host, t.FtpPort)
	Printf("Connecting to %s", ftpAddr)
	// @TODO figure out how to set a connection timeout for ftp, might require
	// sending a patch to goftp.
	ftpConn, err := ftp.Connect(ftpAddr)
	if err != nil {
		return err
	}
	defer ftpConn.Quit()

	Printf("Connected")

	dstDir := t.Dst + ".next"
	ftpConn.MakeDir(dstDir)

	entries, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		baseName := filepath.Base(name)
		if strings.HasPrefix(baseName, cmdName) {
			// don't upload the deployer itself
			continue
		}

		srcName := filepath.Join(t.Src, name)
		dstName := path.Join(dstDir, name)

		file, err := os.Open(srcName)
		if err != nil {
			return err
		}

		Printf("Uploading %s", name)
		if err := ftpConn.Stor(dstName, file); err != nil {
			return err
		}
	}

	telnetAddr := net.JoinHostPort(t.Host, t.TelnetPort)
	Printf("Connecting to %s", telnetAddr)
	telnetConn, err := net.DialTimeout("tcp", telnetAddr, t.NetTimeout)
	if err != nil {
		return err
	}
	defer telnetConn.Close()
	Printf("Connected")

	telnet, err := ptelnet.NewConn(telnetConn)
	if err != nil {
		return err
	}

	Printf("Running start.sh on drone")
	if _, err := telnet.ReadUntil(t.TelnetPrompt); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(telnet, "cd %s/%s && sh start.sh\n", t.FtpDir, dstDir); err != nil {
		return err
	}
	go io.Copy(telnet, os.Stdin)
	if _, err := io.Copy(os.Stdout, telnet); err != nil {
		return err
	}
	return nil
}
