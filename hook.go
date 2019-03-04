package main

import (
    "os"
    "os/exec"
    "strings"
    "strconv"
    "regexp"
    "log"

    "github.com/buty4649/Sibyl/nlc"
    "github.com/prometheus/procfs"
)

type Proc struct {
    PID     int
    CmdLine string
    Command string
    Proc    procfs.Proc
}

func (s *Sibyl) Fork(ev nlc.ForkEvent) error {
    proc := newProc(int(ev.ChildPID))
    if proc == nil {
        return nil
    }

    for _, c := range s.config.Fork {
        r := regexp.MustCompile(c.CmdLine)
        if r.MatchString(proc.CmdLine) == false {
            continue
        }
        exec := exec.Command("sh", "-c", c.Exec)
        exec.Env = proc.buildEnv()
        out, err := exec.CombinedOutput()
        if err != nil {
            return err
        }
        if out != nil {
            log.Printf("Fork Hook: %s", out)
        }
    }
    return nil
}

func (s *Sibyl) Exec(ev nlc.ExecEvent) error {
    proc := newProc(int(ev.ProcessPID))
    if proc == nil {
        return nil
    }

    for _, c := range s.config.Exec {
        r := regexp.MustCompile(c.CmdLine)
        if r.MatchString(proc.CmdLine) == false {
            continue
        }
        exec := exec.Command("sh", "-c", c.Exec)
        exec.Env = proc.buildEnv()
        out, err := exec.CombinedOutput()
        if err != nil {
            return err
        }
        if out != nil {
            log.Printf("Exec Hook: %s", out)
        }
    }
    return nil
}

func (s *Sibyl) Coredump(ev nlc.CoredumpEvent) error {
    proc := newProc(int(ev.ProcessPID))
    if proc == nil {
        return nil
    }

    for _, c := range s.config.Coredump {
        r := regexp.MustCompile(c.CmdLine)
        if r.MatchString(proc.CmdLine) == false {
            continue
        }
        exec := exec.Command("sh", "-c", c.Exec)
        exec.Env = proc.buildEnv()
        out, err := exec.CombinedOutput()
        if err != nil {
            return err
        }
        if out != nil {
            log.Printf("Coredump Hook: %s", out)
        }
    }
    return nil
}

func (s *Sibyl) Exit(ev nlc.ExitEvent) error {
    proc := newProc(int(ev.ProcessPID))
    if proc == nil {
        return nil
    }

    for _, c := range s.config.Exit {
        r := regexp.MustCompile(c.CmdLine)
        if r.MatchString(proc.CmdLine) == false {
            continue
        }
        exec := exec.Command("sh", "-c", c.Exec)
        exec.Env = proc.buildEnv()
        out, err := exec.CombinedOutput()
        if err != nil {
            return err
        }
        if out != nil {
            log.Printf("Exit Hook: %s", out)
        }
    }
    return nil
}

func newProc(pid int) *Proc {
    proc, err := procfs.NewProc(pid)
    if err != nil {
        // /proc/$PID no such file or directory
        return nil
    }

    cl, err := proc.CmdLine()
    if err != nil {
        return nil
    }
    cmdline := strings.Join(cl, " ")

    comm, err := proc.Comm()
    if err != nil {
        return nil
    }

    return &Proc{pid, cmdline, comm, proc}
}

func (p *Proc) buildEnv() []string {
    env := append(os.Environ(), "SIBYL_CMDLINE=" + p.CmdLine)
    env = append(env, "SIBYL_COMM=" + p.Command)
    env = append(env, "SIBYL_PID=" + strconv.Itoa(p.PID))
    return env
}
