// +build darwin freebsd netbsd openbsd
// +build !plan9,!windows,!linux

package flow

import (
	"os"
	"os/signal"
	"syscall"
)

// OnInterrupt subscribe fn to the os.Signal notifications e.g: os.Interrupt
// It returns a detach function to unsubscribe from the notifications.
func OnInterrupt(fn func(), onExitFunc func()) (detach func()) {
	// deal with control+c,etc
	signalChan := make(chan os.Signal, 1)
	// controlling terminal close, daemon not exit
	signal.Ignore(syscall.SIGHUP)
	signal.Notify(signalChan,
		os.Interrupt,
		os.Kill,
		syscall.SIGALRM,
		// syscall.SIGHUP,
		syscall.SIGINFO, // this causes windows to fail
		syscall.SIGINT,
		syscall.SIGTERM,
		// syscall.SIGQUIT, // Quit from keyboard, "kill -3"
	)
	go func() {
		for sig := range signalChan {
			fn()
			if sig != syscall.SIGINFO {
				if onExitFunc != nil {
					onExitFunc()
				}
				os.Exit(0)
			}
		}
	}()

	return func() {
		signal.Stop(signalChan)
		close(signalChan)
	}
}
