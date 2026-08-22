import os

issues = []
with open('audit_results_extra.csv', 'r', encoding='utf-8') as f:
    for line in f:
        parts = line.strip().split('|')
        if len(parts) >= 4:
            issues.append({
                'file': parts[0],
                'line': int(parts[1]),
                'type': parts[2],
                'content': parts[3]
            })

from collections import defaultdict
grouped = defaultdict(list)
for item in issues:
    grouped[item['type']].append(item)

# Append to master plan
with open('.lovable/plans/pending/01-coding-guideline-fixes.md', 'a', encoding='utf-8') as f:
    subtask_idx = 10
    for cat_name, items in grouped.items():
        chunk_size = 150
        chunks = [items[i:i + chunk_size] for i in range(0, len(items), chunk_size)]
        
        for i, chunk in enumerate(chunks):
            slug = cat_name.replace(' ', '-').replace('.', '')
            subtask_filename = f"{subtask_idx:02d}-{slug}-part{i+1}.md"
            f.write(f"- [ ] `{subtask_filename}`: Fix {len(chunk)} `{cat_name}` issues.\n")
            
            subtask_path = f".lovable/plans/subtasks/01-coding-guideline-fixes/{subtask_filename}"
            with open(subtask_path, "w", encoding="utf-8") as sf:
                sf.write(f"# Fix {cat_name} (Part {i+1})\n\n")
                sf.write(f"Total items: {len(chunk)}\n\n")
                sf.write("## Files to Modify\n\n")
                for item in chunk:
                    sf.write(f"- `{item['file']}:{item['line']}`: `{item['content']}`\n")
            
            subtask_idx += 1

print("Updated plans.")
