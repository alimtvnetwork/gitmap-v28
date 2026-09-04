package appfaults_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/appfaults"
	"coding-guidelines/common/pkg/errtype"
)

func TestCollectionBasicOperations(t *testing.T) {
	c := appfaults.New()
	if c.HasError() || !c.IsSuccess() || c.Count() != 0 {
		t.Fatal("expected empty collection to be success")
	}

	c.AddType(errtype.Validation).AddTypeMsg(errtype.Database, "db failed")
	c.AddError(errtype.IO, errors.New("io error")).AddErrorMsg(errtype.Network, errors.New("sock"), "conn drop")
	c.AddTypeMsgf(errtype.Timeout, "timed out after %ds", 5)

	if !c.HasError() || c.Count() != 5 {
		t.Fatalf("expected 5 items, got %d", c.Count())
	}
}

func TestCollectionNilAndNoneSafety(t *testing.T) {
	c := appfaults.New()
	c.Add(nil).AddType(errtype.None).AddTypeMsg(errtype.None, "skip")
	c.AddError(errtype.None, errors.New("skip")).AddError(errtype.IO, nil)
	c.AddErrorMsg(errtype.None, errors.New("skip"), "skip").AddErrorMsg(errtype.IO, nil, "skip")

	if c.HasError() || c.Count() != 0 {
		t.Fatalf("expected 0 items after adding nils and None variations, got %d", c.Count())
	}
}

func TestCollectionFilterAndTransform(t *testing.T) {
	c := appfaults.New()
	c.AddTypeMsg(errtype.Validation, "err1").AddTypeMsg(errtype.Database, "err2")

	filtered := c.FilterByType(errtype.Validation)
	if filtered.Count() != 1 || filtered.First().GetType() != errtype.Validation {
		t.Fatalf("expected 1 validation fault, got %d", filtered.Count())
	}

	comp := c.ToAppError(errtype.Execution, "Batch failure")
	if comp == nil || !comp.HasError() {
		t.Fatal("expected composite AppError")
	}
}

func TestMutexCollectionThreadSafety(t *testing.T) {
	mc := appfaults.NewMutexCollection()
	mc.AddType(errtype.NotFound).AddError(errtype.IO, errors.New("file missing"))
	if !mc.HasError() || mc.Count() != 2 {
		t.Fatalf("expected 2 items in MutexCollection, got %d", mc.Count())
	}
}

func TestContextBinding(t *testing.T) {
	ctx, coll := appfaults.WithFaults(context.Background())
	appfaults.RecordContextError(ctx, appfault.New(errtype.Precondition, "unmet state"))

	if !coll.HasError() || coll.Count() != 1 {
		t.Fatalf("expected 1 recorded error in context, got %d", coll.Count())
	}
}

func TestCollectionJSONMarshaling(t *testing.T) {
	c := appfaults.New().Add(appfault.New(errtype.Validation, "invalid token"))
	data, err := json.Marshal(c)
	if err != nil || len(data) == 0 {
		t.Fatalf("failed to marshal Collection: %v", err)
	}

	restored := appfaults.New()
	if err := json.Unmarshal(data, restored); err != nil || restored.Count() != 1 {
		t.Fatalf("failed to unmarshal Collection: %v", err)
	}
}
