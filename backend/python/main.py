import asyncio
import os
from typing import Any, Dict, List

from dotenv import load_dotenv

from python.automation.growth_engine import GrowthEngine
from python.automation.trend_analyzer import TrendAnalyzer
from python.automation.content_engine import ContentEngine
from python.automation.experiment_engine import ExperimentEngine
from python.automation.scheduler import ContentScheduler
from python.automation.db_adapters import FileExperimentDB, SQLAlchemyCalendarDB


class DummyLLM:
    async def generate(self, prompt: str, **kwargs: Any) -> str:
        return f"{prompt}\n\n# Heading\nBody text with details."


class DummySocial:
    async def get_topic_metrics(self, topic: str) -> Dict[str, Any]:
        return {"engagement": 100, "recency_hours": 24, "related": [topic+" tips", "guide"]}
    # Optional publishing stub for scheduler
    async def publish(self, post: Dict[str, Any]) -> None:
        print(f"[PUBLISH] {post.get('id')} to {post.get('platform')} at {post.get('scheduled_at')}")


class DummySEO:
    async def get_keyword_data(self, topic: str) -> Dict[str, Any]:
        return {"search_volume": 500, "competition": 0.4, "related": [topic+" tutorial"]}


async def main() -> None:
    load_dotenv()
    db_url = os.getenv("DATABASE_URL")

    social = DummySocial()
    seo = DummySEO()

    trend = TrendAnalyzer(social_api=social, seo_client=seo)
    llm = DummyLLM()

    # DB adapters
    exp_db = FileExperimentDB()
    cal_db = SQLAlchemyCalendarDB(db_url)

    content_engine = ContentEngine(llm_client=llm, db=exp_db)
    # Inject trend analyzer for related topics path
    content_engine.trend_analyzer = trend  # type: ignore[attr-defined]
    experiment_engine = ExperimentEngine(db=exp_db, content_analyzer=None)

    growth = GrowthEngine(db=cal_db, analyzer=trend, content_engine=content_engine, experiment_engine=experiment_engine)
    scheduler = ContentScheduler(db=cal_db, publisher=social, analytics=trend)

    single_run = os.getenv("SINGLE_RUN") == "1"
    if single_run:
        await growth._run_cycle()
        await scheduler._process_scheduled_posts()
    else:
        await asyncio.gather(growth.start(), scheduler.run_schedule_cycle())


if __name__ == "__main__":
    asyncio.run(main())

