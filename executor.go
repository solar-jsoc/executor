package executor

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type buff struct {
	bytes.Buffer
}

func (b *buff) Close() error {
	return nil
}

type command struct {
	cmd     string
	args    []string
	env     []string
	inputs  []io.ReadCloser
	outputs []io.WriteCloser
	stderr  io.Writer
}

func (c *command) run(dir string) error {
	inputs := []io.Reader{}
	for i := range c.inputs {
		inputs = append(inputs, io.Reader(c.inputs[i]))
		defer c.inputs[i].Close()
	}

	outputs := []io.Writer{}
	for i := range c.outputs {
		outputs = append(outputs, io.Writer(c.outputs[i]))
		defer c.outputs[i].Close()
	}

	cmd := exec.Command(c.cmd, c.args...)
	cmd.Dir = dir
	cmd.Env = c.env
	cmd.Stdin = io.MultiReader(inputs...)
	cmd.Stdout = io.MultiWriter(outputs...)
	cmd.Stderr = c.stderr
	// log.Printf("%+v\n", *c)
	return cmd.Run()
}

// Exec command
func Exec(cmd, dir string, env []string, stdin io.ReadCloser) ([]byte, []byte, error) {
	if dir == "" {
		tempDir, err := os.MkdirTemp("", "")
		if err != nil {
			return nil, nil, err
		}
		defer os.RemoveAll(tempDir)
		dir = tempDir
	}
	commands := parse(cmd, dir)
	if stdin != nil {
		commands[0].inputs = append(commands[0].inputs, stdin)
	}

	stdout := &buff{}
	stderr := &buff{}
	for i := range commands {
		commands[i].env = env
		commands[i].stderr = stderr
		if len(commands[i].outputs) == 0 {
			commands[i].outputs = append(commands[i].outputs, stdout)
		}
		err := commands[i].run(dir)
		if err != nil {
			return nil, nil, err
		}
	}

	return stdout.Bytes(), stderr.Bytes(), nil
}

func parse(s, dir string) []command {
	commands := []command{}
	for _, c := range strings.Split(s, "&&") {
		pipelineCommands := []command{}
		for _, p := range strings.Split(strings.TrimSpace(c), "|") {
			pipelineCommands = append(pipelineCommands, parseCommand(p, dir))
		}
		for i := 0; i < len(pipelineCommands)-1; i++ {
			b := buff{}
			pipelineCommands[i].outputs = append(pipelineCommands[i].outputs, &b)
			pipelineCommands[i+1].inputs = append([]io.ReadCloser{&b}, pipelineCommands[i+1].inputs...)

		}
		commands = append(commands, pipelineCommands...)
	}

	return commands
}

func parseCommand(s, dir string) command {
	splits := strings.Split(strings.TrimSpace(s), " ")
	cmd := command{
		cmd: strings.TrimSpace(splits[0]),
	}
	if len(splits) > 1 {
		parseArgs(&cmd, splits[1:], dir)
	}

	return cmd
}

func parseArgs(cmd *command, args []string, dir string) {
	cleaned := []string{}
	for i := range args {
		a := strings.TrimSpace(args[i])
		if a == "" {
			continue
		}
		cleaned = append(cleaned, a)
	}

	for i := 0; i < len(cleaned); i++ {
		if cleaned[i] == ">" {
			if len(cleaned) > i+1 {
				cmd.outputs = append(cmd.outputs, &outputFile{join(dir, cleaned[i+1]), nil})
				i++
			}
			continue
		}
		if cleaned[i] == "<" {
			if len(cleaned) > i+1 {
				cmd.inputs = append(cmd.inputs, &inputFile{join(dir, cleaned[i+1]), nil})
				i++
			}
			continue
		}

		cmd.args = append(cmd.args, cleaned[i])
	}
}

func join(dir, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(dir, filename)
}
