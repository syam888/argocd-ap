package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"argocd-ap/internal/appstate"
	"argocd-ap/internal/argo"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	config := appstate.New()
	client := argo.New(os.Getenv("ARGOCD_PROJECT"))

	if len(args) == 0 {
		return selectApp(config, client)
	}

	cmd := args[0]
	passThru := []string{}
	if len(args) > 1 {
		passThru = args[1:]
	}

	switch cmd {
	case "select":
		return selectApp(config, client)
	case "current":
		app, err := config.Get()
		if err != nil {
			return err
		}
		fmt.Println(app)
		return nil
	case "clear":
		return config.Clear()
	case "help", "-h", "--help":
		printUsage()
		return nil
	case "get", "diff", "sync", "history", "resources", "manifests":
		app, err := config.Get()
		if err != nil {
			return err
		}
		return client.RunAppCommand(cmd, app, passThru)
	default:
		app, err := config.Get()
		if err != nil {
			return err
		}
		return client.RunAppCommand(cmd, app, passThru)
	}
}

func selectApp(config *appstate.Config, client *argo.Client) error {
	appNames, err := client.AppList()
	if err != nil {
		return err
	}
	if len(appNames) == 0 {
		return fmt.Errorf("no applications found")
	}

	var selected string

	// Try to use fzf if available
	_, err = exec.LookPath("fzf")
	if err == nil {
		cmd := exec.Command("fzf", "--prompt=Select Argo CD app: ")
		cmd.Stderr = os.Stderr
		cmd.Stdin = strings.NewReader(strings.Join(appNames, "\n"))
		out, err := cmd.Output()
		if err == nil {
			selected = strings.TrimSpace(string(out))
		}
	}

	// Fall back to numbered menu if fzf failed or is unavailable
	if selected == "" {
		fmt.Println("Select an Argo CD application:")
		for i, name := range appNames {
			fmt.Printf("%d. %s\n", i+1, name)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter number or app name: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read selection: %w", err)
		}
		input = strings.TrimSpace(input)

		if idx, err := strconv.Atoi(input); err == nil {
			if idx < 1 || idx > len(appNames) {
				return fmt.Errorf("invalid selection")
			}
			selected = appNames[idx-1]
		} else {
			selected = input
		}
	}

	if selected == "" {
		return fmt.Errorf("no app selected")
	}

	if err := config.Set(selected); err != nil {
		return fmt.Errorf("failed to save selection: %w", err)
	}

	fmt.Printf("Selected Argo CD app: %s\n", selected)
	return nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  ap                    Select an app")
	fmt.Println("  ap select             Select an app")
	fmt.Println("  ap current            Show the selected app")
	fmt.Println("  ap clear              Clear the selected app")
	fmt.Println("  ap get [args]         Run 'argocd app get <selected-app> [args]'")
	fmt.Println("  ap diff [args]        Run 'argocd app diff <selected-app> [args]'")
	fmt.Println("  ap sync [args]        Run 'argocd app sync <selected-app> [args]'")
	fmt.Println("  ap history [args]     Run 'argocd app history <selected-app> [args]'")
	fmt.Println("  ap resources [args]   Run 'argocd app resources <selected-app> [args]'")
	fmt.Println("  ap manifests [args]   Run 'argocd app manifests <selected-app> [args]'")
}
