// The kubedebug binary provides the means to interactively run "kubectl debug"
package main

import (
	"context"
	"log"
	"os"

	"github.com/paulburlumi/kubedebug/internal/kubedebug"
)

func main() {

	kd := kubedebug.NewKubeDebug(os.Args, os.Stderr, kubedebug.NewCommand(os.Stdin, os.Stdout, os.Stderr))
	if err := kd.Run(context.Background()); err != nil {
		log.Fatalf("kubedebug failed: %v", err)
	}
}
