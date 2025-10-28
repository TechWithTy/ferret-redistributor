from __future__ import annotations

import os
from contextlib import contextmanager
from typing import Iterator

from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, Session

from .models import Base


def get_engine(url: str | None = None):
    url = url or os.getenv("DATABASE_URL", "postgresql+psycopg2://user:pass@localhost:5432/ferret")
    return create_engine(url, pool_pre_ping=True)


def create_all(url: str | None = None) -> None:
    engine = get_engine(url)
    Base.metadata.create_all(engine)


def get_session(url: str | None = None) -> Session:
    engine = get_engine(url)
    return sessionmaker(bind=engine, autoflush=False, autocommit=False)()


@contextmanager
def session_scope(url: str | None = None) -> Iterator[Session]:
    session = get_session(url)
    try:
        yield session
        session.commit()
    except Exception:
        session.rollback()
        raise
    finally:
        session.close()

