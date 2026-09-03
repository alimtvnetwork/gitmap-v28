package installer

import (
	"context"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestExecuteOrder(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script := &model.InstallerScript{
		Name:      "Ordered App",
		Slug:      "ordered-app",
		TargetOS:  "ubuntu",
		OrderMode: constants.OrderUnixFirst,
		Version:   "v1.0.0",
	}

	db.CreateInstaller(script)

	ctx := context.Background()
	if err := mgr.ExecuteOrdered(ctx, "ordered-app", "ubuntu"); err != nil {
		t.Fatalf("ExecuteOrdered failed: %v", err)
	}

	if errEmpty := mgr.ExecuteOrdered(ctx, "", "ubuntu"); errEmpty == nil {
		t.Fatal("expected error on empty slug")
	}
}
