package graphql

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// withTestResponseContext sets up a minimal response context for testing.
func withTestResponseContext(ctx context.Context) context.Context {
	return WithResponseContext(ctx, func(ctx context.Context, err error) *gqlerror.Error {
		return &gqlerror.Error{Message: err.Error()}
	}, DefaultRecover)
}

func TestMarshalSliceConcurrently(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		ctx := withTestResponseContext(context.Background())
		ret := MarshalSliceConcurrently(
			ctx,
			0,
			0,
			false,
			func(ctx context.Context, i int) Marshaler {
				t.Fatal("should not be called")
				return Null
			},
		)
		assert.Empty(t, ret)
	})

	t.Run("single element", func(t *testing.T) {
		ctx := withTestResponseContext(context.Background())
		var called bool
		ret := MarshalSliceConcurrently(
			ctx,
			1,
			0,
			false,
			func(ctx context.Context, i int) Marshaler {
				called = true
				assert.Equal(t, 0, i)
				fc := GetFieldContext(ctx)
				require.NotNil(t, fc)
				assert.Equal(t, 0, *fc.Index)
				return MarshalString("hello")
			},
		)
		assert.True(t, called)
		require.Len(t, ret, 1)
		var buf bytes.Buffer
		ret[0].MarshalGQL(&buf)
		assert.Equal(t, `"hello"`, buf.String())
	})

	t.Run("multiple elements", func(t *testing.T) {
		ctx := withTestResponseContext(context.Background())
		n := 10
		callCount := 0
		ret := MarshalSliceConcurrently(
			ctx,
			n,
			0,
			false,
			func(ctx context.Context, i int) Marshaler {
				callCount++
				fc := GetFieldContext(ctx)
				require.NotNil(t, fc)
				assert.Equal(t, i, *fc.Index)
				return MarshalString(fmt.Sprintf("item-%d", i))
			},
		)
		assert.Equal(t, n, callCount)
		require.Len(t, ret, n)
		for i := range n {
			var buf bytes.Buffer
			ret[i].MarshalGQL(&buf)
			assert.Equal(t, fmt.Sprintf(`"item-%d"`, i), buf.String())
		}
	})

	t.Run("worker limit parameter accepted but sequential", func(t *testing.T) {
		ctx := withTestResponseContext(context.Background())
		n := 20
		var workerLimit int64 = 3

		ret := MarshalSliceConcurrently(
			ctx,
			n,
			workerLimit,
			false,
			func(ctx context.Context, i int) Marshaler {
				return MarshalString(fmt.Sprintf("item-%d", i))
			},
		)

		require.Len(t, ret, n)
	})

	t.Run("panic recovery sets result to nil", func(t *testing.T) {
		ctx := withTestResponseContext(context.Background())
		ret := MarshalSliceConcurrently(
			ctx,
			1,
			0,
			false,
			func(ctx context.Context, i int) Marshaler {
				panic("test panic")
			},
		)
		assert.Nil(t, ret)
	})

	t.Run("panic recovery in multi-element slice sets result to nil", func(t *testing.T) {
		ctx := withTestResponseContext(context.Background())
		ret := MarshalSliceConcurrently(
			ctx,
			3,
			0,
			false,
			func(ctx context.Context, i int) Marshaler {
				if i == 1 {
					panic("test panic")
				}
				return MarshalString("ok")
			},
		)
		assert.Nil(t, ret)
	})

	t.Run("omit panic handler does not recover", func(t *testing.T) {
		ctx := withTestResponseContext(context.Background())
		assert.Panics(t, func() {
			MarshalSliceConcurrently(ctx, 1, 0, true, func(ctx context.Context, i int) Marshaler {
				panic("test panic")
			})
		})
	})

	t.Run("cancelled context still works without worker limit", func(t *testing.T) {
		ctx := withTestResponseContext(context.Background())
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		ret := MarshalSliceConcurrently(
			ctx,
			5,
			0,
			false,
			func(ctx context.Context, i int) Marshaler {
				return MarshalString("ok")
			},
		)
		require.Len(t, ret, 5)
	})
}
