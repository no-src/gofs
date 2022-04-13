package daemon

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/no-src/gofs/util/osutil"
	"github.com/no-src/log"
)

// SubprocessTag mark the current process is subprocess
const SubprocessTag = "sub"

var shutdown = make(chan bool, 1)

// Daemon running as a daemon process, and create a subprocess for working
func Daemon(recordPid bool, daemonDelay time.Duration, monitorDelay time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("daemon process error. %v", r)
		}
	}()

	for {
		if wait(daemonDelay) {
			return
		}
		p, err := startSubprocess()
		if err == nil && p != nil {
			if recordPid {
				log.ErrorIf(writePidFile(os.Getppid(), os.Getpid(), p.Pid), "write pid info to file error")
			}
			if monitor(p.Pid, monitorDelay) {
				return
			}
		}
	}
}

// startSubprocess start a subprocess for working
func startSubprocess() (*os.Process, error) {
	attr := &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}
	// try to check stdin
	// if compile with [-ldflags="-H windowsgui"] on Windows system, stdin will get error
	if osutil.IsWindows() {
		_, stdInErr := os.Stdin.Stat()
		if stdInErr != nil {
			attr = &os.ProcAttr{Files: []*os.File{nil, nil, nil}}
		}
	}
	args := os.Args
	// use "-sub" to tag sub process
	args = append(args, "-"+SubprocessTag)
	exeFile, err := os.Executable()
	if err == nil {
		args[0] = exeFile
	} else {
		log.Error(err, "get current executable error")
	}
	p, err := os.StartProcess(args[0], args, attr)
	if err == nil && p != nil {
		log.Info("[%d] start subprocess success", p.Pid)
	} else {
		log.Error(err, "start subprocess error")
	}
	return p, err
}

// monitor start to monitor the subprocess, create a new subprocess to work if subprocess is dead
func monitor(pid int, monitorDelay time.Duration) (isShutdown bool) {
	for {
		if wait(monitorDelay) {
			return true
		}

		p, err := os.FindProcess(pid)
		if err != nil {
			log.Error(err, "[%d] subprocess status error", pid)
			if p != nil {
				log.Info("[%d] try to kill the subprocess", pid)
				log.ErrorIf(p.Kill(), "[%d] try to kill the subprocess error", pid)
			}
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
func writePidFile(ppid, pid, subPid int) error {
	fName := "pid"
	f, err := os.Create(fName)
	if err == nil {
		writer := bufio.NewWriter(f)
		if _, err = writer.WriteString(fmt.Sprintf("%d\n%d\n%d\n", ppid, pid, subPid)); err != nil {
			return err
		}
		if err = writer.Flush(); err != nil {
			return err
		}
		err = f.Close()
	}
	return err
}

// KillPPid kill parent process
func KillPPid() {
	ppid := os.Getppid()
	if ppid > 0 {
		p, err := os.FindProcess(ppid)
		if err == nil {
			if p != nil {
				log.ErrorIf(p.Kill(), "[%d] kill parent process error", ppid)
			}
		} else {
			log.Error(err, "[%d] find parent process error", ppid)
		}
	}

}

// Shutdown send a shutdown notify to the current daemon
func Shutdown() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	shutdown <- true
	close(shutdown)
	return err
}

func wait(d time.Duration) (isShutdown bool) {
	select {
	case isShutdown = <-shutdown:
		{
			if isShutdown {
				return isShutdown
			}
		}
	case <-time.After(d):

	}
	return false
}
