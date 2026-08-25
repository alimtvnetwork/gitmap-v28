package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/glyphs"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/theme"
	"github.com/spf13/cobra"
)

// ... I need to replace or append to root.go. I will use replace_file_content to inject dispatchSJ and its logic.
