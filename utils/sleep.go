package utils

import (
	"context"
	"fmt"
	"time"
)

func Sleep(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context expired")
	case <-time.After(d):
		return nil
	}
}
