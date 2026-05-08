# Plan: Form Submissions for Hosted Static Sites

> Source PRD: `plans/0004. form submissions/prd.md`

## Architectural Decisions

Durable decisions that apply across all phases:

- **User-facing upload surface**: form discovery is shown in the existing upload completion UI after an upload session is closed.
- **Form enablement model**: users choose individual detected supported forms. Hostamat rewrites only the selected forms.
- **Existing form actions**: if a selected form already has an `action`, Hostamat overwrites it after explicit opt-in.
- **Form parser/rewriter**: HTML form scanning and rewriting are backend responsibilities and use a real HTML parser rather than regular expressions.
- **Backend domain**: add a dedicated form collection domain for form definitions, opt-in state, submissions, authenticated dashboard APIs, and public submission handling.
- **Schema**: store form definitions separately from submissions. Form definitions belong to a site and include page path, display label, original action, method, field summary, active state, and opaque public token. Submissions belong to a form and site and store submitted fields as JSON plus creation time.
- **Submission data**: store only submitted form fields and timestamps. Do not store IP address, hashed IP, user agent, or referrer in v1.
- **Public route**: expose an unauthenticated public submission endpoint addressed by an opaque form token. The endpoint verifies the token and active form before storing anything.
- **Authenticated routes**: expose authenticated APIs for listing form definitions, enabling selected forms from the latest discovery result, listing submissions, and deleting individual submissions.
- **Authorization**: all authenticated form/submission operations derive ownership from server-side user context and site ownership. Clients cannot choose arbitrary site ids or form ids to access another user's data.
- **Visitor response**: successful public submissions show a simple visually neutral success page in v1 so the response does not feel strongly branded as Hostamat. Configurable redirects are a later enhancement.
- **Upload replacement rule**: every new upload replaces site files, rescans forms, and requires fresh opt-in before rewriting the new HTML.
- **Historical submissions**: submissions from previous uploads remain visible after re-upload. Old form definitions are disabled for new submissions but retained as historical form context for existing submissions.
- **Form labels**: form labels are auto-derived only in v1. Users cannot rename forms yet.
- **Dashboard route**: add a new authenticated dashboard section for form submissions alongside sites and domains.
- **Dashboard inbox layout**: use a table-first layout for scanning submissions, with site/form/page/time columns and expandable or drill-in field details.
- **Legal scope**: legal/privacy copy updates are part of the feature because hosted site visitors may submit personal data into Hostamat storage.
- **Testing stance**: backend scanning, rewriting, persistence, handlers, and ownership checks need automated tests. Frontend behavior is verified manually because the webapp does not use unit tests.

---

## Phase 1: Form Discovery in Upload Results

**User stories**: 1, 2, 3, 4, 5, 7, 8, 12, 14, 15, 16, 17, 34, 38

### What to build

Add the first end-to-end tracer for detecting forms during upload without rewriting anything. After upload close, Hostamat scans uploaded HTML files, returns a structured form discovery summary, and the upload result UI shows what was found with a clear opt-in prompt placeholder.

This phase makes the feature visible and testable while preserving current uploaded files exactly as-is.

### Acceptance criteria

- [ ] Closing an upload session scans successfully uploaded HTML files for ordinary `<form>` elements.
- [ ] The upload close response includes a form discovery summary when forms are found.
- [ ] The summary includes form count, page path, derived form label, method, current action, field summary, and unsupported reason when applicable.
- [ ] Forms in root HTML files, nested HTML files, ZIP uploads, and individually uploaded HTML files are discovered consistently.
- [ ] Forms with unsupported methods, file inputs, or multipart encoding are identified as unsupported rather than silently enabled.
- [ ] Uploads with no forms return an empty discovery summary without changing the existing upload result behavior.
- [ ] No HTML is rewritten during upload close.
- [ ] The upload result UI shows a Russian prompt when supported forms are detected.
- [ ] Declining or closing the prompt leaves uploaded files unchanged.
- [ ] Backend scanner tests cover realistic HTML fixtures, malformed-but-browser-tolerated HTML, unnamed fields, existing actions, missing methods, unsupported forms, and non-HTML files.

---

## Phase 2: Per-Form Opt-In and HTML Rewrite

**User stories**: 6, 7, 8, 9, 10, 11, 13, 14, 15, 16, 17, 18, 32, 34, 35, 37, 38

### What to build

Let the user choose individual supported forms from the upload result and enable collection for only those forms. The backend creates active form definitions, rewrites selected forms in the uploaded HTML to post to Hostamat, and leaves unselected or unsupported forms untouched.

This phase is complete when a user can upload HTML with multiple forms, select one form, and verify that only the selected form's HTML now posts to Hostamat.

### Acceptance criteria

- [ ] The upload result UI lists supported detected forms individually with page, label, field summary, and current action.
- [ ] The user can select one or more supported forms and enable collection.
- [ ] Form labels are auto-derived from the HTML/page context and are not editable in v1.
- [ ] Unsupported forms cannot be selected and explain why collection is unavailable.
- [ ] Enabling collection requires an authenticated owner request.
- [ ] The backend creates active form definitions for selected forms only.
- [ ] Selected forms are rewritten to use `method="post"` and a Hostamat collection endpoint.
- [ ] Existing `action` attributes are overwritten only for selected forms.
- [ ] Unselected supported forms keep their original HTML.
- [ ] Unsupported forms keep their original HTML.
- [ ] Rewritten forms include stable hidden metadata or route data needed to identify the active form without trusting client ownership values.
- [ ] The rewrite preserves the rest of the document behavior as much as practical.
- [ ] Form collection status is visible after enablement.
- [ ] Rewriter tests verify behavior by parsed output: expected action, method, metadata, preserved fields, and untouched unselected forms.
- [ ] Handler/service/repository tests cover opt-in success, unauthenticated rejection, cross-user rejection, and unsupported form rejection.

---

## Phase 3: Public Submission Storage

**User stories**: 9, 10, 11, 13, 24, 25, 30, 31, 32, 33, 35, 36, 39, 40

### What to build

Add the public submission endpoint and storage path. A visitor can submit an enabled form without authentication, Hostamat validates the active form token, stores only form field data and timestamp, enforces payload limits, and returns a simple visually neutral success page.

This phase proves the hosted form is useful before building the owner inbox.

### Acceptance criteria

- [ ] Public visitors can submit enabled rewritten forms without logging in.
- [ ] The public endpoint identifies the active form by opaque token or generated route value.
- [ ] Unknown, inactive, or disabled form tokens are rejected without storing data.
- [ ] Submitted fields are stored as structured JSON.
- [ ] Duplicate field names, checkbox/radio values, select values, text areas, hidden fields, and submit values are handled predictably.
- [ ] Fields without names are not stored.
- [ ] Request body size, field count, and field value length limits are enforced.
- [ ] Stored submissions include site id, form id, page path, field payload, and created timestamp.
- [ ] Stored submissions do not include IP address, hashed IP address, user agent, or referrer.
- [ ] Successful submissions return a simple user-facing success page with visually neutral styling.
- [ ] Failed submissions return a simple user-facing error page, not raw JSON.
- [ ] Backend service/handler/repository tests cover valid posts, invalid tokens, inactive forms, limits, and common field types.

---

## Phase 4: Dashboard Submissions Inbox

**User stories**: 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 29, 35

### What to build

Add the authenticated dashboard section where site owners can review collected submissions. The page uses a table-first layout, lists submissions newest first, supports filtering by site and form, shows readable field payload details, and lets users delete individual submissions.

This phase makes collected data operationally useful for the site owner.

### Acceptance criteria

- [ ] The authenticated dashboard navigation includes a submissions section.
- [ ] The submissions section loads only submissions owned by the current authenticated user.
- [ ] Submissions are displayed in a table-first layout ordered newest first.
- [ ] Users can filter submissions by site.
- [ ] Users can filter submissions by form.
- [ ] The table includes enough context for multiple forms: site, auto-derived form label, page path, received time, and a compact submitted-field preview.
- [ ] Full submitted-field details are available through row expansion, a detail panel, or a detail view.
- [ ] Each submission shows site, form label, page path, received time, and submitted fields.
- [ ] Historical submissions from disabled previous-upload forms remain visible with clear historical or inactive form context.
- [ ] Long values remain readable.
- [ ] Empty values and unchecked option behavior are represented clearly.
- [ ] The page has a clear empty state when there are no submissions.
- [ ] Users can delete individual submissions.
- [ ] Deleted submissions disappear from the inbox and are removed from persistent storage.
- [ ] Authenticated APIs enforce ownership for listing and deletion.
- [ ] Frontend API calls use the existing Result pattern and webapp proxy style.
- [ ] Backend tests cover listing, filtering, deletion, and cross-user access denial.
- [ ] Manual QA covers desktop and mobile inbox layouts, filtering, detail readability, empty state, and deletion.

---

## Phase 5: Lifecycle Cleanup and Re-Upload Behavior

**User stories**: 18, 27, 28, 35, 37, 42

### What to build

Close the lifecycle gaps: new uploads disable old form definitions for future submissions while preserving historical submissions, site deletion removes form data, user cleanup removes form data, and local wipe/admin reset behavior includes form collection records.

This phase prevents stale form endpoints from accepting new submissions while keeping old submissions understandable in the dashboard until the owner deletes them or the owning site/user is deleted.

### Acceptance criteria

- [ ] Starting or completing a new upload for a site prevents previous form endpoints from continuing to collect against stale HTML.
- [ ] A new upload requires fresh form discovery and fresh per-form opt-in before collection is active again.
- [ ] Re-uploading a site with changed forms does not expose stale form labels, fields, or endpoints as active.
- [ ] Submissions collected before re-upload remain visible after re-upload.
- [ ] The dashboard distinguishes historical submissions from currently active form definitions when needed.
- [ ] Deleting a site deletes or deactivates its form definitions and deletes its submissions.
- [ ] User cleanup deletes form definitions and submissions for all of the user's sites.
- [ ] Local wipe/reset behavior includes form collection records.
- [ ] Lifecycle behavior is covered by service/repository tests.
- [ ] Manual QA verifies upload, enable, submit, re-upload, and stale endpoint rejection.

---

## Phase 6: Legal and Privacy Copy

**User stories**: 40, 41

### What to build

Update public legal/privacy copy to explain Hostamat-managed form collection. Draft the first version directly in product/legal copy rather than using placeholders. The copy should make clear that site owners can enable collection, that Hostamat stores visitor-submitted form field data, and that visitors may submit personal data through hosted forms.

This phase should stay factual and aligned with actual v1 behavior: no request metadata storage, no email notifications, no marketing claims, and no file attachments.

### Acceptance criteria

- [ ] Privacy copy states that Hostamat may store form field data submitted through hosted sites when the site owner enables collection.
- [ ] Privacy copy accurately states that v1 form submissions store submitted fields and timestamps.
- [ ] Legal copy does not claim that Hostamat stores IP address, user agent, referrer, attachments, or email delivery metadata for this feature.
- [ ] Legal copy does not imply unrelated marketing use.
- [ ] Legal copy is written as a first usable version, not as placeholders.
- [ ] Any public or dashboard-facing copy about form collection is consistent with the legal pages.
- [ ] Manual QA verifies relevant legal pages render correctly on desktop and mobile.

---

## Phase 7: Final Verification

**User stories**: all

### What to build

Run the full verification pass for the complete feature across backend tests, frontend checks, local manual flows, public submission behavior, dashboard inbox behavior, lifecycle cleanup, and legal-page review. Fix integration issues without expanding scope into redirects, email notifications, spam tooling, exports, attachments, or form builders.

### Acceptance criteria

- [ ] Backend tests pass for scanner, rewriter, form collection domain, migrations, handlers, repositories, ownership checks, and lifecycle cleanup.
- [ ] Frontend lint/type/build checks pass using the repo's existing webapp commands.
- [ ] Manual QA verifies upload with no forms.
- [ ] Manual QA verifies upload with one supported contact form.
- [ ] Manual QA verifies upload with multiple supported forms and individual selection.
- [ ] Manual QA verifies upload with an existing third-party form action and selected-form overwrite.
- [ ] Manual QA verifies upload with unsupported file-upload forms.
- [ ] Manual QA verifies public visitor submission success and public submission error behavior.
- [ ] Manual QA verifies dashboard inbox filtering, readable submission details, empty state, and deletion.
- [ ] Manual QA verifies re-upload disables stale form endpoints and requires fresh opt-in.
- [ ] Manual QA verifies historical submissions remain visible after re-upload and are clearly distinguishable from active forms.
- [ ] Manual QA verifies site deletion cleanup.
- [ ] Manual QA verifies legal/privacy copy consistency.
- [ ] No configurable redirects, email notifications, file attachments, CAPTCHA, export flow, submission status model, or JavaScript form SDK were added.
