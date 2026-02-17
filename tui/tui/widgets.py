"""Custom widgets for the gchat TUI."""

import logging

from textual import events, work
from textual.app import ComposeResult
from textual.await_remove import AwaitRemove
from textual.containers import Horizontal, Vertical
from textual.message import Message
from textual.screen import ModalScreen
from textual.widgets import Input, Label, ListItem, ListView, Static, TextArea

logger = logging.getLogger(__name__)


class VimListView(ListView):
    """A ListView with vim-style keybindings.

    j/k navigate up/down and g/G jump to first/last item.
    """

    def _on_key(self, event: events.Key) -> None:
        if event.key == "j":
            event.prevent_default()
            event.stop()
            self.action_cursor_down()
            return
        elif event.key == "k":
            event.prevent_default()
            event.stop()
            self.action_cursor_up()
            return
        elif event.key == "g":
            event.prevent_default()
            event.stop()
            self.index = 0
            self.scroll_visible()
            return
        elif event.key == "G":
            event.prevent_default()
            event.stop()
            if len(self) > 0:
                self.index = len(self) - 1
                self.scroll_visible()
            return
        super()._on_key(event)


class CategoryHeader(ListItem):
    """A non-selectable section header in the spaces list."""

    def __init__(self, title: str) -> None:
        self.category_title = title
        super().__init__()

    def compose(self) -> ComposeResult:
        yield Label(
            f"[bold dim]{self.category_title}[/bold dim]", classes="category-header"
        )


class SpaceItem(ListItem):
    """A list item representing a chat space."""

    def __init__(
        self,
        space_name: str,
        display_name: str,
        has_unread: bool = False,
        space_type: str = "",
    ) -> None:
        """Initialize with space metadata.

        Args:
            space_name: The space resource name (e.g., "spaces/AAAA")
            display_name: Human-readable name to show in the list
            has_unread: Whether this space has unread messages
            space_type: The space type (e.g., "SPACE", "GROUP_CHAT", "DIRECT_MESSAGE")
        """
        label_text = display_name if display_name else space_name.split("/")[-1]
        self.space_name = space_name
        self.display_name = display_name
        self.has_unread = has_unread
        self.space_type = space_type
        super().__init__()
        self._label_text = label_text

    def compose(self) -> ComposeResult:
        """Create the space item layout with optional unread dot."""
        with Horizontal(classes="space-row"):
            yield Label("● " if self.has_unread else "  ", classes="unread-dot")
            yield Label(self._label_text, classes="space-name")


class GroupsPanel(Static):
    """Left panel displaying the list of chat groups."""

    class SpaceSelected(Message):
        """Posted when a space is selected."""

        def __init__(self, space_name: str, display_name: str) -> None:
            self.space_name = space_name
            self.display_name = display_name
            super().__init__()

    def compose(self) -> ComposeResult:
        """Create the groups list view with search filter."""
        yield Input(placeholder="Filter spaces…", id="space-filter")
        yield VimListView(id="groups-list")

    def on_mount(self) -> None:
        """Set the border title and load spaces."""
        self.border_title = "Groups"
        self._current_spaces: list[dict] = []
        self._all_spaces: list[dict] = []
        self._filter_text: str = ""
        self._category_filter: str = "all"  # "all", "spaces", "dms"
        self.load_spaces()
        # Ensure the list is focused, not the filter input
        self.set_timer(0.1, lambda: self.query_one("#groups-list").focus())

    def on_input_changed(self, event: Input.Changed) -> None:
        """Filter spaces list when the search input changes."""
        if event.input.id != "space-filter":
            return
        self._filter_text = event.value.lower().strip()
        self._apply_space_filter()

    def _apply_space_filter(self) -> None:
        """Re-populate the spaces list based on the current filter text."""
        if not self._filter_text:
            # Show all spaces
            self._populate_spaces(self._all_spaces)
        else:
            # Filter spaces by display name
            filtered = [
                s
                for s in self._all_spaces
                if self._filter_text
                in (s.get("displayName", "") or s.get("name", "")).lower()
            ]
            self._populate_spaces(filtered)

    def load_spaces(self) -> None:
        """Load spaces — show cached data instantly, then refresh from API."""
        # Phase 0: Show cached data immediately (runs on the main thread)
        from tui.cache import get_cache

        cached_spaces = get_cache().get_spaces()
        if cached_spaces:
            # Show cached spaces immediately
            cached_unread = get_cache().get_unread_states()
            unread_set: set[str] = set()
            if cached_unread:
                # Spaces with entries are the ones that had unread messages last time
                unread_set = set(cached_unread.keys())
            self._populate_spaces(cached_spaces, unread_set)

        self._refresh_spaces()

    def _spaces_unchanged(self, new_spaces: list[dict]) -> bool:
        """Return True if new_spaces matches the currently displayed spaces."""
        old = self._current_spaces
        if len(old) != len(new_spaces):
            return False
        # Compare space resource names in order
        return all(o.get("name") == n.get("name") for o, n in zip(old, new_spaces))

    @work(thread=True)
    def _refresh_spaces(self) -> None:
        """Fetch fresh spaces from API and update the UI only if changed."""
        from tui.cli import list_spaces_fresh

        spaces = list_spaces_fresh()

        if not spaces and not self._current_spaces:
            # No spaces returned and nothing cached — might be an error
            self.app.call_from_thread(
                self.app.notify,
                "No spaces found — check connection or auth",
                severity="warning",
                timeout=5,
            )

        # Only re-render if the spaces list has actually changed
        if not self._spaces_unchanged(spaces):
            self.app.call_from_thread(self._populate_spaces, spaces)

        # Always kick off the unread check (unread states may have changed)
        self.app.call_from_thread(self._start_unread_check, spaces)

    def _start_unread_check(self, spaces: list[dict]) -> None:
        """Start the background unread state check."""
        self._load_unread_states(spaces)

    @work(thread=True)
    def _load_unread_states(self, spaces: list[dict]) -> None:
        """Check unread status for all spaces in parallel, then update UI."""
        from concurrent.futures import ThreadPoolExecutor, as_completed
        from tui.cli import get_space_read_state, list_messages

        unread_spaces: set[str] = set()

        def check_space(space_name: str) -> tuple[str, bool]:
            last_read_time = get_space_read_state(space_name)
            if last_read_time is None:
                return (space_name, False)
            messages = list_messages(space_name, limit=1)
            if messages:
                latest_msg_time = messages[-1].get("createTime", "")
                if latest_msg_time and latest_msg_time > last_read_time:
                    return (space_name, True)
            return (space_name, False)

        space_names = [s.get("name", "") for s in spaces if s.get("name")]

        from tui.config import get_config

        with ThreadPoolExecutor(
            max_workers=get_config().unread_check_workers
        ) as executor:
            futures = {executor.submit(check_space, name): name for name in space_names}
            for future in as_completed(futures):
                try:
                    name, is_unread = future.result()
                    if is_unread:
                        unread_spaces.add(name)
                except Exception as exc:
                    logger.warning(
                        "Failed to check unread state for %s: %s",
                        futures[future],
                        exc,
                    )

        # Cache the unread states for instant display next time
        from tui.cache import get_cache

        get_cache().set_unread_states({name: "unread" for name in unread_spaces})

        if unread_spaces:
            self.app.call_from_thread(self._update_unread_indicators, unread_spaces)

    def _update_unread_indicators(self, unread_spaces: set[str]) -> None:
        """Update unread dot indicators for spaces that have unread messages."""
        groups_list = self.query_one("#groups-list", ListView)
        for item in groups_list.children:
            if isinstance(item, SpaceItem) and item.space_name in unread_spaces:
                item.has_unread = True
                try:
                    dot_label = item.query_one(".unread-dot", Label)
                    dot_label.update("● ")
                except Exception:
                    pass

    def _populate_spaces(
        self, spaces: list[dict], unread_spaces: set[str] | None = None
    ) -> None:
        """Populate the ListView with space items grouped by type."""
        # Store the full list (only when called without filtering)
        if not self._filter_text:
            self._all_spaces = spaces

        groups_list = self.query_one("#groups-list", ListView)
        groups_list.clear()

        if unread_spaces is None:
            unread_spaces = set()

        if not spaces:
            groups_list.append(ListItem(Label("No spaces found")))
            return

        # Classify spaces
        dm_spaces = []
        group_spaces = []
        named_spaces = []

        for space in spaces:
            space_type = space.get("spaceType", "")
            if space_type == "DIRECT_MESSAGE":
                dm_spaces.append(space)
            elif space_type == "SPACE":
                named_spaces.append(space)
            else:
                # GROUP_CHAT or unknown
                group_spaces.append(space)

        # Apply category filter
        show_spaces = self._category_filter in ("all", "spaces")
        show_dms = self._category_filter in ("all", "dms")

        if show_spaces and (named_spaces or group_spaces):
            if self._category_filter == "all":
                groups_list.append(CategoryHeader("Spaces"))
            for space in named_spaces + group_spaces:
                space_name = space.get("name", "")
                display_name = space.get("displayName", "")
                has_unread = space_name in unread_spaces
                groups_list.append(
                    SpaceItem(
                        space_name,
                        display_name,
                        has_unread=has_unread,
                        space_type=space.get("spaceType", ""),
                    )
                )

        if show_dms and dm_spaces:
            if self._category_filter == "all":
                groups_list.append(CategoryHeader("Direct Messages"))
            for space in dm_spaces:
                space_name = space.get("name", "")
                display_name = space.get("displayName", "")
                has_unread = space_name in unread_spaces
                groups_list.append(
                    SpaceItem(
                        space_name,
                        display_name,
                        has_unread=has_unread,
                        space_type=space.get("spaceType", ""),
                    )
                )

        self._current_spaces = spaces

    def on_list_view_selected(self, event: ListView.Selected) -> None:
        """Handle space selection."""
        if isinstance(event.item, SpaceItem):
            self.post_message(
                self.SpaceSelected(event.item.space_name, event.item.display_name)
            )

    def mark_space_as_read(self, space_name: str) -> None:
        """Remove the unread indicator from a specific space item."""
        groups_list = self.query_one("#groups-list", ListView)
        for item in groups_list.children:
            if isinstance(item, SpaceItem) and item.space_name == space_name:
                item.has_unread = False
                try:
                    dot_label = item.query_one(".unread-dot", Label)
                    dot_label.update("  ")
                except Exception:
                    pass
                break


class MessageItem(ListItem):
    """A list item representing a single chat message."""

    def __init__(
        self,
        content: str,
        sender_user_id: str | None = None,
        is_name_resolved: bool = True,
        prefix_markup: str = "",
        prefix_width: int = 0,
        body_text: str = "",
        message_name: str = "",
        reactions_markup: str = "",
    ) -> None:
        """Initialize with message content.

        Args:
            content: Full Rich-markup formatted message text (kept for raw storage).
            sender_user_id: The sender's resource name (e.g. "users/123").
            is_name_resolved: False when the displayed name is a raw user ID.
            prefix_markup: Rich markup for the prefix (timestamp + name + padding).
            prefix_width: The visual character width of the prefix.
            body_text: The plain message body text (will wrap independently).
            message_name: The message resource name (e.g. "spaces/AAAA/messages/BBBB").
            reactions_markup: Rich markup string for reactions display.
        """
        # If prefix_markup and body_text are provided, use two-column layout
        if prefix_markup and prefix_width > 0:
            prefix_label = Static(prefix_markup, markup=True, classes="msg-prefix")
            prefix_label.styles.width = prefix_width
            prefix_label.styles.min_width = prefix_width
            prefix_label.styles.max_width = prefix_width
            body_label = Label(body_text, markup=True, classes="msg-body")
            if reactions_markup:
                body_content = Vertical(
                    body_label,
                    Label(reactions_markup, markup=True, classes="msg-reactions"),
                    classes="msg-body-container",
                )
                container = Horizontal(prefix_label, body_content, classes="msg-row")
            else:
                container = Horizontal(prefix_label, body_label, classes="msg-row")
            super().__init__(container)
        else:
            # Fallback for system/status messages
            super().__init__(Label(content, markup=True))
        self.message_content = content
        self.sender_user_id = sender_user_id
        self.is_name_resolved = is_name_resolved
        self.prefix_width = prefix_width
        self.body_text = body_text
        self.message_name = message_name


class ChatLog(VimListView):
    """A ListView that displays chat messages as selectable items."""

    def __init__(self, **kwargs) -> None:
        # Remove any kwargs that are specific to RichLog and not applicable to ListView
        kwargs.pop("highlight", None)
        kwargs.pop("markup", None)
        kwargs.pop("wrap", None)
        kwargs.pop("min_width", None)
        super().__init__(**kwargs)
        self._raw_entries: list[str] = []

    def write_message(self, content: str) -> None:
        """Write a message as a new selectable list item."""
        self._raw_entries.append(content)
        self.append(MessageItem(content))
        self.index = len(self._raw_entries) - 1

    def clear(self) -> AwaitRemove:
        """Clear all messages."""
        self._raw_entries.clear()
        return super().clear()


class ChatPanel(Static):
    """Main panel displaying chat messages with optional search."""

    def compose(self) -> ComposeResult:
        """Create the chat log view with search bar."""
        yield Input(
            placeholder="Search messages…",
            id="message-search",
            classes="search-hidden",
        )
        yield ChatLog(id="chat-log")

    def on_mount(self) -> None:
        """Set the border title when mounted."""
        self.border_title = "Chat"
        self._search_visible = False

    def toggle_search(self) -> None:
        """Toggle the message search bar visibility."""
        search_input = self.query_one("#message-search", Input)
        self._search_visible = not self._search_visible
        if self._search_visible:
            search_input.remove_class("search-hidden")
            search_input.add_class("search-visible")
            search_input.focus()
        else:
            search_input.add_class("search-hidden")
            search_input.remove_class("search-visible")
            search_input.value = ""
            # Restore all messages visibility
            self._clear_message_filter()

    def on_input_changed(self, event: Input.Changed) -> None:
        """Filter messages when the search input changes."""
        if event.input.id != "message-search":
            return
        search_text = event.value.lower().strip()
        self._filter_messages(search_text)

    def _filter_messages(self, search_text: str) -> None:
        """Show/hide messages based on search text."""
        chat_log = self.query_one("#chat-log", ChatLog)
        for item in chat_log.children:
            if isinstance(item, MessageItem):
                if not search_text:
                    item.display = True
                else:
                    # Search in the raw message content
                    content = item.message_content.lower()
                    body = (item.body_text or "").lower()
                    if search_text in content or search_text in body:
                        item.display = True
                    else:
                        item.display = False

    def _clear_message_filter(self) -> None:
        """Show all messages (clear the filter)."""
        chat_log = self.query_one("#chat-log", ChatLog)
        for item in chat_log.children:
            if isinstance(item, MessageItem):
                item.display = True


class MessageInput(TextArea):
    """A TextArea that sends on Enter and inserts newline on Shift+Enter."""

    class Submitted(Message):
        """Posted when the user presses Enter to send a message."""

        def __init__(self, value: str) -> None:
            self.value = value
            super().__init__()

    def _on_key(self, event: events.Key) -> None:
        if event.key == "enter":
            event.prevent_default()
            event.stop()
            text = self.text.strip()
            if text:
                self.post_message(self.Submitted(text))
                self.clear()
            return
        super()._on_key(event)


class InputPanel(Static):
    """Bottom panel for composing messages."""

    def compose(self) -> ComposeResult:
        """Create the message input area."""
        yield MessageInput(id="message-input")

    def on_mount(self) -> None:
        """Set the border title when mounted."""
        self.border_title = "Message"


class NameInputScreen(ModalScreen[str | None]):
    """Modal dialog to assign a display name to an unknown user."""

    BINDINGS = [("escape", "cancel", "Cancel")]

    def __init__(self, user_id: str) -> None:
        self.user_id = user_id
        super().__init__()

    def compose(self) -> ComposeResult:
        with Vertical(id="name-dialog"):
            yield Label(
                f"Enter name for [bold]{self.user_id}[/bold]:",
                id="name-dialog-title",
            )
            yield Input(placeholder="Display name…", id="name-input")

    def on_mount(self) -> None:
        self.query_one("#name-input", Input).focus()

    def on_input_submitted(self, event: Input.Submitted) -> None:
        value = event.value.strip()
        if value:
            self.dismiss(value)
        else:
            self.dismiss(None)

    def action_cancel(self) -> None:
        self.dismiss(None)


class ReactionScreen(ModalScreen[str | None]):
    """Modal dialog to select an emoji reaction for a message."""

    BINDINGS = [("escape", "cancel", "Cancel")]

    REACTIONS = ["👍", "👎", "❤️", "😂", "😮", "😢", "🎉", "🤔"]

    def compose(self) -> ComposeResult:
        with Vertical(id="reaction-dialog"):
            yield Label(
                "Add reaction  " + "  ".join(self.REACTIONS),
                id="reaction-dialog-title",
            )
            yield Input(placeholder="Type emoji or paste…", id="reaction-input")

    def on_mount(self) -> None:
        self.query_one("#reaction-input", Input).focus()

    def on_input_submitted(self, event: Input.Submitted) -> None:
        value = event.value.strip()
        if value:
            self.dismiss(value)
        else:
            self.dismiss(None)

    def action_cancel(self) -> None:
        self.dismiss(None)


class MessageActionScreen(ModalScreen[str | None]):
    """Modal dialog showing available actions for a message."""

    BINDINGS = [("escape", "cancel", "Cancel")]

    def __init__(
        self, message_name: str, body_text: str, is_own_message: bool = False
    ) -> None:
        self.message_name = message_name
        self.body_text = body_text
        self.is_own_message = is_own_message
        super().__init__()

    def compose(self) -> ComposeResult:
        with Vertical(id="action-dialog"):
            yield Label("Message Actions", id="action-dialog-title")
            items = []
            if self.is_own_message:
                items.append(ListItem(Label("📝  Edit message"), id="action-edit"))
                items.append(ListItem(Label("🗑️  Delete message"), id="action-delete"))
            items.append(ListItem(Label("💬  Quote reply"), id="action-quote"))
            yield ListView(*items, id="action-list")

    def on_mount(self) -> None:
        self.query_one("#action-list", ListView).focus()

    def on_list_view_selected(self, event: ListView.Selected) -> None:
        action_id = event.item.id
        if action_id == "action-edit":
            self.dismiss("edit")
        elif action_id == "action-delete":
            self.dismiss("delete")
        elif action_id == "action-quote":
            self.dismiss("quote")

    def action_cancel(self) -> None:
        self.dismiss(None)


class EditMessageScreen(ModalScreen[str | None]):
    """Modal dialog to edit a message's text."""

    BINDINGS = [
        ("escape", "cancel", "Cancel"),
        ("ctrl+s", "save", "Save"),
    ]

    def __init__(self, current_text: str) -> None:
        self.current_text = current_text
        super().__init__()

    def compose(self) -> ComposeResult:
        with Vertical(id="edit-dialog"):
            yield Label(
                "Edit message [dim](Ctrl+S to save, Escape to cancel)[/dim]:",
                id="edit-dialog-title",
            )
            yield TextArea(self.current_text, id="edit-input")

    def on_mount(self) -> None:
        self.query_one("#edit-input", TextArea).focus()

    def action_save(self) -> None:
        """Save the edited message."""
        text = self.query_one("#edit-input", TextArea).text.strip()
        if text:
            self.dismiss(text)

    def action_cancel(self) -> None:
        self.dismiss(None)


class ConfirmDeleteScreen(ModalScreen[bool]):
    """Modal confirmation dialog for deleting a message."""

    BINDINGS = [
        ("escape", "cancel", "Cancel"),
        ("y", "confirm", "Yes"),
        ("n", "cancel", "No"),
    ]

    def compose(self) -> ComposeResult:
        with Vertical(id="confirm-dialog"):
            yield Label(
                "Delete this message? [bold](y/n)[/bold]",
                id="confirm-dialog-title",
            )

    def action_confirm(self) -> None:
        self.dismiss(True)

    def action_cancel(self) -> None:
        self.dismiss(False)
