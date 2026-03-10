"""Simple REST API module for task management."""

from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional
from enum import Enum


class TaskStatus(Enum):
    """Status of a task."""
    PENDING = "pending"
    IN_PROGRESS = "in_progress"
    DONE = "done"


@dataclass
class Task:
    """Represents a task with title, status, and timestamps."""
    id: int
    title: str
    status: TaskStatus = TaskStatus.PENDING
    created_at: datetime = field(default_factory=datetime.now)
    completed_at: Optional[datetime] = None


class TaskRepository:
    """In-memory repository for managing tasks."""

    def __init__(self):
        self._tasks: dict[int, Task] = {}
        self._next_id: int = 1

    def create(self, title: str) -> Task:
        """Create a new task with the given title."""
        task = Task(id=self._next_id, title=title)
        self._tasks[self._next_id] = task
        self._next_id += 1
        return task

    def get(self, task_id: int) -> Optional[Task]:
        """Get a task by ID. Returns None if not found."""
        return self._tasks.get(task_id)

    def complete(self, task_id: int) -> bool:
        """Mark a task as done. Returns False if task not found."""
        task = self._tasks.get(task_id)
        if task is None:
            return False
        task.status = TaskStatus.DONE
        task.completed_at = datetime.now()
        return True

    def list_by_status(self, status: TaskStatus) -> list[Task]:
        """List all tasks with the given status."""
        return [t for t in self._tasks.values() if t.status == status]


def format_task(task: Task) -> str:
    """Format a task for display."""
    status_icon = {"pending": "⏳", "in_progress": "🔄", "done": "✅"}
    icon = status_icon.get(task.status.value, "❓")
    return f"{icon} [{task.id}] {task.title}"
