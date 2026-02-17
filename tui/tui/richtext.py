"""Rich text formatting for Google Chat messages.

Converts Google Chat text formatting (bold, italic, strikethrough, code,
links) into Rich markup for display in the Textual TUI.
"""

import re

# Regex for URLs — match common URL patterns
_URL_RE = re.compile(r"(https?://[^\s<>\[\]()]+)")

# Code block: triple backticks (possibly with content spanning lines)
_CODE_BLOCK_RE = re.compile(r"```(.*?)```", re.DOTALL)

# Inline code: single backticks (not preceded/followed by backtick)
_INLINE_CODE_RE = re.compile(r"(?<!`)`(?!`)(.*?)(?<!`)`(?!`)")

# Bold: *text* (not preceded/followed by *)
_BOLD_RE = re.compile(r"(?<!\*)\*(?!\*)(.+?)(?<!\*)\*(?!\*)")

# Italic: _text_ (not preceded/followed by _)
_ITALIC_RE = re.compile(r"(?<!_)_(?!_)(.+?)(?<!_)_(?!_)")

# Strikethrough: ~text~ (not preceded/followed by ~)
_STRIKE_RE = re.compile(r"(?<!~)~(?!~)(.+?)(?<!~)~(?!~)")


def format_message_text(text: str) -> str:
    """Convert Google Chat formatted text to Rich markup.

    Processing order matters: code blocks first (to protect their content),
    then inline code, then URLs, then bold/italic/strikethrough.

    Args:
        text: Raw message text from Google Chat API.

    Returns:
        Rich-markup formatted text string.
    """
    if not text:
        return text

    # Phase 1: Extract and protect code blocks with placeholders
    code_blocks: list[str] = []

    def _replace_code_block(match: re.Match) -> str:
        idx = len(code_blocks)
        code_blocks.append(match.group(1).strip())
        return f"\x00CODEBLOCK{idx}\x00"

    result = _CODE_BLOCK_RE.sub(_replace_code_block, text)

    # Phase 2: Extract and protect inline code with placeholders
    inline_codes: list[str] = []

    def _replace_inline_code(match: re.Match) -> str:
        idx = len(inline_codes)
        inline_codes.append(match.group(1))
        return f"\x00INLINECODE{idx}\x00"

    result = _INLINE_CODE_RE.sub(_replace_inline_code, result)

    # Phase 3: Convert URLs to links
    def _replace_url(match: re.Match) -> str:
        url = match.group(1)
        # Escape any Rich markup chars in the URL display
        display = _escape_rich(url)
        return f"[underline #5b9bd5]{display}[/underline #5b9bd5]"

    result = _URL_RE.sub(_replace_url, result)

    # Phase 4: Apply text formatting
    result = _BOLD_RE.sub(r"[bold]\1[/bold]", result)
    result = _ITALIC_RE.sub(r"[italic]\1[/italic]", result)
    result = _STRIKE_RE.sub(r"[strike]\1[/strike]", result)

    # Phase 5: Restore inline code (with styling)
    for idx, code in enumerate(inline_codes):
        escaped = _escape_rich(code)
        result = result.replace(
            f"\x00INLINECODE{idx}\x00",
            f"[reverse] {escaped} [/reverse]",
        )

    # Phase 6: Restore code blocks (with styling)
    for idx, code in enumerate(code_blocks):
        escaped = _escape_rich(code)
        result = result.replace(
            f"\x00CODEBLOCK{idx}\x00",
            f"\n[on #2a2a2a]{escaped}[/on #2a2a2a]\n",
        )

    return result


def format_attachments(message: dict) -> str:
    """Build a display string for message attachments.

    Args:
        message: The full message dict from the Google Chat API.

    Returns:
        A Rich-markup string showing attachment indicators, or empty string.
    """
    attachments = message.get("attachment", [])
    if not attachments:
        return ""

    parts: list[str] = []
    for att in attachments:
        content_name = att.get("contentName", "")
        content_type = att.get("contentType", "")
        name = att.get("attachmentDataRef", {}).get("resourceName", "")

        # Determine icon based on content type
        if content_type.startswith("image/"):
            icon = "\U0001f5bc"
        elif content_type.startswith("video/"):
            icon = "\U0001f3ac"
        elif content_type.startswith("audio/"):
            icon = "\U0001f50a"
        elif content_type == "application/pdf":
            icon = "\U0001f4c4"
        else:
            icon = "\U0001f4ce"

        display_name = content_name or name or "attachment"
        parts.append(f"{icon} {display_name}")

    return "[dim]" + " | ".join(parts) + "[/dim]"


def _escape_rich(text: str) -> str:
    """Escape Rich markup special characters in text.

    Prevents user content from accidentally being interpreted as
    Rich tags.
    """
    # Rich uses square brackets for markup
    return text.replace("[", "\\[")
