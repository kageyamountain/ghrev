package runtimeoption

import "flag"

type Options struct {
	Mode            Mode
	RepositoryOwner RepositoryOwner
	RepositoryName  RepositoryName
	CreatedAtFrom   CreatedAtFrom
	CreatedAtTo     CreatedAtTo
	IgnoreLabels    []string
}

func NewOptions() (*Options, error) {
	modeFlag := flag.String(OptionNameMode, ModeTwoApprove.String(), "execution mode")
	repositoryOwnerFlag := flag.String(OptionNameRepositoryOwner, "", "target repository owner")
	repositoryNameFlag := flag.String(OptionNameRepositoryName, "", "target repository name")
	createdAtFromFlag := flag.String(OptionNameCreatedAtFrom, "", "pull request's created at from")
	createdAtToFlag := flag.String(OptionNameCreatedAtTo, "", "pull request's created at to")
	ignoreLabelsFlag := flag.String(OptionNameIgnoreLabels, "", "ignore labels")
	flag.Parse()

	mode, err := ParseMode(*modeFlag)
	if err != nil {
		return &Options{}, err
	}

	repositoryOwner, err := ParseRepositoryOwner(*repositoryOwnerFlag)
	if err != nil {
		return &Options{}, err
	}

	repositoryName, err := ParseRepositoryName(*repositoryNameFlag)
	if err != nil {
		return &Options{}, err
	}

	createdAtFrom, err := ParseCreatedAtFrom(*createdAtFromFlag)
	if err != nil {
		return &Options{}, err
	}

	createdAtTo, err := ParseCreatedAtTo(*createdAtToFlag)
	if err != nil {
		return &Options{}, err
	}

	ignoreLabels, err := ParseIgnoreLabel(*ignoreLabelsFlag)
	if err != nil {
		return &Options{}, err
	}

	return &Options{
		Mode:            mode,
		RepositoryOwner: repositoryOwner,
		RepositoryName:  repositoryName,
		CreatedAtFrom:   createdAtFrom,
		CreatedAtTo:     createdAtTo,
		IgnoreLabels:    ignoreLabels,
	}, nil
}
