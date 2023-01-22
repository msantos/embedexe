package fdexec_test

import (
	"log"
	"os"

	"codeberg.org/msantos/embedexe"
	"codeberg.org/msantos/embedexe/fdexec"
)

func ExampleCommand() {
	exe, err := os.ReadFile("/bin/echo")
	if err != nil {
		log.Fatalln(err)
	}

	fd, err := embedexe.Open(exe, "echo")
	if err != nil {
		log.Fatalln(err)
	}

	cmd := fdexec.Command(fd, "test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}

	// Output: test
}
