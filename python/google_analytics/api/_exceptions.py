from __future__ import annotations

class ApiError(Exception):
    def __init__(self, status: int, message: str, reason: str | None = None) -> None:
        super().__init__(f"API {status}: {message}" + (f" (reason={reason})" if reason else ""))
        self.status = status
        self.message = message
        self.reason = reason

class AuthError(Exception):
    pass

class RateLimitError(ApiError):
    pass

class RetryableError(ApiError):
    pass
