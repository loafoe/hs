package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestS3CredsCmd(t *testing.T) {
	cmd := s3credsCmd
	if cmd.Use != "s3creds" {
		t.Errorf("expected command use to be 's3creds', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected command to have a short description")
	}

	if cmd.Run == nil {
		t.Error("expected command to have a Run function")
	}
}

func TestS3CredsCmdExecution(t *testing.T) {
	cmd := s3credsCmd
	cmd.Run = func(cmd *cobra.Command, args []string) {
		// Simulate command execution
	}

	err := cmd.Execute()
	if err != nil {
		t.Errorf("command execution failed: %v", err)
	}
}
