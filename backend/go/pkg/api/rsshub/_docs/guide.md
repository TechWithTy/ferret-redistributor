# RSSHub Authoring Notes

Collected notes from the official RSSHub docs so we can build/maintain custom routes consistently.

## Feed Fundamentals

- Channel fields: `title`, `link`, `description`, `language`, `image`, `icon`, `logo`, `subtitle`, `author`, iTunes metadata (`itunes_author`, `itunes_category`, `itunes_explicit`), `allowEmpty`. Compatibility varies by output format (Atom/JSON Feed/RSS 2.0).
- Item fields: `title`, `link`, `description`, `author`, `category`, `guid`, `pubDate`, `updated`, iTunes fields (`itunes_item_image`, `itunes_duration`), enclosures (`enclosure_url`, `enclosure_length`, `enclosure_type`), `media.*`, `doi`, and interaction counters (`upvotes`, `downvotes`, `comments`).
- Special feeds: magnet/podcast/media/Sci-Hub require extra fields (`enclosure_*`, `doi`, `media.*`) and documentation flags (`supportBT`, `supportScihub`, `supportPodcast`).
- Trim whitespace in `title`, `subtitle`, `author`, and use `<br>` for intentional line breaks in descriptions.

## Cache Usage

- RSSHub exposes `cache.tryGet(key, asyncFn, maxAge, refresh)`; results are cached per key with defaults governed by `CACHE_CONTENT_EXPIRE`.
- Only mutate objects inside the `tryGet` callback. Assignments outside are skipped on cache hits.
- Advanced helpers: `cache.get(key, refresh)` (returns raw string; use `JSON.parse`) and `cache.set(key, value, maxAge)`.

## Script Standard

- 4-space indent, semicolons, single quotes, template literals, `const/let`, `for...of`, arrow functions.
- Namespace layout under `lib/routes/<namespace>` must include `router.ts`, `maintainer.ts`, `radar.ts`, optional `templates/*.art`.
- Use SLD as namespace, keep entries sorted, prefer HTTPS/WebP.

## Date Handling

- If upstream has no date, leave `pubDate` undefined. Do not invent times.
- Prefer `parseDate(value, format?)` / `parseRelativeDate(value)` from `@/utils/parse-date`. Adjust to server time via `timezone(date, offsetHours)`.
- Pass actual `Date` objects whenever possible to avoid inconsistent parsing.

## Debugging Helpers

- `ctx.set('json', obj)` + `?format=debug.json` dumps custom objects for inspection (requires `debugInfo=true`).
- `?format={index}.debug.html` renders `data.item[index].description` directly for quick preview.

Use these references when adding new RSSHub integrations or documenting feature support.


