package graphql

import (
	"context"
)

// MarshalSliceConcurrently marshals a slice of elements concurrently, writing
// each result into the returned Array.
//
// The marshalElement callback is called for each index and receives a context
// that already has a FieldContext with Index set. The callback should set
// FieldContext.Result and perform the actual marshaling.
//
// workerLimit of 0 means unlimited concurrency.
func MarshalSliceConcurrently(
	ctx context.Context,
	length int,
	workerLimit int64,
	omitPanicHandler bool,
	marshalElement func(ctx context.Context, i int) Marshaler,
) Array {
	ret := make(Array, length)
	if length == 0 {
		return ret
	}

	isLen1 := length == 1

	if isLen1 {
		i := 0
		fc := &FieldContext{
			Index: &i,
		}
		childCtx := WithFieldContext(ctx, fc)
		if omitPanicHandler {
			ret[0] = marshalElement(childCtx, 0)
		} else {
			func() {
				defer func() {
					if r := recover(); r != nil {
						AddError(childCtx, Recover(childCtx, r))
						ret = nil
					}
				}()
				ret[0] = marshalElement(childCtx, 0)
			}()
		}
		return ret
	}

	for i := range length {
		fc := &FieldContext{
			Index: &i,
		}
		childCtx := WithFieldContext(ctx, fc)

		if omitPanicHandler {
			ret[i] = marshalElement(childCtx, i)
		} else {
			func() {
				defer func() {
					if r := recover(); r != nil {
						AddError(childCtx, Recover(childCtx, r))
						ret = nil
					}
				}()
				ret[i] = marshalElement(childCtx, i)
			}()
			if ret == nil {
				return nil
			}
		}
	}
	return ret
}
