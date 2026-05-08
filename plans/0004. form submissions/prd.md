# PRD: Form Submissions for Hosted Static Sites

## Problem Statement

Hostamat users can upload static sites, but static hosting does not give ordinary HTML forms anywhere useful to submit data. A user may upload a landing page, portfolio, event page, waitlist, or contact page that contains one or more `<form>` tags, then only discover after sharing the site that form submissions do not reach them.

Users need Hostamat to notice forms during upload, explain that form collection can be enabled, and provide a simple dashboard inbox for received submissions. The feature should preserve the simplicity of static-site hosting: users should not need to write backend code, configure mail delivery, or manually edit their HTML before upload.

## Solution

When an upload session completes successfully, Hostamat will scan uploaded HTML files for form tags. If any forms are found, the upload result UI will show a prompt asking whether the user wants Hostamat to collect submissions from those forms.

If the user opts in for selected forms, Hostamat will rewrite the uploaded HTML so those supported forms post to a Hostamat-managed endpoint. The endpoint will accept public submissions from the hosted site, store only the submitted form fields in PostgreSQL, and show a simple success response. The site owner will see submissions in a new dashboard section tied to their sites and forms, and can delete individual submissions when they are no longer needed.

The first version should focus on ordinary server-submitted HTML forms. JavaScript-heavy forms, custom client-side validation flows, file upload fields, email notifications, and spam protection beyond basic safeguards are out of scope for the first release.

## User Stories

1. As a site owner, I want Hostamat to detect forms in uploaded HTML, so that I know my static site includes interactive forms.
2. As a site owner, I want form detection to run automatically after upload, so that I do not have to start a separate scan.
3. As a site owner, I want the upload result to tell me which uploaded pages contain forms, so that I understand what Hostamat found.
4. As a site owner, I want the upload result to show how many forms were found, so that I can make a quick decision.
5. As a site owner, I want Hostamat to ask before modifying my HTML, so that my uploaded site is not changed without consent.
6. As a site owner, I want to choose which detected forms should collect replies, so that only the forms I approve are modified.
7. As a site owner, I want to decline form collection, so that existing third-party form handlers or custom behavior remain untouched.
8. As a site owner, I want declined forms to keep their original HTML, so that upload behavior is predictable.
9. As a site owner, I want enabled forms to post to a Hostamat endpoint, so that submissions can be stored without my own backend.
10. As a site owner, I want existing form fields to keep their names, so that submitted data remains understandable.
11. As a site owner, I want text inputs, text areas, selects, checkboxes, radio buttons, hidden fields, and submit values to be captured, so that common contact and lead forms work.
12. As a site owner, I want file inputs to be rejected or ignored clearly in the first version, so that users are not surprised when attachments are not stored.
13. As a site owner, I want the rewritten form to include a stable form identifier, so that submissions can be grouped by form.
14. As a site owner, I want detected forms to have useful display names when possible, so that I can distinguish a contact form from a newsletter form.
15. As a site owner, I want anonymous or unnamed forms to still work, so that simple HTML from site builders is supported.
16. As a site owner, I want forms in nested HTML pages to be supported, so that `/contact` and `/landing/signup` can collect responses too.
17. As a site owner, I want forms in ZIP uploads and individual file uploads to be handled consistently, so that the source upload method does not matter.
18. As a site owner, I want form collection status to be visible after upload, so that I know whether collection is enabled for a site.
19. As a site owner, I want to see all submissions for my sites in a dashboard section, so that I can act on messages and leads.
20. As a site owner, I want submissions grouped or filterable by site, so that data from different hosted sites does not mix together.
21. As a site owner, I want submissions grouped or filterable by form, so that multiple forms on one site stay understandable.
22. As a site owner, I want each submission to show the page URL, form name, submitted fields, and received time, so that I can understand where it came from.
23. As a site owner, I want submissions ordered newest first, so that recent responses are easiest to review.
24. As a site owner, I want long values to be readable in the dashboard, so that detailed messages are not lost.
25. As a site owner, I want empty values and unchecked options represented clearly, so that the stored response matches what the visitor submitted.
26. As a site owner, I want a clear empty state when there are no submissions, so that I understand collection is enabled but no one has replied yet.
27. As a site owner, I want only my own site submissions to be visible to me, so that another user's form data is private.
28. As a site owner, I want deleting a site to delete or otherwise remove access to its form definitions and submissions, so that stale personal data does not remain visible.
29. As a site owner, I want to delete individual submissions, so that I can remove outdated, test, spam, or personal data from my dashboard.
30. As a site visitor, I want submitting a form to work without creating a Hostamat account, so that public website forms remain public.
31. As a site visitor, I want the submission response to be understandable, so that I know whether my message was sent.
32. As a site visitor, I want basic browser form behavior to keep working, so that required fields and input types still help me submit valid data.
33. As a site visitor, I want my submission not to expose Hostamat internals in the page, so that the form feels like part of the hosted site.
34. As a developer, I want HTML scanning and rewriting isolated in a testable module, so that parsing rules do not leak into handlers or Svelte components.
35. As a developer, I want form collection persistence isolated in a domain module, so that submission storage, access checks, and deletion can be tested independently.
36. As a developer, I want the public submission endpoint to resolve the site and form from server-side identifiers, so that clients cannot submit into another user's inbox by changing a site id.
37. As a developer, I want the opt-in operation to be explicit and authenticated, so that only the site owner can enable collection.
38. As a developer, I want upload scanning to reuse the upload session lifecycle, so that detection is based on the files that were just uploaded.
39. As a developer, I want public submission limits, body-size limits, and field-count limits, so that a public endpoint cannot accept unbounded payloads.
40. As a privacy-conscious user, I want the product to make clear that form replies are stored by Hostamat, so that collection is transparent.
41. As a privacy-conscious user, I want legal pages to explain hosted form submission storage, so that visitor data handling is documented.
42. As a Hostamat operator, I want collected form data included in local wipe/user cleanup behavior, so that test and deletion workflows stay consistent.

## Implementation Decisions

- Form detection should run after a successful upload session close, alongside the existing upload summary behavior.
- The close-upload response should include a structured form discovery summary when uploaded HTML contains forms.
- The discovery summary should include the total number of supported forms, the pages where forms were found, and whether any forms could not be supported in the first version.
- Hostamat must not rewrite HTML automatically during the initial upload close operation. Rewriting happens only after an authenticated owner opt-in action.
- The upload result UI should show a prompt when forms are detected. The prompt should explain that Hostamat can collect form replies and that enabling it modifies form actions in uploaded HTML.
- The prompt copy should be written in Russian to match the current authenticated dashboard.
- The first version should let the user choose individual supported detected forms from the latest upload before enabling collection.
- The opt-in UI should show each supported detected form with its page, derived label, field summary, and current action if one exists.
- If a site owner declines or closes the prompt, Hostamat should leave the uploaded files as-is and should not create active collection endpoints for those forms.
- If a form already has an external `action` attribute, enabling collection for that specific form should overwrite it after explicit opt-in.
- If a form has no `method`, Hostamat should treat it as `GET` in discovery but rewrite supported collected forms to `POST`.
- If a form uses `method="post"`, Hostamat should preserve the user intent while changing the `action` to the Hostamat collection endpoint.
- Forms using unsupported methods should be listed as unsupported and not rewritten in the first version.
- Forms with `enctype="multipart/form-data"` or file inputs should be listed as unsupported in the first version unless all file fields can be ignored with clear dashboard messaging.
- Form field names should come from HTML `name` attributes. Fields without names are not submitted by normal browser form behavior and should not be stored.
- The scanner should detect ordinary `<form>` elements in HTML files using a real HTML parser, not regular expressions.
- The rewriter should preserve the rest of the HTML document as much as practical, including non-form markup, scripts, styles, and text.
- Each discovered form should receive a stable internal form id tied to the site, page path, and form position or derived signature.
- Rewriting should add the Hostamat collection endpoint as the form `action`.
- Rewriting should add any hidden metadata needed by the backend, such as a form id token or return-path token, without trusting client-provided ownership values.
- The public submission endpoint should accept unauthenticated browser form posts.
- The public submission endpoint should identify the target site and form from an opaque token or route value generated during opt-in.
- The public submission endpoint should verify that the target form is active before storing a submission.
- Submission storage should include site id, form id, page path, submitted field payload, and creation time.
- The first version should not store request metadata such as IP address, hashed IP address, user agent, or referrer.
- Field payload should be stored as structured JSON so arbitrary form schemas do not require migrations.
- A maximum request body size should be enforced. A reasonable first-version limit is 64 KB per submission.
- A maximum number of submitted fields should be enforced. A reasonable first-version limit is 100 fields per submission.
- A maximum field value length should be enforced. A reasonable first-version limit is 10,000 characters per field.
- The first version should not send email notifications when a form is submitted.
- The first version should not add webhook delivery or integrations.
- The submission response should use a simple Hostamat-generated success page in the first version.
- Configurable post-submit redirects are a likely later enhancement, but are not part of v1.
- If the original form uses an existing redirect target, success-page behavior should be preferred for the first version.
- Failed public submissions should return a simple user-facing error page rather than raw JSON.
- The authenticated dashboard should add a new section for form submissions.
- The submissions section should use the existing dashboard shell and API-client Result pattern.
- The submissions section should list submissions newest first and allow filtering by site.
- If multiple forms exist on a site, the submissions section should allow filtering by form or show form labels in each row.
- Submission detail should show submitted fields in a readable key/value layout.
- The dashboard should not show raw internal tokens to the site owner unless they are useful for support.
- The dashboard should allow deleting individual submissions.
- Deleting a submission should remove it from the owner's inbox and from persistent storage.
- The backend should provide authenticated APIs to list form definitions for the user's sites and list submissions for owned sites.
- The backend should provide an authenticated API to delete individual submissions owned by the current user.
- The backend must enforce ownership for every authenticated form/submission list request.
- Form definitions and submissions should be deleted when the owning site is deleted.
- User cleanup should delete form definitions and submissions for all of the user's sites.
- Existing uploaded-page discovery should continue to work; form discovery should extend it without removing page links or missing-home-page warnings.
- Form collection should be disabled or replaced on each new upload session because new files replace old files. The latest upload should be the source of truth for active forms.
- If a user uploads a new site version after enabling collection, Hostamat should rescan forms and require a new opt-in before rewriting the new HTML.
- The feature should not require scanning S3 outside the upload flow for the first version unless the opt-in rewrite needs to fetch and update uploaded HTML objects.
- The implementation should be compatible with generated subdomains and custom domains, because public submissions may originate from either host.
- Legal/privacy copy updates are part of this feature because hosted site visitors may submit personal data into Hostamat storage.
- Legal/privacy copy should explain that site owners can enable Hostamat-managed form collection and that Hostamat stores submitted form field data.

## Testing Decisions

- Good tests should verify externally visible behavior: uploaded HTML with forms produces a discovery summary, opting in rewrites supported forms, public submissions are stored, and dashboard APIs show only the owner's submissions.
- Backend scanner tests should use realistic HTML fixtures and cover one form, multiple forms, nested pages, malformed-but-browser-tolerated HTML, forms without action, forms with existing action, forms without method, unsupported methods, file inputs, unnamed fields, and non-HTML files.
- Backend rewrite tests should verify output behavior rather than byte-for-byte formatting: rewritten forms should contain the expected action, method, hidden metadata, and preserved fields.
- Submission service tests should cover valid posts, inactive form tokens, unknown form tokens, request-size limits, field-count limits, field-value limits, duplicate field names, checkbox/radio values, and ownership checks for listing.
- Repository tests should cover migrations, creating form definitions, activating/deactivating forms, storing JSON payloads, listing submissions newest first, filtering by site/form, deleting individual submissions, and deleting by site/user cleanup.
- Handler tests should cover authenticated opt-in, unauthenticated opt-in rejection, public submission success, public submission validation failures, authenticated list behavior, and cross-user access denial.
- Frontend testing should focus on manual QA because this webapp does not use unit tests.
- Manual QA should upload a site with no forms, a site with one contact form, a site with multiple forms across pages, a site with an existing third-party form action, a ZIP with nested HTML pages, and a site with unsupported file-upload forms.
- Manual QA should verify desktop and mobile dashboard layouts, the upload prompt, individual form selection, the opt-in flow, successful visitor submission, error visitor submission, empty submission inbox, filtered inbox, individual submission deletion, and site deletion cleanup.
- Security QA should verify that public submissions cannot choose arbitrary site ids, that authenticated users cannot view another user's form submissions, and that oversized requests fail predictably.

## Out of Scope

- Email notifications for new submissions.
- SMTP, SendGrid, Telegram, webhook, CRM, or spreadsheet integrations.
- File upload fields and attachment storage.
- CAPTCHA, advanced bot detection, spam scoring, or moderation queues.
- Per-form custom confirmation messages.
- Per-form redirect configuration in v1.
- Editing form schemas in the dashboard.
- Building a drag-and-drop form builder.
- JavaScript SDKs for custom client-side forms.
- Parsing or executing JavaScript to discover dynamically generated forms.
- Supporting non-HTML page formats.
- Exporting submissions to CSV in the first version.
- Marking submissions as read, archived, replied, assigned, or resolved.
- Visitor-facing consent management beyond preserving the site's own visible form copy and updating Hostamat legal pages.
- Automatic rewriting without explicit site-owner opt-in.
- Retrofitting existing already-uploaded sites outside the next upload flow.

## Further Notes

- Current upload behavior replaces old site files when a new upload session is created, then returns an upload summary when the session is closed.
- Current upload summary already discovers uploaded HTML pages and warns when root `index.html` is missing. Form detection can extend that backend summary pattern.
- Current public file serving resolves the host to a site and serves from object storage. Public form submission should similarly resolve/validate site and form ownership server-side instead of trusting submitted ids.
- Current dashboard copy is Russian, while existing local PRDs are written in English. This PRD follows the existing PRD language pattern but calls out Russian UI copy explicitly.
- Form collection is opt-in after each upload because the upload replaces HTML and form structure may change.
- Product decisions from follow-up discussion: users choose individual forms; existing form actions are overwritten only for selected forms; v1 uses a Hostamat success page while configurable redirects can come later; submissions store only form fields and timestamps; individual submission deletion is included; legal/privacy updates are included.
