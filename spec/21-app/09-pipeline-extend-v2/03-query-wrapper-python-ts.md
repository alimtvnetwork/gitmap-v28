# Centralized Query Wrappers (Python/TS)

To prevent scattered `try/catch` logic and guarantee uniform error logging, this repository strictly mandates the use of generic query wrappers in both Python and TypeScript.

## 1. Python Wrapper Implementation

Python scripts (like those in `.github/scripts/`) MUST NOT use inline `try/except` blocks. They must import and execute `.github/scripts/query_wrapper.py`.

### Structure

```python
import traceback
from typing import Callable, Any, Dict

def query_wrapper(operation: Callable[..., Any], *args, **kwargs) -> Dict[str, Any]:
    try:
        data = operation(*args, **kwargs)
        return {"data": data, "error": None, "is_fail": False}
    except Exception as e:
        print(f"[QueryWrapper Error]: {str(e)}")
        return {"data": None, "error": e, "is_fail": True}
```

### Usage

```python
from query_wrapper import query_wrapper

def _read_file():
    with open("file.json", "r") as f:
        return json.load(f)

res = query_wrapper(_read_file)
if res["is_fail"]:
    return []
data = res["data"]
```

## 2. TypeScript Wrapper Implementation

TypeScript code MUST use `src/lib/queryWrapper.ts`.

### Strict Boolean State Checking

The wrapper returns a typed response with an `isFail` boolean property. AI Agents must strictly evaluate `isFail` rather than writing inverted logic like `!isSuccess` or `!response`.
