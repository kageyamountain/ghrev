package twoapprove

import (
	"flag"

	"github.com/kageyamountain/ghrev/internal/common/globaloption"
	"github.com/kageyamountain/ghrev/internal/subcommand"
)

type RuntimeOptions struct {
	RepositoryOwner globaloption.RepositoryOwner
	RepositoryName  globaloption.RepositoryName
	CreatedAtFrom   globaloption.CreatedAtFrom
	CreatedAtTo     globaloption.CreatedAtTo
	IgnoreLabels    globaloption.IgnoreLabels
}

func NewRuntimeOptions(optionArgs []string) (*RuntimeOptions, error) {
	flagSet := flag.NewFlagSet(subcommand.TwoApprove.String(), flag.ContinueOnError)
	repositoryOwnerFlag := flagSet.String(globaloption.OptionNameRepositoryOwner, "", "target repository owner")
	repositoryNameFlag := flagSet.String(globaloption.OptionNameRepositoryName, "", "target repository name")
	createdAtFromFlag := flagSet.String(globaloption.OptionNameCreatedAtFrom, "", "pull request's created at from")
	createdAtToFlag := flagSet.String(globaloption.OptionNameCreatedAtTo, "", "pull request's created at to")
	ignoreLabelsFlag := flagSet.String(globaloption.OptionNameIgnoreLabels, "", "ignore labels")
	err := flagSet.Parse(optionArgs)
	if err != nil {
		return nil, err
	}

	repositoryOwner, err := globaloption.ParseRepositoryOwner(*repositoryOwnerFlag)
	if err != nil {
		return nil, err
	}

	repositoryName, err := globaloption.ParseRepositoryName(*repositoryNameFlag)
	if err != nil {
		return nil, err
	}

	createdAtFrom, err := globaloption.ParseCreatedAtFrom(*createdAtFromFlag)
	if err != nil {
		return nil, err
	}

	createdAtTo, err := globaloption.ParseCreatedAtTo(*createdAtToFlag)
	if err != nil {
		return nil, err
	}

	ignoreLabels, err := globaloption.ParseIgnoreLabels(*ignoreLabelsFlag)
	if err != nil {
		return nil, err
	}

	return &RuntimeOptions{
		RepositoryOwner: repositoryOwner,
		RepositoryName:  repositoryName,
		CreatedAtFrom:   createdAtFrom,
		CreatedAtTo:     createdAtTo,
		IgnoreLabels:    ignoreLabels,
	}, nil
}
