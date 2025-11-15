class TrendsApiError(Exception):
    pass


class TrendsAuthError(TrendsApiError):
    pass


class TrendsRateLimitError(TrendsApiError):
    pass


class TrendsRetryableError(TrendsApiError):
    pass
