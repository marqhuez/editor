package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

type state struct {
	screenRows  int
	screenCols  int
	origTermios unix.Termios
}

var editorState *state = &state{}

func main() {
	enterRawMode()
	initEditor()

	fmt.Println(editorState.screenRows)

	for {
		editorClearScreen()
		editorProcessKeypress()
	}
}

func initEditor() {
	getWindowSize(&editorState.screenRows, &editorState.screenCols)
}

func getWindowSize(rows *int, cols *int) {
	size, err := unix.IoctlGetWinsize(int(os.Stdin.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		panic(err)
	}

	*rows = int(size.Row)
	*cols = int(size.Col)
}

func editorDrawRows() {
	for y := 0; y < editorState.screenRows; y++ {
		fmt.Fprintf(os.Stdout, "~\r\n")
	}
}

func editorClearScreen() {
	fmt.Fprintf(os.Stdout, "\x1b[2J")
	fmt.Fprintf(os.Stdout, "\x1b[H")

	editorDrawRows()

	fmt.Fprintf(os.Stdout, "\x1b[H")
}

func editorProcessKeypress() {
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

func enterRawMode() error {
	fd := int(os.Stdin.Fd())

	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		panic(err)
	}

	editorState.origTermios = *termios

	termios.Iflag &^= unix.IXON | unix.IEXTEN | unix.ICRNL
	termios.Lflag &^= unix.ECHO | unix.ICANON | unix.ISIG
	termios.Oflag &^= unix.OPOST

	// turn off other random flags bc of raw mode tradition
	termios.Iflag &^= unix.BRKINT | unix.INPCK | unix.ISTRIP
	termios.Cflag &^= unix.CS8

	unix.IoctlSetTermios(fd, unix.TCSETS, termios)

	return nil
}

func disableRawMode() {
	fd := int(os.Stdin.Fd())
	unix.IoctlSetTermios(fd, unix.TCSETS, &editorState.origTermios)
}

func exit(code int) {
	disableRawMode()
	fmt.Fprintf(os.Stdout, "\x1b[2J")
	fmt.Fprintf(os.Stdout, "\x1b[H")
	os.Exit(code)
}
