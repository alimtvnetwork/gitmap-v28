package cmdpipeline

import (
	"fmt"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestQueryRepoSuggestions(t *testing.T) {
	conn, err := openSuggestDB()
	fmt.Printf("DEBUG: openSuggestDB err=%v, path=%s\n", err, store.DefaultDBPath())
	if conn != nil {
		rows, err := conn.Query("SELECT COUNT(*) FROM Repo")
		fmt.Printf("DEBUG: Repo count query err=%v\n", err)
		if rows != nil {
			var count int
			if rows.Next() {
				_ = rows.Scan(&count)
			}
			rows.Close()
			fmt.Printf("DEBUG: Repo count = %d\n", count)
		}
	}
	suggs := QueryRepoSuggestionsFromDB("movi-cli")
	fmt.Printf("DEBUG: suggs for movi-cli = %v\n", suggs)
}
