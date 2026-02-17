"""Configuration file support for the gchat TUI application.

Reads user settings from ``~/.config/gogchat/tui.toml`` and exposes them
via a :class:`TuiConfig` dataclass.  If the file is missing or contains
invalid values, sensible defaults are used.
"""

from __future__ import annotations

import logging
from dataclasses import dataclass
from pathlib import Path

try:
    import tomllib  # Python 3.11+
except ImportError:
    import tomli as tomllib  # type: ignore[no-redef]

logger = logging.getLogger(__name__)

_CONFIG_PATH = Path.home() / ".config" / "gogchat" / "tui.toml"

_VALID_NOTIFICATIONS = {"off", "bell", "desktop", "both"}


@dataclass
class TuiConfig:
    """All user-configurable settings for the TUI."""

    # Polling
    poll_interval: float = 15.0  # seconds between message polls (0 = disabled)

    # Messages
    message_page_size: int = 25  # number of messages to fetch

    # Cache
    members_cache_ttl: float = 3600.0  # seconds before members cache expires

    # Unread
    unread_check_workers: int = 8  # thread pool size for unread checks

    # Keybindings
    vim_mode: bool = False  # enable vim-style keybindings

    # Notifications
    notifications: str = "off"  # "off", "bell", "desktop", "both"

    # UI
    theme: str = "dark"  # for future use


def _load_config(path: Path = _CONFIG_PATH) -> TuiConfig:
    """Read and parse the TOML config file, returning a :class:`TuiConfig`.

    Missing files or individual invalid values are silently handled with
    defaults; warnings are logged for values that fail validation.
    """
    defaults = TuiConfig()

    if not path.is_file():
        logger.debug("Config file not found at %s, using defaults", path)
        return defaults

    try:
        raw = path.read_bytes()
        data = tomllib.loads(raw.decode("utf-8"))
    except (OSError, tomllib.TOMLDecodeError) as exc:
        logger.warning("Failed to read config file %s: %s — using defaults", path, exc)
        return defaults

    # -- polling ---------------------------------------------------------------
    poll_interval = defaults.poll_interval
    polling = data.get("polling", {})
    if isinstance(polling, dict) and "interval" in polling:
        val = polling["interval"]
        if isinstance(val, (int, float)) and val >= 0:
            poll_interval = float(val)
        else:
            logger.warning(
                "Invalid polling.interval=%r (must be >= 0), using default %.1f",
                val,
                defaults.poll_interval,
            )

    # -- messages --------------------------------------------------------------
    message_page_size = defaults.message_page_size
    messages = data.get("messages", {})
    if isinstance(messages, dict) and "page_size" in messages:
        val = messages["page_size"]
        if isinstance(val, int) and val > 0:
            message_page_size = val
        else:
            logger.warning(
                "Invalid messages.page_size=%r (must be int > 0), using default %d",
                val,
                defaults.message_page_size,
            )

    # -- cache -----------------------------------------------------------------
    members_cache_ttl = defaults.members_cache_ttl
    cache = data.get("cache", {})
    if isinstance(cache, dict) and "members_ttl" in cache:
        val = cache["members_ttl"]
        if isinstance(val, (int, float)) and val >= 0:
            members_cache_ttl = float(val)
        else:
            logger.warning(
                "Invalid cache.members_ttl=%r (must be >= 0), using default %.1f",
                val,
                defaults.members_cache_ttl,
            )

    # -- unread ----------------------------------------------------------------
    unread_check_workers = defaults.unread_check_workers
    unread = data.get("unread", {})
    if isinstance(unread, dict) and "check_workers" in unread:
        val = unread["check_workers"]
        if isinstance(val, int) and val > 0:
            unread_check_workers = val
        else:
            logger.warning(
                "Invalid unread.check_workers=%r (must be int > 0), using default %d",
                val,
                defaults.unread_check_workers,
            )

    # -- keybindings -----------------------------------------------------------
    vim_mode = defaults.vim_mode
    keybindings = data.get("keybindings", {})
    if isinstance(keybindings, dict) and "vim_mode" in keybindings:
        val = keybindings["vim_mode"]
        if isinstance(val, bool):
            vim_mode = val
        else:
            logger.warning(
                "Invalid keybindings.vim_mode=%r (must be bool), using default %s",
                val,
                defaults.vim_mode,
            )

    # -- notifications ---------------------------------------------------------
    notifications = defaults.notifications
    notif = data.get("notifications", {})
    if isinstance(notif, dict) and "mode" in notif:
        val = notif["mode"]
        if isinstance(val, str) and val in _VALID_NOTIFICATIONS:
            notifications = val
        else:
            logger.warning(
                "Invalid notifications.mode=%r (must be one of %s), using default %r",
                val,
                _VALID_NOTIFICATIONS,
                defaults.notifications,
            )

    # -- ui --------------------------------------------------------------------
    theme = defaults.theme
    ui = data.get("ui", {})
    if isinstance(ui, dict) and "theme" in ui:
        val = ui["theme"]
        if isinstance(val, str) and val:
            theme = val
        else:
            logger.warning(
                "Invalid ui.theme=%r (must be non-empty string), using default %r",
                val,
                defaults.theme,
            )

    return TuiConfig(
        poll_interval=poll_interval,
        message_page_size=message_page_size,
        members_cache_ttl=members_cache_ttl,
        unread_check_workers=unread_check_workers,
        vim_mode=vim_mode,
        notifications=notifications,
        theme=theme,
    )


# -- singleton ----------------------------------------------------------------

_instance: TuiConfig | None = None


def get_config() -> TuiConfig:
    """Return the lazily-created singleton :class:`TuiConfig` instance."""
    global _instance
    if _instance is None:
        _instance = _load_config()
    return _instance
