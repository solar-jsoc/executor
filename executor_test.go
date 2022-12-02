package executor

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommand(t *testing.T) {
	cmd := parseCommand("ps", "")
	require.Equal(t, "ps", cmd.cmd)
	require.Equal(t, 0, len(cmd.args))

	cmd = parseCommand("echo test", "")
	require.Equal(t, "echo", cmd.cmd)
	require.Equal(t, []string{"test"}, cmd.args)

	cmd = parseCommand("echo test >   out.txt", "")
	require.Equal(t, "echo", cmd.cmd)
	require.Equal(t, []string{"test"}, cmd.args)
	require.Equal(t, 1, len(cmd.outputs))
}

func TestJoin(t *testing.T) {
	require.Equal(t, "/var/1.txt", join("/var", "1.txt"))
	require.Equal(t, "/1.txt", join("/var", "/1.txt"))
}

func TestParseArgs(t *testing.T) {
	cmd := command{
		cmd:     "",
		args:    []string{},
		inputs:  []io.ReadCloser{},
		outputs: []io.WriteCloser{},
		stderr:  nil,
	}
	args := []string{" ", " <", "> ", "< >", ">   <"}
	dir := "test_dir"
	parseArgs(&cmd, args, dir)
	require.Equal(t, len(cmd.args), 2)
	require.Equal(t, len(cmd.inputs), 1)
}

func TestParse(t *testing.T) {
	dir := "test_dit"
	commands := parse("echo 1", dir)
	require.Equal(t, 1, len(commands))
	require.Equal(t, "echo", commands[0].cmd)
	require.Equal(t, []string{"1"}, commands[0].args)

	commands = parse("echo 1 && cat test.txt", dir)
	require.Equal(t, 2, len(commands))
	require.Equal(t, "echo", commands[0].cmd)
	require.Equal(t, []string{"1"}, commands[0].args)
	require.Equal(t, "cat", commands[1].cmd)
	require.Equal(t, []string{"test.txt"}, commands[1].args)

	commands = parse("echo 2 | cat", dir)
	require.Equal(t, 2, len(commands))
	require.Equal(t, "echo", commands[0].cmd)
	require.Equal(t, []string{"2"}, commands[0].args)
	require.Equal(t, "cat", commands[1].cmd)
	require.Equal(t, 0, len(commands[1].args))
}

func TestExec(t *testing.T) {
	r := io.NopCloser(strings.NewReader(""))
	_, _, err := Exec("echo 1 && cat file.go", "", nil, r)
	require.Equal(t, err.Error(), "exit status 1")
	_, _, err = Exec("echo 1 && cat file.go", "./", nil, r)
	require.Nil(t, err)
}

func TestRun(t *testing.T) {
	commands := parse("echo 2 | cat", "")
	cmd := commands[0]
	err := cmd.run("")
	require.Nil(t, err)
	cmd = commands[1]
	err = cmd.run("")
	require.Nil(t, err)
	cmd.cmd = "echoo"
	err = cmd.run("")
	require.NotNil(t, err)
}
