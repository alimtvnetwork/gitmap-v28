import re

def replace_in_file(path, old, new):
    with open(path, 'r', encoding='utf-8') as f:
        content = f.read()
    if old in content:
        content = content.replace(old, new)
        with open(path, 'w', encoding='utf-8') as f:
            f.write(content)

replace_in_file('cli/cmdagy/agy_cmd.go', 'runAgyLs()', 'runAgyLs(5)')
replace_in_file('cli/cmdagy/agy_conv_ls.go', 'queues := DiscoverAllWorkspaceQueues()', 'queues, _ := DiscoverAllWorkspaceQueues()')
replace_in_file('cli/cmdagy/agy_inject_prompts.go', 'apperror.NewSimple("ip requires an argument")', 'apperror.NewSimple("ip requires an argument", "EIP001")')
replace_in_file('cli/cmdagy/agy_look_prompts.go', 'queues := DiscoverAllWorkspaceQueues()', 'queues, _ := DiscoverAllWorkspaceQueues()')
replace_in_file('cli/cmdagy/agy_misc.go', 'runAgyLs()', 'runAgyLs(5)')
replace_in_file('cli/cmdagy/agy_pins_runners.go', 'processAgyLsProjects(workspaceURI)', 'processAgyLsProjects(workspaceURI, 5)')
replace_in_file('cli/cmdagy/agy_watch_prompts.go', 'queues := DiscoverAllWorkspaceQueues()', 'queues, _ := DiscoverAllWorkspaceQueues()')

with open('cli/cmdagy/agy_most_conv.go', 'r', encoding='utf-8') as f:
    content = f.read()
content = re.sub(r'"database/sql"\n', '', content)
content = re.sub(r'"sort"\n', '', content)
with open('cli/cmdagy/agy_most_conv.go', 'w', encoding='utf-8') as f:
    f.write(content)

print("Fixed compilation errors.")
