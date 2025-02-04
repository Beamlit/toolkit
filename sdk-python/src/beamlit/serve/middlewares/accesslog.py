"""Module: accesslog

Defines the AccessLogMiddleware for logging HTTP requests and responses.
"""

from logging import getLogger

from starlette.middleware.base import BaseHTTPMiddleware


class AccessLogMiddleware(BaseHTTPMiddleware):
    """Middleware for logging access information of each HTTP request."""

    async def dispatch(self, request, call_next):
        """Processes each request to log access details.

        Args:
            request (Request): The incoming HTTP request.
            call_next (Callable): The next middleware or endpoint handler.

        Returns:
            Response: The HTTP response generated by the next handler.
        """
        logger = getLogger(__name__)
        response = await call_next(request)
        process_time = response.headers.get("X-Process-Time")
        rid_header = response.headers.get("X-Request-Id")
        request_id = rid_header or response.headers.get("X-Beamlit-Request-Id")
        logger.info(
            f"{request.method} {request.url.path} {response.status_code} {process_time}ms rid={request_id}"
        )
        return response
