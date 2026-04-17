"""
Logging configuration.
Ticket PM: Logging and Monitoring — request/response logging, routing info, error log.
"""

from __future__ import annotations

import logging
import os
import sys
from logging.handlers import RotatingFileHandler
from pathlib import Path


def setup_logging(
    log_dir: str = "outputs/logs",
    log_level: str = "INFO",
) -> None:
    """
    Setup structured logging:
    - Console: INFO level, human-readable
    - File (rotating): DEBUG level, untuk observability
    - Error file: ERROR level, untuk alerting
    
    Ticket PM: Error log disediakan untuk observability.
    """
    log_level_int = getattr(logging, log_level.upper(), logging.INFO)

    # Format: timestamp | level | logger | message
    fmt = "%(asctime)s | %(levelname)-8s | %(name)s | %(message)s"
    datefmt = "%Y-%m-%d %H:%M:%S"
    formatter = logging.Formatter(fmt, datefmt=datefmt)

    root_logger = logging.getLogger()
    root_logger.setLevel(logging.DEBUG)

    # Console handler
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setLevel(log_level_int)
    console_handler.setFormatter(formatter)
    root_logger.addHandler(console_handler)

    # File handler (rotating, max 5MB x 3 files)
    Path(log_dir).mkdir(parents=True, exist_ok=True)
    file_handler = RotatingFileHandler(
        filename=os.path.join(log_dir, "app.log"),
        maxBytes=5 * 1024 * 1024,
        backupCount=3,
        encoding="utf-8",
    )
    file_handler.setLevel(logging.DEBUG)
    file_handler.setFormatter(formatter)
    root_logger.addHandler(file_handler)

    # Error-only file handler
    error_handler = RotatingFileHandler(
        filename=os.path.join(log_dir, "error.log"),
        maxBytes=2 * 1024 * 1024,
        backupCount=2,
        encoding="utf-8",
    )
    error_handler.setLevel(logging.ERROR)
    error_handler.setFormatter(formatter)
    root_logger.addHandler(error_handler)

    # Suppress noisy third-party loggers
    logging.getLogger("urllib3").setLevel(logging.WARNING)
    logging.getLogger("requests").setLevel(logging.WARNING)