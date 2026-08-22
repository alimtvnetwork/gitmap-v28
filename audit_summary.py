import collections

counter = collections.Counter()
with open('audit_results_go.csv', 'r', encoding='utf-8') as f:
    for line in f:
        parts = line.strip().split('|')
        if len(parts) >= 3:
            counter[parts[2]] += 1

for k, v in counter.items():
    print(f"{k}: {v}")
