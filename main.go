package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

type state struct {
	termios unix.Termios
}

func main() {
	oldState, err := enterRawMode()
	if err != nil {
		panic(err)
	}

	var c byte

	for {
		buf := make([]byte, 1)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			os.Exit(1)
		}

		c = buf[0]

		if c == ctrlKey('q') {
			break
		}

		fmt.Printf("%v\r\n", string(c))
	}

	disableRawMode(oldState)
}

func ctrlKey(k byte) byte {
	return k & 0x1f
}

func enterRawMode() (*state, error) {
	fd := int(os.Stdin.Fd())

	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		panic(err)
	}

	oldState := &state{termios: *termios}

	termios.Iflag &^= unix.IXON | unix.IEXTEN | unix.ICRNL
	termios.Lflag &^= unix.ECHO | unix.ICANON | unix.ISIG
	termios.Oflag &^= unix.OPOST

	// turn off other random flags bc of raw mode tradition
	termios.Iflag &^= unix.BRKINT | unix.INPCK | unix.ISTRIP
	termios.Cflag &^= unix.CS8

	unix.IoctlSetTermios(fd, unix.TCSETS, termios)

	return oldState, nil
}

func disableRawMode(state *state) {
	fd := int(os.Stdin.Fd())
	unix.IoctlSetTermios(fd, unix.TCSETS, &state.termios)
}
