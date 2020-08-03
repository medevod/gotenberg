package core

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func Exec(ctx context.Context, bin string, args ...string) error {
	cmd := exec.CommandContext(ctx, bin, args...)
	// TODO: pipe command output.

	// See https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773.
	kill := func() {
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err == nil {
			return
		}
		if !strings.Contains(err.Error(), "no such process") {
			fmt.Printf("%s\n", err) // TODO: use logging.
		}
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err := cmd.Start()
	if err != nil {
		return err
	}

	result := make(chan error, 1)
	go func() {
		result <- cmd.Wait()
	}()

	select {
	case err := <-result:
		kill()
		return err
	case <-ctx.Done():
		kill()
		return ctx.Err()
	}
}
