package help

import (
	"context"
	"fmt"
)

type UseCase struct{}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (u *UseCase) Do(ctx context.Context) error {
	fmt.Println("help is not implemented yet")
	return nil
}
