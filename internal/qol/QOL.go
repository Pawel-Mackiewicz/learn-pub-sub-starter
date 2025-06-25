package qol

import (
	"os"
	"os/signal"
)

func WaitForSignalToKill() {
	//wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
