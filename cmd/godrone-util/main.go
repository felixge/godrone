package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
)

var (
	addr    = flag.String("addr", "192.168.1.1", "Addr of the drone.")
	verbose = flag.Int("verbose", 0, "Verbose flag; values gt 0 are passed onwards to the program on the drone.")
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
	killCmd      = "program.elf program.elf.respawner.sh " + godroneBin
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	tmpDir, err := ioutil.TempDir("", tmpDirPrefix)
	if err != nil {
		log.Fatalf("Could not create tmp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Expected a command: run")
		return
	}

	switch cmd := flag.Arg(0); cmd {
	case "run":
		pkg := flag.Arg(1)
		if pkg == "" {
			pkg = godronePkg
		}
		run(pkg, tmpDir)
	default:
		log.Fatalf("Unknown command: %s", cmd)
	}
}

func run(pkg, buildDir string) {
	binName := filepath.Base(pkg)
	log.Printf("Getting %s", pkg)
	get := exec.Command("go", "get", pkg)
	get.Dir = buildDir
	if output, err := get.CombinedOutput(); err != nil {
		log.Fatalf("Compile error: %s: %s", err, output)
	}
	log.Printf("Cross compiling")
	build := exec.Command("go", "build", pkg)
	build.Env = append(os.Environ(), "GOOS="+goOs, "GOARCH="+goArch)
	build.Dir = buildDir
	if output, err := build.CombinedOutput(); err != nil {
		log.Printf("Compile error: %s: %s", err, output)
		log.Print("If you need help setting up Go cross-compiling see:")
		log.Fatal("  http://godrone.io/en/latest/contributor/install_from_source.html")
	}
	log.Printf("Establishing telnet connection")
	telnet, err := DialTelnet(net.JoinHostPort(*addr, telnetPort))
	if err != nil {
		log.Printf("Telnet connect error: %s", err)
		log.Fatal("Check to be sure that your computer is attached to the drone's wifi network.")
	}
	defer telnet.Close()
	log.Printf("Killing firmware (restart drone to get it back)")
	if out, err := telnet.Exec("killall -q -KILL " + killCmd); err != nil {
		if string(out) != "" {
			log.Fatalf("Failed to kill firmware: %s: %s", err, out)
		}
	}
	file, err := os.Open(filepath.Join(buildDir, binName))
	if err != nil {
		log.Fatalf("Could not open godrone file: %s", err)
	}
	defer file.Close()
	log.Printf("Establishing ftp connection")
	ftp, err := ftp.Connect(net.JoinHostPort(*addr, ftpPort))
	if err != nil {
		log.Fatalf("FTP connect error: %s", err)
	}
	defer ftp.Quit()
	dstPath := path.Join(godroneDir, godroneBin)
	log.Printf("Uploading %s to %s", binName, dstPath)
	ftp.MakeDir(godroneDir)
	if err := ftp.Stor(dstPath, file); err != nil {
		log.Fatalf("Failed to upload: %s", err)
	}
	ftp.Quit()
	file.Close()
	// otherwise the drone starts counting time from Jan 1st 2000 after restart
	// which is annoying when trying to correlate log output to observed behavior
	log.Printf("Syncing drone clock with host clock")
	now := time.Now().Format("2006-01-02 15:04:05")
	if out, err := telnet.Exec(fmt.Sprintf("date -s '%s'", now)); err != nil {
		log.Fatalf("Failed to sync clock: %s: %s", err, out)
	}
	log.Printf("Starting target program")
	if out, err := telnet.Exec("cd '" + path.Join(ftpDir, godroneDir) + "'"); err != nil {
		log.Fatalf("Failed to change directory: %v: %v", err, out)
	}
	if out, err := telnet.Exec("chmod +x '" + godroneBin + "'"); err != nil {
		log.Fatalf("Failed to make godrone executable: %v: %v", err, out)
	}

	cmd := fmt.Sprintf("./%v", godroneBin)
	if *verbose != 0 {
		cmd += fmt.Sprintf(" -verbose=%v", *verbose)
	}

	log.Printf("Running %v", cmd)
	if err := telnet.ExecRawWriter(cmd, os.Stdout); err != nil {
		log.Printf("Failed to run %v: %v", godroneBin, err)
	}
	telnet.Close()
}
