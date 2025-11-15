from __future__ import annotations

import enum
from datetime import datetime
from typing import Optional

from sqlalchemy import (
    Column,
    String,
    Text,
    DateTime,
    Boolean,
    Enum,
    ForeignKey,
    JSON,
    Integer,
    func,
)
from sqlalchemy.orm import declarative_base, relationship


Base = declarative_base()


class Platform(enum.Enum):
    instagram = "instagram"
    linkedin = "linkedin"
    twitter = "twitter"
    facebook = "facebook"
    youtube = "youtube"
    behiiv = "behiiv"


class CampaignStatus(enum.Enum):
    draft = "draft"
    active = "active"
    paused = "paused"
    completed = "completed"


class ScheduledStatus(enum.Enum):
    scheduled = "scheduled"
    published = "published"
    failed = "failed"
    canceled = "canceled"


class ContentItem(Base):
    __tablename__ = "content_items"
    id = Column(String(36), primary_key=True)
    title = Column(String(255), nullable=False)
    description = Column(Text, nullable=True)
    canonical_url = Column(String(1024), nullable=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now(), nullable=False)

    assets = relationship("Asset", back_populates="content", cascade="all, delete-orphan")
    hashtags = relationship("Hashtag", back_populates="content", cascade="all, delete-orphan")


class Campaign(Base):
    __tablename__ = "campaigns"
    id = Column(String(36), primary_key=True)
    name = Column(String(255), nullable=False)
    description = Column(Text, nullable=True)
    status = Column(Enum(CampaignStatus), nullable=False, default=CampaignStatus.draft)
    starts_at = Column(DateTime(timezone=True), nullable=True)
    ends_at = Column(DateTime(timezone=True), nullable=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now(), nullable=False)

    posts = relationship("ScheduledPost", back_populates="campaign", cascade="all, delete-orphan")


class ScheduledPost(Base):
    __tablename__ = "scheduled_posts"
    id = Column(String(36), primary_key=True)
    campaign_id = Column(String(36), ForeignKey("campaigns.id", ondelete="CASCADE"), nullable=False)
    content_id = Column(String(36), ForeignKey("content_items.id", ondelete="SET NULL"), nullable=True)
    platform = Column(Enum(Platform), nullable=False)
    caption = Column(Text, nullable=True)
    hashtags = Column(Text, nullable=True)
    scheduled_at = Column(DateTime(timezone=True), nullable=False)
    status = Column(Enum(ScheduledStatus), nullable=False, default=ScheduledStatus.scheduled)
    external_id = Column(String(255), nullable=True)  # Published ID/URN
    published_at = Column(DateTime(timezone=True), nullable=True)
    metadata = Column(JSON, nullable=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now(), nullable=False)

    campaign = relationship("Campaign", back_populates="posts")
    content = relationship("ContentItem")


class AssetType(enum.Enum):
    image = "image"
    video = "video"
    document = "document"


class Asset(Base):
    __tablename__ = "assets"
    id = Column(String(36), primary_key=True)
    content_id = Column(String(36), ForeignKey("content_items.id", ondelete="CASCADE"), nullable=False)
    kind = Column(Enum(AssetType), nullable=False)
    url = Column(String(2048), nullable=True)
    local_path = Column(String(2048), nullable=True)
    position = Column(Integer, nullable=True)  # ordering within carousel
    created_at = Column(DateTime(timezone=True), server_default=func.now(), nullable=False)

    content = relationship("ContentItem", back_populates="assets")


class Hashtag(Base):
    __tablename__ = "hashtags"
    id = Column(String(36), primary_key=True)
    content_id = Column(String(36), ForeignKey("content_items.id", ondelete="CASCADE"), nullable=False)
    tag = Column(String(100), nullable=False)
    created_at = Column(DateTime(timezone=True), server_default=func.now(), nullable=False)

    content = relationship("ContentItem", back_populates="hashtags")


class CalendarEvent(Base):
    __tablename__ = "calendar_events"
    id = Column(String(36), primary_key=True)
    title = Column(String(255), nullable=False)
    notes = Column(Text, nullable=True)
    start_at = Column(DateTime(timezone=True), nullable=False)
    end_at = Column(DateTime(timezone=True), nullable=True)
    all_day = Column(Boolean, default=False, nullable=False)
    created_at = Column(DateTime(timezone=True), server_default=func.now(), nullable=False)

