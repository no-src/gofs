package daemon

import (
	"bufio"
	"fmt"
	"github.com/no-src/log"
	"os"
	"runtime"
	"time"
)

const SubprocessTag = "sub"

func Daemon(recordPid bool, daemonDelay time.Duration, monitorDelay time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("daemon process error. %v", r)
		}
	}()
	for {
		<-time.After(daemonDelay)
		p, err := startSubprocess()
		if err == nil && p != nil {
			if recordPid {
				writePidFile(os.Getppid(), os.Getpid(), p.Pid)
			}
			monitor(p.Pid, monitorDelay)
		}
	}
}

// startSubprocess start a subprocess to work
func startSubprocess() (*os.Process, error) {
	attr := &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}
	// try to check stdin
	// if compile with [-ldflags="-H windowsgui"] on Windows system, stdin will get error
	if isWindows() {
		_, stdInErr := os.Stdin.Stat()
		if stdInErr != nil {
			attr = &os.ProcAttr{Files: []*os.File{nil, nil, nil}}
		}
	}
	args := os.Args
	// use "-sub" to tag sub process
	args = append(args, "-"+SubprocessTag)
	p, err := os.StartProcess(os.Args[0], args, attr)
	if err == nil && p != nil {
		log.Info("[%d] start subprocess success", p.Pid)
	} else {
		log.Error(err, "start subprocess error")
	}
	return p, err
}

// monitor start to monitor the subprocess, create a new subprocess to work if subprocess is dead
func monitor(pid int, monitorDelay time.Duration) {
	for {
		<-time.After(monitorDelay)
		p, err := os.FindProcess(pid)
		if err != nil {
			if p != nil {
				log.Info("[%d] try to kill the subprocess", pid)
				p.Kill()
			}
			log.Error(err, "[%d] subprocess status error", pid)
			return
		}
		if p == nil {
			log.Info("[%d] subprocess is not found", pid)
			return
		}

		// wait subprocess to exit
		stat, err := p.Wait()
		if err != nil || stat.Exited() {
			log.Info("[%d] subprocess is exited", pid)
			return
		}

	}
}

// writePidFile write current process and subprocess pid to pid file
// row 1: record parent process pid (bash,cmd,explorer etc.)
// row 2: record current process pid (daemon)
// row 3: record subprocess pid (worker)
func writePidFile(ppid, pid, subPid int) {
	fName := "pid"
	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		log.Error(err, "open pid file error [%s]", fName)
	} else {
		writer := bufio.NewWriter(f)
		writer.WriteString(fmt.Sprintf("%d\n", ppid))
		writer.WriteString(fmt.Sprintf("%d\n", pid))
		writer.WriteString(fmt.Sprintf("%d\n", subPid))
		writer.Flush()
		f.Close()
	}
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

// KillPPid kill parent process
func KillPPid() {
	ppid := os.Getppid()
	if ppid > 0 {
		p, err := os.FindProcess(ppid)
		if err == nil {
			if p != nil {
				err = p.Kill()
				if err != nil {
					log.Error(err, "[%d] kill parent process error", ppid)
				}
			}
		} else {
			log.Error(err, "[%d] find parent process error", ppid)
		}
	}

}
