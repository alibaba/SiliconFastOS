package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"oskit/util"
	"strings"

	"github.com/containerd/containerd/platforms"
	log "github.com/sirupsen/logrus"
)

func isExecErrNotFound(err error) bool {
	eerr, ok := err.(*exec.Error)
	if !ok {
		return false
	}
	return eerr.Err == exec.ErrNotFound
}

func run(args []string) {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	runCmd.Usage = func() {
		fmt.Printf("USAGE: %s run [options] <img-name>\n\n", os.Args[0])
		fmt.Printf("Options:\n\n")
		runCmd.PrintDefaults()
	}

	local := runCmd.Bool("local", false, "use local qemu binary instead of docker")
	ramdisk := runCmd.Bool("ramdisk", false, "Use <name>-initrd.ramdisk as initrd")
	_ = runCmd.Parse(args)

	for _, img := range runCmd.Args() {
		kernel := img + "-kernel"
		rootfs := img

		content, err := ioutil.ReadFile(img + "-cmdline")
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		// convert content to string
		cmdline := string(content)
		qemuImg := fmt.Sprintf("%s/silicon_fast/qemu_%s:7.1.0", util.DOCKER_REPO, platforms.DefaultSpec().Architecture)

		if *ramdisk {
			rootfs = rootfs + "-initrd.ramdisk"
			f, err := os.Stat(rootfs)
			if err != nil {
				fmt.Println("Can not find initrd :", err)
				return
			}
			cmdline = fmt.Sprintf("%s root=/dev/ram0 keepinitrd ramdisk_start=0x48000000 ramdisk_size=0x%X", cmdline, int(math.Ceil(float64(f.Size())/1024)))
		} else {
			rootfs = rootfs + "-initrd.cpio.gz"
		}

		var c []string
		if *local {
			c = []string{}
		} else {
			c = []string{"docker", "run", "--network", "host", "-v", "$(pwd):/workspace", "-it", qemuImg}
		}
		c = append(c, "qemu-system-aarch64 -smp 4 -m 8G -cpu neoverse-n1 -machine virt")
		c = append(c, "-kernel", kernel, "-initrd", rootfs, "-device virtio-net-pci,netdev=t0 -netdev user,id=t0", "-append", "\""+cmdline+"\"", "-nographic")

		exec := exec.CommandContext(context.Background(), "/bin/sh", "-c", strings.Join(c, " "))

		exec.Stdout = os.Stdout
		exec.Stderr = os.Stderr
		exec.Stdin = os.Stdin
		exec.Env = os.Environ()

		log.Debugf("Executing: %v", exec.Args)

		err = exec.Run()
		if err != nil {
			if isExecErrNotFound(err) {
				fmt.Errorf("linuxkit pkg requires docker to be installed")
			}
		}

	}
}
