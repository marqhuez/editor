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
			break
		}

		c = buf[0]

		fmt.Println(buf, string(buf[:]))

		if c == 'q' {
			break
		}
	}

	disableRawMode(oldState)
}

func enterRawMode() (*state, error) {
	fd := int(os.Stdin.Fd())

	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		panic(err)
	}

	oldState := &state{termios: *termios}

	termios.Iflag &^= unix.IXON
	termios.Iflag &^= unix.IEXTEN
	termios.Iflag &^= unix.ICRNL

	termios.Lflag &^= unix.ECHO
	termios.Lflag &^= unix.ICANON
	termios.Lflag &^= unix.ISIG

	unix.IoctlSetTermios(fd, unix.TCSETS, termios)

	return oldState, nil
}

func disableRawMode(state *state) {
	fd := int(os.Stdin.Fd())
	unix.IoctlSetTermios(fd, unix.TCSETS, &state.termios)
}
