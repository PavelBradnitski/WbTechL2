// main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(path + ">")

		if !scanner.Scan() {
			fmt.Println()
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		input = expandEnvVars(input)

		ctx, cancel := context.WithCancel(context.Background())

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGQUIT)

		go func() {
			s := <-sigCh
			if s == syscall.SIGQUIT {
				os.Exit(0)
			}
			cancel()
		}()

		exitCode := runCommand(ctx, input)

		fmt.Printf("exit code: %d\n", exitCode)

		close(sigCh)
		cancel()
	}
}

func expandEnvVars(line string) string {
	for _, part := range strings.Fields(line) {
		if strings.HasPrefix(part, "$") && len(part) > 1 {
			val := os.Getenv(part[1:])
			line = strings.ReplaceAll(line, part, val)
		}
	}
	return line
}

func runCommand(ctx context.Context, input string) int {
	if strings.Contains(input, "|") {
		return runPipeline(ctx, strings.Split(input, "|"))
	}

	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		return cmdCd(args)
	case "pwd":
		return cmdPwd(os.Stdout)
	case "echo":
		return cmdEcho(args, os.Stdout)
	case "kill":
		return cmdKill(ctx, args)
	case "ps":
		return cmdPs(ctx, args, os.Stdout)
	}

	return 0
}

func cmdCd(args []string) int {
	if len(args) < 2 {
		fmt.Println("cd: missing operand")
		return 1
	}
	path := args[1]
	if !filepath.IsAbs(path) {
		cwd, _ := os.Getwd()
		path = filepath.Join(cwd, path)
	}
	if err := os.Chdir(path); err != nil {
		fmt.Println("cd:", err)
		return 1
	}
	return 0
}

func cmdPwd(w io.Writer) int {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, "pwd:", err)
		return 1
	}
	fmt.Fprintln(w, dir)
	return 0
}

func cmdEcho(args []string, w io.Writer) int {

	fmt.Fprintln(w, strings.Join(args[1:], " "))
	return 0
}

func cmdKill(ctx context.Context, args []string) int {
	if len(args) < 2 {
		fmt.Println("kill: missing pid")
		return 1
	}

	cmdName := ""
	if runtime.GOOS == "windows" {
		cmdName = "taskkill"
		args = append([]string{"/PID"}, args[1:]...)
		args = append(args, "/F")
	} else {
		cmdName = "kill"
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func cmdPs(ctx context.Context, args []string, w io.Writer) int {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "tasklist")
	} else {
		cmd = exec.CommandContext(ctx, "ps", args[1:]...)
	}
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func runPipeline(ctx context.Context, cmds []string) int {
	var commands []*exec.Cmd
	for _, c := range cmds {
		args := strings.Fields(strings.TrimSpace(c))
		if len(args) == 0 {
			return 1
		}
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		commands = append(commands, cmd)
	}

	for i := 0; i < len(commands)-1; i++ {
		stdout, _ := commands[i].StdoutPipe()
		commands[i+1].Stdin = stdout
	}

	commands[len(commands)-1].Stdout = os.Stdout
	commands[len(commands)-1].Stderr = os.Stderr

	for _, cmd := range commands {
		start := time.Now()
		if err := cmd.Run(); err != nil {
			fmt.Println("error:", err)
			return 1
		}
		dur := time.Since(start)
		fmt.Printf("end cmd: %v, dur: %v\n", cmd, dur)
	}

	return 0
}
