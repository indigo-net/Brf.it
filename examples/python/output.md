# Code Summary: examples/python

*brf.it dev*

---

## Files

### examples/python/api.py

```python
class TaskStatus(Enum)
class Task
class TaskRepository
def __init__(self)
def create(self, title: str) -> Task
def get(self, task_id: int) -> Optional[Task]
def complete(self, task_id: int) -> bool
def list_by_status(self, status: TaskStatus) -> list[Task]
def format_task(task: Task) -> str
```

---

