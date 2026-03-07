---
sidebar_position: 3
title: Python
---

# Python Support

Brf.it provides full support for Python with Tree-sitter based parsing.

## Supported Features

- Function definitions
- Class definitions
- Method definitions
- Decorators
- Type hints
- Docstrings
- Import statements (with `--include-imports`)

## Example

### Input

```python
from typing import List, Optional
from dataclasses import dataclass

@dataclass
class User:
    """Represents a system user."""
    id: int
    name: str
    email: str

class UserRepository:
    """Repository for user data access."""

    def __init__(self, db_connection):
        self.db = db_connection

    async def get_user(self, user_id: int) -> Optional[User]:
        """Fetch a user by ID."""
        # implementation...
        pass

    async def list_users(self, limit: int = 100) -> List[User]:
        """List all users with optional limit."""
        # implementation...
        pass
```

### Output (Brf.it)

```python
# user_repository.py
from typing import List, Optional
from dataclasses import dataclass

@dataclass
class User: id: int; name: str; email: str

class UserRepository:
    def __init__(self, db_connection)
    async def get_user(self, user_id: int) -> Optional[User]
    async def list_users(self, limit: int = 100) -> List[User]
```

## Extension

Files with `.py` extension are processed.
