package main


import (
	"io/ioutil"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("help")
	}
}

func run() {
	fmt.Printf("Running %v\n", os.Args[2:])
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr =  &syscall.SysProcAttr {
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID,
	}

	must(cmd.Run())
}

func child() {
	cg()
	fmt.Printf("Running %v\n", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot("/home/asfandyar/ubuntufs"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(syscall.Mount("something", "/mytemp", "tmpfs", 0, ""))


	must(cmd.Run())
	must(syscall.Unmount("/proc", 0))
	must(syscall.Unmount("/mytemp", 0))
}

func cg() {
	cgroups := "/sys/fs/cgroup/"

	mem := filepath.Join(cgroups, "memory")
	os.Mkdir(filepath.Join(mem, "asfandyar"), 0755)
	must(ioutil.WriteFile(filepath.Join(mem, "asfandyar/memory.limit_in_bytes"), []byte("999424"), 0700))
	must(ioutil.WriteFile(filepath.Join(mem, "asfandyar/memory.memsw.limit_in_bytes"), []byte("999424"), 0700))
	must(ioutil.WriteFile(filepath.Join(mem, "asfandyar/notify_on_release"), []byte("1"), 0700))

	pid := strconv.Itoa(os.Getpid())
	must(ioutil.WriteFile(filepath.Join(mem, "asfandyar/cgroup.procs"), []byte(pid), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}