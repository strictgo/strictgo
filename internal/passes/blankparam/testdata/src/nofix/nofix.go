package nofix

import (
	"context"
	"testing"
)

var ctx = context.Background()

func duplicate(_, _ context.Context) {} // want "blank parameter obscures its purpose" "blank parameter obscures its purpose"

func existing(_ context.Context, ctx context.Context) {} // want "blank parameter obscures its purpose"

func result(_ context.Context) (ctx context.Context) { return nil } // want "blank parameter obscures its purpose"

func local(_ context.Context) { // want "blank parameter obscures its purpose"
	ctx := context.Background()
	_ = ctx
}

func captured(_ context.Context) { // want "blank parameter obscures its purpose"
	_ = ctx
}

func nested(_ context.Context) { // want "blank parameter obscures its purpose"
	callback := func() { _ = ctx }
	callback()
}

type record struct{}

func (ctx record) method(_ context.Context) {} // want "blank parameter obscures its purpose"

type Type struct{}

func keyword(_ Type) {} // want "blank parameter obscures its purpose"

func commonCollision(_ testing.TB, tb testing.TB) {} // want "blank parameter obscures its purpose"
