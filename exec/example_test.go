package exec_test

import (
	"log"
	"os"

	"codeberg.org/msantos/embedexe/exec"
)

func ExampleCommand() {
	exe, err := os.ReadFile("/bin/echo")
	if err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command(exe, []string{"test"})
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}

	// Output: test
}
