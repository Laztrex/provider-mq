from uuid import UUID, uuid4

from pydantic import Field, BaseModel


class Job(BaseModel):
    """
    Класс задачи
    uid: идентификатор задачи
    status: этап выполнения задачи
    result: результат выполенения задачи
    """

    uid: UUID = Field(default_factory=uuid4)
    status: str = "in_progress"
    result: int = None

    no_return: bool = False

    def asdict(self):
        return {
            "uid": str(self.uid),
            "status": self.status,
            "result": self.result
        }
