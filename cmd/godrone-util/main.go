package main

import "flag"
import goftp "github.com/jlaffaye/goftp"
import "io/ioutil"

import "log"
import "net"

import (
	"os"
	"os/exec"
	"path"
)
import "path/filepath"

var (
	addr = flag.String("addr", "192.168.1.1", "Addr of the drone.")
)

const (
	godronePkg   = "github.com/felixge/godrone/cmd/godrone"
	godroneBin   = "godrone"
	godroneDir   = "godrone"
	goOs         = "linux"
	goArch       = "arm"
	ftpPort      = "21"
	ftpDir       = "/data/video"
	telnetPort   = "23"
	tmpDirPrefix = "godrone-util"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	flag.Parse()
	tmpDir, err := ioutil.TempDir("", tmpDirPrefix)
	if err != nil {
		log.Fatalf("Could not create tmp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)
	switch cmd := flag.Arg(0); cmd {
	case "run":
		run(tmpDir)
	default:
		log.Fatalf("Unknown command: %s", cmd)
	}
}

func run(dir string) {
	log.Printf("Cross compiling %s", godroneBin)
	build := exec.Command("go", "build", godronePkg)
	build.Env = append(os.Environ(), "GOOS="+goOs, "GOARCH="+goArch)
	build.Dir = dir
	if output, err := build.CombinedOutput(); err != nil {
		log.Fatalf("Compile error: %s: %s", err, output)
	}
	file, err := os.Open(filepath.Join(dir, godroneBin))
	if err != nil {
		log.Fatalf("Could not open godrone file: %s", err)
	}
	defer file.Close()
	log.Printf("Uploading %s", godroneBin)
	ftp, err := goftp.Connect(net.JoinHostPort(*addr, ftpPort))
	if err != nil {
		log.Fatalf("FTP connect error: %s", err)
	}
	defer ftp.Quit()
	ftp.MakeDir(godroneDir)
	if err := ftp.Stor(path.Join(godroneDir, godroneBin), file); err != nil {
		log.Fatalf("Failed to upload: %s", err)
	}
	ftp.Quit()
	file.Close()
	log.Printf("Starting %s", godroneBin)
	telnet, err := DialTelnet(net.JoinHostPort(*addr, telnetPort))
	if err != nil {
		log.Fatalf("Telnet connect error: %s", err)
	}
	defer telnet.Close()
	if out, err := telnet.Exec("cd '" + path.Join(ftpDir, godroneDir) + "'"); err != nil {
		log.Fatalf("Failed to change directory: %s: %s", err, out)
	}
	if out, err := telnet.Exec("chmod +x '" + godroneBin + "'"); err != nil {
		log.Fatalf("Failed to make godrone executable: %s: %s", err, out)
	}
	log.Printf("Running %s", godroneBin)
	if err := telnet.ExecRawWriter("./"+godroneBin, os.Stdout); err != nil {
		log.Printf("Failed to run %s: %s", godroneBin, err)
	}
	telnet.Close()
}
