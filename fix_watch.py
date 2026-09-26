import re

with open('cli/cmdagy/agy_watch_prompts.go', 'r') as f:
    content = f.read()

content = content.replace('''
func runAgyWatchPrompts(args []string) *apperror.AppError {
	for {
		// Implementation placeholder for loop and ssh logic.
		// Refresh every 30s as per requirements.
		time.Sleep(30 * time.Second)
	}
	return nil
}
''', '''
import "fmt"
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"

func runAgyWatchPrompts(args []string) *apperror.AppError {
	if isWatchPromptsJSON {
		fmt.Println("{ \"status\": \"watching\", \"interval\": 30 }")
		return nil
	}

	fmt.Println("Watching prompts (press Ctrl+C to stop)...")
	for {
		queues := cmdagy.DiscoverAllWorkspaceQueues()
		fmt.Printf("Discovered %d active queues\\n", len(queues))
		
		time.Sleep(30 * time.Second)
	}
	return nil
}
''')

with open('cli/cmdagy/agy_watch_prompts.go', 'w') as f:
    f.write(content)
