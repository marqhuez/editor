package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

type state struct {
	termios unix.Termios
}

var oldState *state = &state{}

func main() {
	enterRawMode()

	for {
		processKeypress()
	}
}

func processKeypress() {
	c := readKey()

	switch c {
	case ctrlKey('q'):
		exit(0)
		break
	}
}

func readKey() byte {
	buf := make([]byte, 1)
	_, err := os.Stdin.Read(buf)
	if err != nil {
		fmt.Println("Error reading from stdin:", err)
		exit(1)
	}

	fmt.Printf("%v\r\n", string(buf[0]))
	return buf[0]
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

	oldState = &state{termios: *termios}

	termios.Iflag &^= unix.IXON | unix.IEXTEN | unix.ICRNL
	termios.Lflag &^= unix.ECHO | unix.ICANON | unix.ISIG
	termios.Oflag &^= unix.OPOST

	// turn off other random flags bc of raw mode tradition
	termios.Iflag &^= unix.BRKINT | unix.INPCK | unix.ISTRIP
	termios.Cflag &^= unix.CS8

	unix.IoctlSetTermios(fd, unix.TCSETS, termios)

	return oldState, nil
}

func disableRawMode() {
	fd := int(os.Stdin.Fd())
	unix.IoctlSetTermios(fd, unix.TCSETS, &oldState.termios)
}

func exit(code int) {
	disableRawMode()
	os.Exit(code)
}
