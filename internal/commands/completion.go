package commands

import (
	"github.com/spf13/cobra"
)

func NewCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:
  $ source <(k8ctl completion bash)
  
  # To load completions for each session, execute once:
  # Linux:
  $ k8ctl completion bash > /etc/bash_completion.d/k8ctl
  # macOS:
  $ k8ctl completion bash > /usr/local/etc/bash_completion.d/k8ctl

Zsh:
  $ source <(k8ctl completion zsh)
  
  # To load completions for each session, execute once:
  $ k8ctl completion zsh > "${fpath[1]}/_k8ctl"

Fish:
  $ k8ctl completion fish | source
  
  # To load completions for each session, execute once:
  $ k8ctl completion fish > ~/.config/fish/completions/k8ctl.fish

PowerShell:
  PS> k8ctl completion powershell | Out-String | Invoke-Expression
  
  # To load completions for each session, execute once:
  PS> k8ctl completion powershell > k8ctl.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletion(cmd.OutOrStdout())
		}
		return nil
	}

	return cmd
}
