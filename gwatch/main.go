package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func run() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var paras []string
	if len(os.Args) > 2 {
		paras = os.Args[2:]
	}

	cmd := exec.CommandContext(ctx, os.Args[1], paras...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
		if err == context.Canceled {
			// Kill it:
			if err := cmd.Process.Kill(); err != nil {
				log.Printf("failed to kill process: %s", err)
			}
		}
		return
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//t := time.Now()
	//log.Printf("%s ", t.Format("2006-01-02 15:04:05"))
	if outStr != "" {
		log.Printf("stdout: %s", outStr)
	}
	if errStr != "" {
		log.Printf("stderr: %s", errStr)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("no enough argument")
		return
	}

	closeSignalChan := make(chan os.Signal, 1)
	signal.Notify(closeSignalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			run()
		case <-closeSignalChan:
			os.Exit(0)
		}
	}
}
