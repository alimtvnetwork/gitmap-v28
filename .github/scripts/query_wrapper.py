import traceback
from typing import Callable, Any, Dict

def query_wrapper(operation: Callable[..., Any], *args, **kwargs) -> Dict[str, Any]:
    """
    Wraps an operation (API, database, or command logic), catching exceptions,
    logging errors explicitly to stderr without scattered try/catch, 
    and returning a structured result with explicit `is_success` and `is_fail` properties.
    """
    try:
        data = operation(*args, **kwargs)
        return {
            "data": data,
            "error": None,
            "is_success": True,
            "is_fail": False
        }
    except Exception as e:
        # We explicitly log the caught error as per the error management guidelines.
        print(f"[QueryWrapper Error]: {str(e)}")
        # If needed for deeper debugging, we could dump traceback
        # traceback.print_exc()
        
        return {
            "data": None,
            "error": e,
            "is_success": False,
            "is_fail": True
        }
