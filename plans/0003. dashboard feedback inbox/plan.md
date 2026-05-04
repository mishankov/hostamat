# Plan: Dashboard Feedback Inbox

> Source PRD: `plans/0003. dashboard feedback inbox/prd.md`

## Architectural Decisions

Durable decisions that apply across all phases:

- **User-facing route surface**: authenticated feedback submission is available from the dashboard shell and is proxied from the webapp to the backend as a `POST` operation.
- **Backend API route**: expose an authenticated `POST /api/feedback` endpoint.
- **Frontend proxy method**: add a `createFeedback` proxy operation that maps to the backend feedback endpoint.
- **Schema**: add a `feedback` table with `id`, `user_id`, `username`, `email`, `url`, `message`, and `created_at`.
- **Identity source**: derive user id, username, and email on the backend from authenticated server-side context or trusted database lookup. The client must not be able to spoof sender identity or email.
- **Email requirement**: all users are expected to have email. If an unexpected authenticated account has no email, reject feedback with a clear error.
- **Technical context**: store only the current dashboard URL in the first version. Do not store user agent, site id, domain id, screenshots, or browser diagnostics.
- **Message validation**: trim feedback text, reject empty messages, and enforce a 5,000 character maximum on both client and server.
- **Admin model**: show a simple newest-first feedback list. Do not add read/resolved status, assignment, tagging, or reply threading.
- **Email operations**: all replies, welcome messages, and onboarding emails remain manual. No email provider integration is part of this plan.
- **Legal/contact stance**: keep the external Yandex feedback form for logged-out users while making the dashboard prefer the in-product feedback form.
- **Testing stance**: backend behavior should be covered with service, handler, repository, and migration tests where appropriate. Frontend validation is manual QA because the webapp does not use unit tests.

---

## Phase 1: Working Submission Path

**User stories**: 1, 2, 3, 8, 9, 10, 11, 12, 19, 20, 21, 28, 29, 30, 31, 32

### What To Build

Build the narrowest end-to-end path for a logged-in dashboard user to open a feedback modal, see their read-only account email, enter a message, submit it, and have the backend store a feedback record in PostgreSQL with server-derived user context.

This slice should include the feedback data model, authenticated create endpoint, frontend proxy/client operation, and a basic dashboard-shell entry point. The success state can be simple, but the flow must be demonstrable without manual database inserts.

### Acceptance Criteria

- [ ] Authenticated dashboard users can open a feedback modal from the dashboard shell.
- [ ] The modal shows the logged-in user's email as read-only.
- [ ] Submitting a non-empty message creates a database record with user id, username snapshot, email snapshot, message, current URL, and creation timestamp.
- [ ] The backend derives user identity and email; user id, username, and email from the client are ignored or not accepted.
- [ ] Unauthenticated requests to the feedback endpoint are rejected.
- [ ] Feedback records preserve the email snapshot from submission time.
- [ ] The frontend submission uses the existing Result-style API-client pattern.
- [ ] Backend service/handler/repository tests cover the happy path and unauthenticated rejection.

---

## Phase 2: Dashboard UX And Validation

**User stories**: 4, 5, 6, 7, 13, 14, 15, 16, 17, 18

### What To Build

Harden the user-facing dashboard feedback experience so it behaves like a polished production control rather than a bare form. The button should stay visible across desktop and mobile dashboard pages, the modal should be accessible, validation should be clear, and submission state should prevent duplicate sends while preserving typed text on errors.

This slice completes the user-facing form behavior while keeping the same backend contract from Phase 1.

### Acceptance Criteria

- [ ] The feedback button remains visible while scrolling dashboard pages.
- [ ] The button placement works on desktop and mobile without covering primary dashboard actions.
- [ ] The button and modal have clear accessible names and keyboard focus behavior.
- [ ] Empty feedback is rejected before submission with a clear message.
- [ ] Messages over 5,000 characters are rejected before submission with a clear message.
- [ ] The Send button is disabled while submission is in progress.
- [ ] Successful submission shows a clear received state.
- [ ] Failed submission preserves the typed message for retry.
- [ ] The modal does not hide failure state through accidental close-after-submit behavior.
- [ ] Manual QA covers desktop, mobile, keyboard navigation, screen-reader labels, success, validation errors, and retry after API failure.

---

## Phase 3: Admin Feedback Inbox

**User stories**: 22, 23, 24, 25, 26, 27

### What To Build

Add feedback visibility to the internal admin dashboard so operators can review recent messages and manually reply by email. Keep the admin model deliberately simple: newest-first records, sender context, message body, URL, and count context. Include feedback in environment wipe/reset behavior.

This slice makes submitted feedback operationally useful without building a ticket system.

### Acceptance Criteria

- [ ] The internal admin dashboard shows recent feedback ordered newest first.
- [ ] The admin list shows username, email, current URL, created time, and message body.
- [ ] Long feedback messages remain readable in the admin dashboard.
- [ ] The admin dashboard has a clear empty state when there is no feedback.
- [ ] Feedback count is visible either in dashboard metrics or in the feedback section header.
- [ ] Feedback query is limited to a practical recent window, such as the latest 100 records.
- [ ] Local wipe/reset behavior deletes feedback records with the rest of dashboard-managed data.
- [ ] Backend query or rendering behavior is tested where the existing dashboard structure allows it.

---

## Phase 4: Legal And Contact Copy

**User stories**: 33, 34, 35, 36, 37, 38

### What To Build

Update the legal and contact-facing copy so email processing and feedback storage are transparent. The privacy policy should mention collected emails and feedback content. The user agreement should describe manual service communications, including support replies, welcome/onboarding emails, and important service notices, while avoiding any implication of unrelated marketing without opt-in.

Keep the external Yandex feedback form available for logged-out users and update stale in-product references so authenticated users are directed toward the dashboard feedback form.

### Acceptance Criteria

- [ ] Privacy copy explicitly lists account email and feedback message content as collected data.
- [ ] Privacy copy explicitly includes support replies, onboarding, service notifications, and service improvement as processing purposes.
- [ ] Terms copy states that Hostamat may use the account email for manual service-related messages.
- [ ] Legal copy does not imply permission for unrelated promotional newsletters.
- [ ] Logged-out/public legal contact paths keep the external Yandex feedback form as a fallback.
- [ ] Authenticated dashboard feedback references prefer the in-product form.
- [ ] Existing console or public feedback hints are updated if they would otherwise point users to the wrong primary channel.
- [ ] Manual QA verifies legal pages and contact references are consistent.

---

## Phase 5: Final Verification

**User stories**: all

### What To Build

Run the full verification pass for the completed feature across backend tests, frontend checks, local manual flow, admin dashboard review, and legal-page review. Fix integration issues found during verification without expanding scope into email automation or ticket management.

### Acceptance Criteria

- [ ] Backend tests pass for feedback domain, migrations, handlers, and related repository behavior.
- [ ] Frontend lint/type/build checks pass using the repo's existing webapp commands.
- [ ] A logged-in user can submit feedback from the dashboard and see a success state.
- [ ] The stored feedback appears in the internal admin dashboard with sender email and URL.
- [ ] Empty, too-long, unauthenticated, and no-email edge cases fail clearly.
- [ ] Desktop and mobile dashboard layouts are manually checked.
- [ ] Privacy, terms, and contact references are manually checked for consistency.
- [ ] No automated email sending, status model, attachment flow, or diagnostics collection was added.
