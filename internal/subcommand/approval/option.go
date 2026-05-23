package approval

import (
	"flag"
	"time"

	"github.com/kageyamountain/ghrev/internal/common/globaloption"
	"github.com/kageyamountain/ghrev/internal/subcommand"
)

type RuntimeOptions struct {
	Owner             string
	Name              string
	CreatedAtFrom     time.Time
	CreatedAtTo       time.Time
	IgnoreLabels      []string
	Assignees         []string
	RequiredApprovals int
}

func NewRuntimeOptions(optionArgs []string) (*RuntimeOptions, error) {
	flagSet := flag.NewFlagSet(subcommand.Approval.String(), flag.ContinueOnError)
	repositoryOwnerFlag := flagSet.String(globaloption.OptionNameRepositoryOwner, "", "target repository owner")
	repositoryNameFlag := flagSet.String(globaloption.OptionNameRepositoryName, "", "target repository name")
	createdAtFromFlag := flagSet.String(globaloption.OptionNameCreatedAtFrom, "", "pull request's created at from")
	createdAtToFlag := flagSet.String(globaloption.OptionNameCreatedAtTo, "", "pull request's created at to")
	ignoreLabelsFlag := flagSet.String(globaloption.OptionNameIgnoreLabels, "", "ignore labels")
	assigneesFlag := flagSet.String(globaloption.OptionNameAssignees, "", "filter PRs by assignees")
	requiredApprovalsFlag := flagSet.Int(OptionNameRequiredApprovals, 0, "number of approvals required to complete review")
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

	assignees, err := globaloption.ParseAssignees(*assigneesFlag)
	if err != nil {
		return nil, err
	}

	requiredApprovals, err := ParseRequiredApprovals(*requiredApprovalsFlag)
	if err != nil {
		return nil, err
	}

	return &RuntimeOptions{
		Owner:             repositoryOwner.String(),
		Name:              repositoryName.String(),
		CreatedAtFrom:     createdAtFrom.Time(),
		CreatedAtTo:       createdAtTo.Time(),
		IgnoreLabels:      ignoreLabels.Strings(),
		Assignees:         assignees.Strings(),
		RequiredApprovals: requiredApprovals.Int(),
	}, nil
}
