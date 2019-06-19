// +build !main

package completion

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

//go:generate aduu generate

// RunCompletion checks given arguments and executes command
func RunCompletion(out io.Writer, boilerPlate string, cmd *cobra.Command, args []string) error {
	run, found := completionShells[args[0]]
	if !found {
		return fmt.Errorf("Unsupported shell type %q", args[0])
	}

	return run(out, boilerPlate, cmd.Parent())
}

func GenerateShellCompletion(cmd *cobra.Command, boilerPlate string, shell string) {
	if cmd.Parent() == nil {
		panic("give one command under root, not root itself")
	}
	if cmd.Parent() != cmd.Root() {
		panic("parent != root")
	}

	var buffer bytes.Buffer



	if err := RunCompletion(&buffer, boilerPlate, cmd, []string{shell}); err != nil {
		panic(err)
	}

	out := strings.ReplaceAll(buffer.String(), "kubectl", cmd.Parent().Name())
	out = strings.ReplaceAll(out, "_bash_comp", "_bash_complete")
	fmt.Print(out)
}

func NewCompletionCMD() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "completion SHELL",
		Short: "completions for zsh or bash hopefully.",
		Run: func(cmd *cobra.Command, args []string) {
			boilerPlate := `# Shell completion for ` + cmd.Parent().Name() + "."
			GenerateShellCompletion(cmd, boilerPlate, args[0])
		},
		ValidArgs: []string{"zsh", "bash"},
		Args:      cobra.ExactValidArgs(1),
	}

	return cmd
}


/*
run, found := completionShells[args[0]]
if !found {
return cmdutil.UsageErrorf(cmd, "Unsupported shell type %q.", args[0])
}

return run(out, boilerPlate, cmd.Parent())

*/