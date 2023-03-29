package daemon

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/no-src/gofs/util/osutil"
	"github.com/no-src/gofs/wait"
	"github.com/no-src/log"
)

// SubprocessTag mark the current process is subprocess
const SubprocessTag = "sub"

// Daemon support to running daemon process and create subprocess for working
type Daemon struct {
	shutdown chan struct{}
}

// New create an instance of Daemon
func New() *Daemon {
	return &Daemon{
		shutdown: make(chan struct{}, 1),
	}
}

// Run running as a daemon process, and create a subprocess for working, the first argument must be an absolute path of the program name
func (d *Daemon) Run(args []string, recordPid bool, daemonDelay time.Duration, monitorDelay time.Duration, wd wait.Done) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("daemon process error. %v", r)
			log.Error(err, "daemon exited by panic")
			wd.DoneWithError(err)
		}
	}()

	for {
		if d.waitShutdown(daemonDelay) {
			wd.Done()
			log.Info("daemon exited by shutdown")
			return
		}
		p, err := d.startSubprocess(args)
		if err == nil && p != nil {
			if recordPid {
				log.ErrorIf(d.writePidFile(os.Getppid(), os.Getpid(), p.Pid), "write pid info to file error")
			}
			if d.monitor(p.Pid, monitorDelay) {
				wd.Done()
				log.Info("daemon exited by shutdown")
				return
			}
		}
	}
}

// startSubprocess start a subprocess for working
func (d *Daemon) startSubprocess(args []string) (*os.Process, error) {
	attr := &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}
	// try to check stdin
	// if compile with [-ldflags="-H windowsgui"] on Windows system, stdin will get error
	if osutil.IsWindows() {
		_, stdInErr := os.Stdin.Stat()
		if stdInErr != nil {
			attr = &os.ProcAttr{Files: []*os.File{nil, nil, nil}}
		}
	}
	// use "-sub" to tag sub process
	args = append(args, "-"+SubprocessTag)
	p, err := os.StartProcess(args[0], args, attr)
	if err == nil && p != nil {
		log.Info("[%d] start subprocess success", p.Pid)
	} else {
		log.Error(err, "start subprocess error")
	}
	return p, err
}

// monitor start to monitor the subprocess, create a new subprocess to work if subprocess is dead
func (d *Daemon) monitor(pid int, monitorDelay time.Duration) (isShutdown bool) {
	for {
		if d.waitShutdown(monitorDelay) {
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
func (d *Daemon) writePidFile(ppid, pid, subPid int) error {
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
func (d *Daemon) KillPPid() {
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
func (d *Daemon) Shutdown() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	close(d.shutdown)
	return err
}

func (d *Daemon) waitShutdown(du time.Duration) (isShutdown bool) {
	select {
	case <-d.shutdown:
		return true
	case <-time.After(du):

	}
	return false
}
