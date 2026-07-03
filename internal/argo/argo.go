package argo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Client struct {
	Project string
}

func New(project string) *Client {
	return &Client{Project: project}
}

func (c *Client) AppList() ([]string, error) {
	args := []string{"app", "list", "-o", "name"}
	if c.Project != "" {
		args = []string{"app", "list", "--project", c.Project, "-o", "name"}
	}

	cmd := exec.Command("argocd", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("argocd app list failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}

func (c *Client) RunAppCommand(subcommand, app string, passThru []string) error {
	args := []string{"app", subcommand, app}
	args = append(args, passThru...)
	cmd := exec.Command("argocd", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
