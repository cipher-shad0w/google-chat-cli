"""Notification support for the gchat TUI.

Sends terminal bell and/or desktop notifications when new unread
messages arrive.  The notification mode is read from the TUI config.
"""

import logging
import platform
import subprocess
import sys

logger = logging.getLogger(__name__)


def send_notification(title: str, body: str, mode: str) -> None:
    """Send a notification based on the configured mode.

    Args:
        title: The notification title (e.g., space name).
        body: The notification body (e.g., message preview).
        mode: One of "off", "bell", "desktop", "both".
    """
    if mode == "off":
        return

    if mode in ("bell", "both"):
        _send_bell()

    if mode in ("desktop", "both"):
        _send_desktop(title, body)


def _send_bell() -> None:
    """Send a terminal bell character."""
    try:
        sys.stdout.write("\a")
        sys.stdout.flush()
    except OSError:
        pass


def _send_desktop(title: str, body: str) -> None:
    """Send an OS-level desktop notification.

    Uses ``osascript`` on macOS and ``notify-send`` on Linux.
    Failures are logged but never raised.
    """
    system = platform.system()

    if system == "Darwin":
        _send_macos_notification(title, body)
    elif system == "Linux":
        _send_linux_notification(title, body)
    else:
        logger.debug("Desktop notifications not supported on %s", system)


def _send_macos_notification(title: str, body: str) -> None:
    """Send a notification via macOS ``osascript``."""
    # Escape double quotes in the strings for AppleScript
    safe_title = title.replace('"', '\\"')
    safe_body = body.replace('"', '\\"')
    script = (
        f'display notification "{safe_body}" with title "gchat" subtitle "{safe_title}"'
    )
    try:
        subprocess.run(
            ["osascript", "-e", script],
            capture_output=True,
            timeout=5,
        )
    except (subprocess.TimeoutExpired, OSError) as exc:
        logger.debug("Failed to send macOS notification: %s", exc)


def _send_linux_notification(title: str, body: str) -> None:
    """Send a notification via ``notify-send`` on Linux."""
    try:
        subprocess.run(
            ["notify-send", "--app-name=gchat", title, body],
            capture_output=True,
            timeout=5,
        )
    except FileNotFoundError:
        logger.debug(
            "notify-send not found — install libnotify for desktop notifications"
        )
    except (subprocess.TimeoutExpired, OSError) as exc:
        logger.debug("Failed to send Linux notification: %s", exc)
