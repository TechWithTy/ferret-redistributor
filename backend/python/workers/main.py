import logging
import os
import time
from typing import Callable

logging.basicConfig(
    level=os.getenv("SOCIAL_SCALE_LOG_LEVEL", "INFO"),
    format="%(asctime)s %(levelname)s %(message)s",
)
LOGGER = logging.getLogger("social-scale.worker")


def poll_work(fetch_job: Callable[[], dict | None]) -> None:
    """Simple polling loop placeholder for future queue integration."""
    poll_interval = float(os.getenv("SOCIAL_SCALE_POLL_INTERVAL", "5"))

    while True:
        job = fetch_job()
        if job:
            LOGGER.info("Processing job: %s", job["id"])
            time.sleep(1)
            LOGGER.info("Job %s complete", job["id"])
        else:
            LOGGER.debug("No jobs available; sleeping %.1fs", poll_interval)
            time.sleep(poll_interval)


def _fake_fetch() -> dict | None:
    return None


if __name__ == "__main__":
    LOGGER.info("Social Scale worker booted")
    poll_work(_fake_fetch)


