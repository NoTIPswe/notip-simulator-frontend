package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start an interactive shell session",
	Long: `Opens a persistent REPL where you can run multiple commands without
restarting the container. Type 'help' for available commands, 'exit' to quit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printWelcomeBanner()

		// Silence cobra's usage print on every error inside the REPL —
		// errors are shown by the individual commands themselves.
		rootCmd.SilenceUsage = true
		defer func() { rootCmd.SilenceUsage = false }()

		if canUseLineEditor() {
			err := runShellWithLineEditor()
			if err == nil {
				return nil
			}
			pterm.Warning.Printf("line editor disabled: %v\n", err)
		}

		reader := bufio.NewReader(os.Stdin)
		for {
			printPrompt()

			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println()
					pterm.Info.Println("Goodbye!")
				}
				return nil
			}

			if shouldExit := processShellLine(line); shouldExit {
				return nil
			}
		}
	},
}

func canUseLineEditor() bool {
	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd()))
}

func runShellWithLineEditor() error {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer func() {
		_ = term.Restore(fd, oldState)
	}()

	lineEditor := term.NewTerminal(struct {
		io.Reader
		io.Writer
	}{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}, "sim-cli> ")

	for {
		line, readErr := lineEditor.ReadLine()
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				fmt.Println()
				pterm.Info.Println("Goodbye!")
				return nil
			}
			return readErr
		}

		// Leave raw mode while the command runs so output renders correctly.
		if err := term.Restore(fd, oldState); err != nil {
			return err
		}
		shouldExit := processShellLine(line)
		if _, err := term.MakeRaw(fd); err != nil {
			return err
		}
		if shouldExit {
			return nil
		}
	}
}

func processShellLine(line string) bool {
	line = strings.TrimSpace(line)
	if line == "" {
		return false
	}
	if line == "exit" || line == "quit" {
		pterm.Info.Println("Goodbye!")
		return true
	}

	// Prevent the user from nesting shells.
	args := strings.Fields(line)
	if args[0] == "shell" {
		pterm.Warning.Println("Already inside a shell session.")
		return false
	}

	resetAllCommandFlags(rootCmd)
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		pterm.Error.Println(err)
	}

	return false
}

func printPrompt() {
	if pterm.RawOutput {
		fmt.Print("sim-cli> ")
	} else {
		fmt.Print(pterm.Green("sim-cli") + pterm.Gray(" › ") + " ")
	}
}

func printWelcomeBanner() {
	if !pterm.RawOutput {
		if err := pterm.DefaultBigText.WithLetters(
			putils.LettersFromStringWithStyle("sim", pterm.NewStyle(pterm.FgGreen)),
			putils.LettersFromStringWithStyle("-cli", pterm.NewStyle(pterm.FgGray)),
		).Render(); err != nil {
			pterm.Warning.Printf("failed to render banner: %v\n", err)
		}
	}
	pterm.DefaultBox.
		WithTitle("Interactive Shell").
		WithTitleTopCenter().
		Println("Type a command to execute it.\n" +
			pterm.Gray("  gateways list\n") +
			pterm.Gray("  gateways get <uuid>\n") +
			pterm.Gray("  sensors list <gateway-id-or-uuid>\n") +
			pterm.Gray("  anomalies disconnect <uuid> --duration 5\n") +
			"\nType " + pterm.Bold.Sprint("exit") + " or press Ctrl+D to quit.")
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
