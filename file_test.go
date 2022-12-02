package executor

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	content := []byte(`Lorem ipsum dolor sit amet, consectetur 
	adipiscing elit, sed do eiusmod tempor incididunt ut labore et 
	dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation 
	ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor 
	in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. 
	Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`)

	tmpfile, err := os.CreateTemp("./", "example")
	require.Nil(t, err)
	defer os.Remove(tmpfile.Name()) // clean up

	tmpfile.Write(content)

	tmpfile.Close()

	f := inputFile{
		path: "error_path",
		file: nil,
	}

	buf := make([]byte, len(content))

	_, err1 := f.Read(buf)
	require.NotNil(t, err1)

	f.path = tmpfile.Name()
	n2, _ := f.Read(buf)
	require.Equal(t, n2, 455)

	n3, err3 := f.Read(buf)
	require.Equal(t, n3, 0)
	require.Equal(t, err3.Error(), "EOF")
	f.Close()
}

func TestWrite(t *testing.T) {
	tmpfile, err := os.CreateTemp("./", "example")
	require.Nil(t, err)
	defer os.Remove(tmpfile.Name()) // clean up

	f := outputFile{
		path: tmpfile.Name(),
		file: tmpfile,
	}
	buf := []byte("Lorem ipsum dolor sit amet, consectetur")
	n, err := f.Write(buf)
	require.Equal(t, n, len(buf))
	require.Nil(t, err)
	f.Close()
}
