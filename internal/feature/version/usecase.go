package version

import (
	"context"
	"fmt"
)

type UseCase struct {
	version string
}

func NewUseCase(version string) *UseCase {
	return &UseCase{
		version: version,
	}
}

func (u *UseCase) Do(ctx context.Context) error {
	fmt.Println(u.version)
	return nil
}
