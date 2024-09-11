package main

import (
    "fmt"
    "os"
    "os/exec"
    "golang.org/x/sys/unix"
)

func main() {
    switch os.Args[1] {
    case "run":
        parent()
    case "child":
        child()
    default:
        panic("what do i do?")
    }
}

func parent() {
    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
    cmd.SysProcAttr = &unix.SysProcAttr{
        Cloneflags: unix.CLONE_NEWUTS | unix.CLONE_NEWPID | unix.CLONE_NEWNS,
    }
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if err := cmd.Run(); err != nil {
        fmt.Println("ERROR: ", err)
        os.Exit(1)
    }
}

func child() {
    must(unix.Mount("rootfs", "rootfs", "", unix.MS_BIND, nil))
    must(os.MkdirAll("rootfs/oldrootfs", 0700))
    must(unix.PivotRoot("rootfs", "rootfs/oldrootfs"))
    must(os.Chdir("/"))

    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        fmt.Println("ERROR: ", err)
        os.Exit(1)
    }
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}