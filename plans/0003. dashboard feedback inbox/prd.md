# PRD: Dashboard Feedback Inbox

## Problem Statement

Hostamat users currently have no first-class way to send feedback from inside the authenticated dashboard. The site still points users to an external feedback form in some places, which separates feedback from the logged-in account context and makes it harder to connect a message to a real user, email address, sites, domains, and recent activity.

This makes early alpha feedback harder to collect and act on. It also leaves Hostamat without a clear legal and product foundation for using the email address collected during Yandex OAuth to reply to feedback or send service-oriented onboarding messages such as a welcome email.

## Solution

Add an always-visible feedback button inside the authenticated dashboard shell. The button opens a modal with a read-only email field prefilled from the logged-in user's stored email, a feedback text area, and a Send button. Submitting the form saves the feedback to PostgreSQL and gives the user a clear success or error state without leaving the dashboard.

Add a backend feedback domain that validates authenticated submissions, stores the message with user and email context, and returns a small response contract. Show submitted feedback in the existing internal admin dashboard so operators can review recent feedback, identify the sender, and manually reply by email.

Update the legal pages so they explicitly mention email addresses, feedback messages, support replies, and service communications. The first version should not send emails automatically; all replies and welcome/onboarding messages remain manual. The legal copy should still make it clear that Hostamat may use stored user emails for direct service communication.

## User Stories

1. As a logged-in user, I want a feedback button that is always visible in the dashboard, so that I can contact Hostamat without searching for a link.
2. As a logged-in user, I want the feedback button to be available from the sites dashboard, so that I can report upload, hosting, or UX issues while I am working.
3. As a logged-in user, I want the feedback button to be available from the domains dashboard, so that I can report DNS or custom-domain issues in context.
4. As a logged-in user, I want the feedback button to remain visible after scrolling, so that I can open it at any point in the workflow.
5. As a mobile user, I want the feedback button to be reachable without covering primary dashboard actions, so that it does not make the dashboard harder to use.
6. As a keyboard user, I want the feedback button and modal controls to be focusable, so that I can submit feedback without a mouse.
7. As a screen-reader user, I want the feedback button and modal to have clear accessible names, so that the feedback flow is understandable.
8. As a logged-in user, I want the modal to show my account email, so that I know where Hostamat can reply.
9. As a logged-in user, I want the email field to be read-only, so that feedback is tied to my verified account email rather than an arbitrary address.
10. As a logged-in user, I want the modal to show the email from my account, so that I understand which address Hostamat can use for a reply.
11. As a logged-in user, I want to type a free-form feedback message, so that I can describe bugs, requests, questions, or confusion in my own words.
12. As a logged-in user, I want the feedback text area to support multiline text, so that I can include steps to reproduce or longer context.
13. As a logged-in user, I want the Send button disabled while my feedback is being submitted, so that I do not accidentally send duplicates.
14. As a logged-in user, I want a clear validation message when the feedback text is empty, so that I know what to fix.
15. As a logged-in user, I want a clear validation message when the feedback text is too long, so that the form fails predictably.
16. As a logged-in user, I want a success state after sending feedback, so that I know Hostamat received it.
17. As a logged-in user, I want errors to preserve my typed text, so that I can retry without retyping.
18. As a logged-in user, I want the modal to close only when I choose to close it or after a successful send state, so that accidental dismissal does not hide failures.
19. As a Hostamat operator, I want feedback saved in the database, so that it is not lost when a browser session or external form changes.
20. As a Hostamat operator, I want each feedback record to include the user id, username, stored email, message, and creation time, so that I can understand who wrote it and when.
21. As a Hostamat operator, I want feedback records to preserve the email value at submission time, so that later account email changes do not rewrite historical feedback context.
22. As a Hostamat operator, I want to see recent feedback in the internal admin dashboard, so that I can triage messages without querying the database manually.
23. As a Hostamat operator, I want the admin dashboard feedback list to show sender identity and email, so that I can manually reply.
24. As a Hostamat operator, I want long feedback messages to be readable in the admin dashboard, so that detailed reports are not truncated beyond usefulness.
25. As a Hostamat operator, I want feedback ordered newest first, so that the most recent messages are easiest to review.
26. As a Hostamat operator, I want feedback count included in dashboard metrics or admin context, so that I can notice whether users are writing in.
27. As a Hostamat operator, I want feedback data included in local wipe behavior, so that test environments can be reset consistently.
28. As a developer, I want feedback persistence isolated in a small domain module, so that storage, validation, and handler behavior can be tested without coupling to dashboard UI details.
29. As a developer, I want the frontend feedback client to use the existing Result pattern, so that submission errors match the rest of the webapp.
30. As a developer, I want the feedback route to require authentication, so that anonymous spam is not accepted by the authenticated dashboard endpoint.
31. As a developer, I want the backend to derive user identity and email from the session, so that the client cannot spoof sender identity.
32. As a developer, I want the feedback API contract to be small, so that future admin workflows can evolve without changing the user-facing form.
33. As a privacy-conscious user, I want the legal pages to say that email addresses are processed, so that account data usage is transparent.
34. As a privacy-conscious user, I want the legal pages to say that feedback messages are stored, so that I understand what happens when I submit the form.
35. As a privacy-conscious user, I want the legal pages to say that Hostamat may reply to feedback by email, so that I understand why my email is shown in the form.
36. As a privacy-conscious user, I want the legal pages to distinguish service communications from unrelated marketing, so that account emails are not used beyond the stated purposes.
37. As a Hostamat operator, I want the legal pages to cover manual welcome or onboarding emails, so that basic service communication is documented before I contact users.
38. As a Hostamat operator, I want the first version to keep outbound replies manual, so that feedback collection ships without email provider integration.

## Implementation Decisions

- The feedback entry point should live in the authenticated dashboard shell so it appears across dashboard pages, not inside a single page component.
- The feedback trigger should be visually persistent. On desktop it can be a fixed-position button near the lower-right dashboard area. On mobile it should remain reachable while avoiding overlap with primary navigation and content.
- The feedback modal should reuse the existing modal and button primitives.
- The modal copy should be written in Russian to match the authenticated dashboard.
- The modal should contain a read-only email field sourced from the authenticated user's stored email.
- The modal should contain one required multiline feedback text area.
- The frontend should not allow editing or overriding the email address.
- The backend should not trust any email value sent by the client. It should resolve the user and email from authenticated server-side context.
- All authenticated users are expected to have a stored email from the current auth flow.
- If an unexpected account state has no stored email, the feedback API should reject submission with a clear server-side error rather than accepting feedback without a reply address.
- Feedback submission should use the existing SvelteKit API proxy and typed API-client pattern.
- Add a dedicated backend feedback domain rather than placing feedback storage in the site domain, because feedback is account-level support data and not part of site lifecycle.
- The feedback domain should expose a create-feedback operation for authenticated users.
- The feedback record should store an id, user id, username snapshot, email snapshot, current dashboard URL, message body, and created timestamp.
- The first version should not add feedback categories, severity, status, read/resolved state, assignment, or threaded replies.
- Message validation should be enforced on both frontend and backend.
- The backend should trim leading and trailing whitespace before storing feedback.
- Empty messages should be rejected.
- A maximum message length should be enforced. A reasonable first-version limit is 5,000 characters.
- Duplicate prevention should be limited to UI-level disabled sending in the first version. Server-side deduplication is out of scope unless spam becomes a real issue.
- The internal admin dashboard should query the latest feedback records and render them in a dedicated section.
- Admin feedback should be ordered by creation time descending and limited to a practical recent window, such as the latest 100 messages.
- The admin dashboard should show sender username, sender email, created time, and message body.
- The admin dashboard should include a clear empty state when no feedback exists.
- Feedback should be included in admin dashboard counts if the existing metric row has room; otherwise the feedback section header should show the count loaded for that page.
- The database wipe operation used by the internal dashboard should delete feedback records together with the rest of dashboard-managed data.
- Feedback migrations should be managed through the same platforma migration conventions used by existing domains.
- Privacy policy updates should explicitly list email address and feedback message content as collected data.
- Privacy policy purpose updates should explicitly include replying to feedback, technical support, onboarding, service notifications, and improving Hostamat.
- Terms or user agreement updates should mention that Hostamat may send service-related messages to the email connected to the account, including replies to support requests and onboarding messages.
- Legal copy should avoid implying consent for unrelated promotional newsletters unless a separate opt-in mechanism is implemented.
- Existing external Yandex Forms references should be reviewed. The dashboard should prefer the in-product feedback form after this feature ships, while legal and public contact links should keep the external form as a fallback for logged-out users.
- The implementation should preserve Sentry console feedback hints only if they remain useful; otherwise they should be updated to reference the in-product feedback flow.
- Outbound email sending is not part of the first implementation. Manual replies and manual welcome/onboarding emails using the stored user email are sufficient for the first version.
- No email provider or sender address is selected in this PRD.
- The data model and legal copy should be compatible with manual service emails now and a future email delivery provider later.

## Testing Decisions

- Good tests should verify externally visible behavior: authenticated users can submit feedback, unauthenticated requests are rejected, invalid messages are rejected, valid messages are stored with server-derived user context, and admin dashboard data includes saved feedback.
- Backend service tests should cover message validation, trimming, max length enforcement, and storage of user/email snapshots.
- Backend handler tests should cover request decoding, authentication behavior, success response shape, and expected error codes.
- Repository tests should cover the migration and the ability to insert and list feedback records ordered newest first.
- Admin dashboard behavior can be covered with focused query or rendering checks if the current dashboard test style supports it; otherwise it should be manually verified against local test infrastructure.
- Frontend testing should focus on manual QA because this webapp does not use unit tests.
- Manual QA should verify desktop and mobile dashboard layouts, modal focus behavior, read-only email display, current URL capture, empty message validation, long message validation, successful submission, retry after API failure, and visibility of the submitted record in the admin dashboard.
- Legal-page QA should verify that privacy, terms, and any contact references consistently mention email usage, feedback storage, support replies, and service communications.

## Out of Scope

- Sending automated emails in the first version.
- Building a full support ticketing system.
- Marking feedback as read, resolved, assigned, archived, or tagged.
- Threaded admin replies inside Hostamat.
- User-facing history of previously submitted feedback.
- Anonymous feedback from logged-out visitors in the in-product form.
- Attaching screenshots, files, site ids, domain ids, browser user agents, or other browser diagnostics to feedback submissions.
- Feedback categories, ratings, NPS, or surveys.
- Spam scoring, rate limits, or abuse tooling beyond authentication and basic validation.
- Marketing newsletters or promotional campaigns.
- Editing the user's stored email from the feedback modal.
- Replacing all legal contact mechanisms if an external public contact form is still needed for unauthenticated users.

## Open Questions

1. What exact Russian legal wording should be used to describe manual email contact: support replies, welcome/onboarding emails, and important service notices?
2. Does the Russian legal copy need review by counsel before shipping email communication language?

## Further Notes

- The current authenticated user model already stores email addresses from Yandex OAuth when Yandex provides them.
- The current internal admin dashboard already joins and displays user emails for users, sites, custom domains, and Caddy host context.
- The current privacy page lists Yandex account identifier/login, session cookies, site/file information, and technical data, but it does not explicitly list email addresses or feedback message content.
- The current terms page does not explicitly describe service emails, feedback replies, or welcome/onboarding emails.
- The current public layout logs an external Yandex Forms feedback link to the browser console. This should be revisited once an in-product feedback form exists.
- Product decisions from the PRD interview: all users are expected to have email; feedback should store the current URL only as technical context; the external Yandex feedback form stays available for logged-out users; all email activity remains manual for now; no feedback status model should be added in the first version.
- Because the user asked for legal-page review, the PRD intentionally calls out legal copy changes but does not provide legal advice. Final wording should be reviewed for the target jurisdiction and actual email behavior before production use.
