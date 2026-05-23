package help

import (
	"context"
	"fmt"
)

const helpText = `ghrev - A CLI tool for measuring GitHub Pull Request review metrics.

Usage:
  ghrev <subcommand> [options]

For details on subcommands and options, see:
  https://github.com/kageyamountain/ghrev
`

type UseCase struct{}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (u *UseCase) Do(ctx context.Context) error {
	fmt.Print(helpText)
	return nil
}
