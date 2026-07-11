# v2 Ideas (parked, not v1)

Do not build these until v1 is shipped and used for a few weeks of real sessions.

- Session-completion prediction (classifier on time-of-day, interval length, genre/playlist → complete vs. abandon). Needs real accumulated session data first — not viable on day one.
- Focus-pattern clustering / descriptive stats ("your most consistent block is 9-11am", completion rate by playlist).
- Dark mode.
- Now-playing album art driving accent color.
- Multi-device sync polish beyond default Postgres behavior.
- Ambient sound layering / non-Spotify audio sources.
- Native Mac App Store app (separate goal: learning Swift/SwiftUI, not a v1 priority).

## Bigger-picture vision: AIO focus workspace

The longer-term idea is for Loft to become an all-in-one workspace that jumbles Spotify, Google Calendar, timers, notes/todos, etc. together — so getting "in the zone" doesn't mean tabbing between four separate apps. Each piece below is a second (or third) OAuth provider on top of Spotify, so this is real added architectural complexity, not a small bolt-on — evaluate one at a time, after v1 is shipped.

- **Google Calendar integration** — full read/write access (not just read-only display). Show upcoming events during a session (so you know how much real focus time you have), and write completed sessions back to the calendar as events, so Loft session history becomes visible in the calendar you already look at. Longer term: a general "what's coming up" view so you don't have to hold your whole schedule in your head.
- **Notes/todo integration** — simple note-taker or todo list, service not yet decided (Notion is one candidate given its API; a local-file/Obsidian-style approach is a different technical shape — no hosted API, more like reading/writing markdown files — and would need to be designed separately if chosen). Use case: pull tasks in, and/or push session summaries out as a lightweight work log.
- **LLM assistant with app-wide access** — an LLM that can navigate the app and execute commands on the user's behalf (e.g. "start a 25-minute session on the essay draft," "what did I work on yesterday"). Flagged by Benson himself as the hardest one to generalize: for a multi-user app, someone has to provide the LLM access/tokens and eat the cost — needs a real decision (bring-your-own-key? Loft-hosted with usage limits?) before this is buildable, not just a wiring problem like the others.
