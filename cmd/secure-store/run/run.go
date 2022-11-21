package run

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func ExecuteArgs() {
	// Execute the main process requested by the user. It is run here as this app needs to
	// keep running to serve the encryption, and to allow both processes to remain connected to
	// PID 1

	if len(flag.Args()) == 0 {
		fmt.Println("No arguments were passed to execute.")
		return
	}

	fmt.Println("Executing the passed arguments...")

	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)

	// Connect the output of the requested programme to the output of this programme running as PID 1
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		// Raise non-zero exit code to ensure Docker's restart on failure policy works
		log.Fatal(err)
	}
}
