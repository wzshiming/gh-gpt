package test

import (
	"context"
	"os"
	"testing"

	"github.com/wzshiming/gh-gpt/pkg/run"
)

func TestA(t *testing.T) {
	err := example()
	if err != nil {
		t.Fatal(err)
	}
}

func example() error {
	ctx := context.Background()

	err := run.RunStream(ctx, "Who are you?", os.Stdout)
	if err != nil {
		return err
	}
	return nil
}
