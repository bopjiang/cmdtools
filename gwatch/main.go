package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	paras := []string{"-c", fmt.Sprintf(`"%s"`, strings.Join(os.Args[1:], " "))}
	cmd := exec.CommandContext(ctx, "bash", paras...)
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
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			run()
		}
	}
}
