// +build !android

package main

import (
	"flag"
	"github.com/remogatto/application"
	"github.com/remogatto/gorgasm"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// sigterm is a type for handling a SIGTERM signal.
type sigterm int

func (h *sigterm) HandleSignal(s os.Signal) {
	switch ss := s.(type) {
	case syscall.Signal:
		switch ss {
		case syscall.SIGTERM, syscall.SIGINT:
			application.Exit()
		}
	}
}

func main() {
	var signal sigterm

	verbose := flag.Bool("verbose", false, "produce verbose output")
	debug := flag.Bool("debug", false, "produce debug output")
	size := flag.String("size", "320x480", "set the size of the window")

	flag.Parse()

	if *verbose {
		application.Verbose = true
	}

	if *debug {
		application.Debug = true
	}

	dims := strings.Split(strings.ToLower(*size), "x")
	width, err := strconv.Atoi(dims[0])
	if err != nil {
		panic(err)
	}
	height, err := strconv.Atoi(dims[1])
	if err != nil {
		panic(err)
	}

	// Enable CTRL-C shortcut to kill the application
	application.InstallSignalHandler(&signal)

	// Initialize EGL for xorg
	gorgasm.XorgInitialize(width, height)

	go application.Run()
	for {
		select {
		case eglState := <-gorgasm.Init:
			renderLoop := newRenderLoop(eglState, FRAMES_PER_SECOND)
			eventsLoop := newEventsLoop(renderLoop)
			application.Register("renderLoop", renderLoop)
			application.Register("eventsLoop", eventsLoop)
			application.Start("renderLoop")
			application.Start("eventsLoop")
		case <-application.ExitCh:
			return
		case err := <-application.ErrorCh:
			application.Logf(err.(application.Error).Error())
		}
	}
}
