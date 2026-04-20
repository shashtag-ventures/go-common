package gitutil

import (
	"context"

	"github.com/shashtag-ventures/go-common/executil"
)

// GitCloner defines the interface for cloning git repositories.
type GitCloner interface {
	Clone(ctx context.Context, repoURL, targetPath string) ([]byte, error)
}

// DefaultGitCloner is the standard implementation of GitCloner using the git CLI.
type DefaultGitCloner struct{}

// NewDefaultGitCloner creates a new DefaultGitCloner.
func NewDefaultGitCloner() *DefaultGitCloner {
	return &DefaultGitCloner{}
}

// Clone performs a shallow clone of a git repository to the specified target path.
func (g *DefaultGitCloner) Clone(ctx context.Context, repoURL, targetPath string) ([]byte, error) {
	cmd := executil.CommandContext(ctx, "git", "clone", "--depth", "1", "--", repoURL, targetPath)
	return cmd.CombinedOutput()
}
