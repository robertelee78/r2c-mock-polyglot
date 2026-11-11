import os

from pydantic import BaseModel
from sqlalchemy import Boolean, String, create_engine
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


class Base(DeclarativeBase):
    pass


class Account(Base):
    __tablename__ = "crossfile_accounts"

    id: Mapped[int] = mapped_column(primary_key=True)
    display: Mapped[str] = mapped_column(String(120))
    is_admin: Mapped[bool] = mapped_column(Boolean, default=False)


class AccountPayload(BaseModel):
    display: str
    is_admin: bool = False


engine = create_engine(os.environ["O84_DATABASE_URL"])
