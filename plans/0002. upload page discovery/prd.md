# PRD: Uploaded Page Discovery

## Problem Statement

Hostamat users often upload a static site and do not understand which URLs are now available. The current upload success state tells users which files were uploaded or skipped, but it does not translate uploaded HTML files into visitable pages. Users also may upload a site without a root `index.html`, then open the site home page and see a missing page without understanding why.

This creates confusion immediately after upload, especially for users who upload ZIP archives from site builders or manually selected files. They need Hostamat to explain what happened in terms of pages they can open, not only files that were stored.

## Solution

After every completed upload, Hostamat will analyze the uploaded files on the backend when the upload session is closed and identify the uploaded HTML files. The close-upload response will return a structured upload summary. The upload result UI will show the user a clear list of available pages as regular links with full URLs for the site.

If the uploaded files do not include a root `index.html`, Hostamat will show a warning explaining that the site does not have a home page at `/`, and that the user should upload `index.html` in the root of the site to make the main URL work.

The first version should focus on post-upload feedback in the dashboard upload modal only. It should extend the existing upload session close response, because the backend already has the uploaded and skipped file records needed to build this summary.

## User Stories

1. As a site owner, I want to see which pages were created after upload, so that I know what I can open immediately.
2. As a site owner, I want page links to be shown as URLs, so that I do not have to manually convert filenames into paths.
3. As a site owner, I want available pages to be shown as regular links, so that I can open them without learning any additional controls.
4. As a site owner, I want to know whether my home page exists, so that I understand whether the main site URL will work.
5. As a site owner, I want a warning when `index.html` is missing from the root, so that I can fix the upload package before sharing the site.
6. As a site owner, I want the warning to distinguish root `index.html` from nested `index.html` files, so that `about/index.html` does not make me think the home page exists.
7. As a site owner, I want only successfully uploaded HTML files to become available pages, so that skipped or invalid files do not appear as working pages.
8. As a site owner, I want skipped files to remain visible separately, so that I can understand what was ignored.
9. As a site owner, I want uploaded CSS, JS, image, and font files not to clutter the page list, so that the result focuses on navigable pages.
10. As a site owner, I want nested HTML files to be represented as their public page paths, so that `about/index.html` appears as `/about` when that clean path resolves to it.
11. As a site owner, I want simple HTML files to be represented as clean page paths, so that `about.html` appears as `/about`.
12. As a site owner, I want the root `index.html` to be shown as the home page, so that I understand the main domain URL maps to that file.
13. As a site owner, I want the upload result to handle ZIP uploads, so that extracted HTML files are discovered automatically.
14. As a site owner, I want the upload result to handle multiple individually selected files, so that manually selected HTML files are discovered too.
15. As a site owner, I want the message to remain visible after upload completes, so that I can review results before closing the modal.
16. As a site owner, I want upload errors to remain clearly separate from page discovery, so that failed uploads do not appear successful.
17. As a support operator, I want users to see actionable upload feedback, so that fewer users ask why their site URL does not work.
18. As a developer, I want the page discovery logic isolated and testable, so that URL mapping and `index.html` detection do not become fragile UI-only string logic.
19. As a developer, I want backend and frontend contracts to stay small, so that the upload flow remains easy to maintain.
20. As a developer, I want the feature to use the existing upload session lifecycle, so that no separate file listing endpoint is required for the first version.

## Implementation Decisions

- The feature is triggered when an upload session is closed successfully.
- Page discovery logic must live in the backend.
- The close-upload response should return a structured upload summary instead of only raw file records.
- The upload summary should include uploaded files, skipped files, available pages, and whether the root home page is present.
- Page discovery considers only files with `uploaded: true`.
- Page discovery considers only files whose normalized names end in `.html`.
- The root home page exists only when there is an uploaded file whose normalized name is exactly `index.html`.
- Nested `index.html` files are treated as nested pages, not as the root home page.
- The upload result should include a dedicated “Available pages” section separate from the existing uploaded files and skipped files sections.
- Available pages should be rendered as regular links, not as custom action buttons.
- The missing-home-page warning should appear when the upload succeeds and root `index.html` is absent.
- The missing-home-page warning should be displayed as a separate warning-colored alert card, distinct from the success alert.
- The warning should be actionable and concise: users should be told to place `index.html` at the root of the uploaded site or ZIP archive.
- Page URLs should be derived from the selected site domain and the same routing behavior used by public file serving.
- The backend should return both a display path and a full absolute URL for each available page.
- For `index.html`, the available page URL is the site root without a trailing slash.
- For `folder/index.html`, the preferred available page URL is the folder path without a trailing slash.
- For `page.html`, the available page URL should be the clean extensionless path, for example `/about` for `about.html`.
- Page URL generation must match current public file resolution precedence: extension paths resolve exactly, while extensionless paths try exact path, then `path.html`, then `path/index.html`.
- If two uploaded files would compete for the same clean URL, the file that public serving would resolve for that URL gets the clean URL; the other file should be linked by an exact path that resolves to that file.
- HTML extension matching in page discovery should follow the same case-insensitive behavior as the existing file-extension allowlist.
- The UI should not remove the existing raw uploaded/skipped file details; those remain useful for debugging.
- The feature should not require storing a permanent manifest of site pages in the database for the first version.
- The feature should not scan S3 after upload for the first version; upload session records are the source of truth for the just-completed upload.
- The dashboard should not persist or show the discovered page list after the upload modal result is gone.
- The page list should handle empty results gracefully, for example when a user uploads only assets and no HTML files.
- The UI should be written in Russian to match the current dashboard copy.

## Testing Decisions

- Good tests should verify externally visible behavior: given uploaded file records and a site domain, the system returns or renders the expected page links and home-page warning state.
- Backend page discovery rules should be covered with focused table tests because filename-to-page conversion is the highest-risk logic.
- Test cases should include root `index.html`, missing root `index.html`, nested `index.html`, ordinary `.html` files, uppercase `.HTML` if supported, skipped HTML files, non-HTML assets, and mixed upload results.
- Test the backend service method that builds the upload summary and keep handler tests limited to response shape and errors.
- Frontend testing should focus on manual QA because this webapp does not use unit tests.
- Existing upload session close behavior should remain covered: closing a session returns the upload result, checks site access, and closes the session.
- Manual QA should upload a valid ZIP with root `index.html`, a ZIP without root `index.html`, a ZIP with nested pages, and individual HTML files.

## Out of Scope

- Building a persistent file browser for all files currently stored on a site.
- Showing available pages outside the post-upload result UI or remembering them after the modal is closed.
- Crawling uploaded HTML to discover linked pages.
- Validating whether HTML content is well formed.
- Detecting broken links, missing assets, or JavaScript runtime errors in uploaded sites.
- Generating screenshots or previews of uploaded pages.
- Changing public routing behavior.
- Supporting non-HTML page formats.
- Adding analytics for page visits.
- Changing upload size limits or allowed file extensions.

## Further Notes

- Current backend behavior already stores uploaded and skipped files in the upload session and returns them from close upload session.
- Current public file serving resolves `/` to `index.html`. Paths with known extensions resolve exactly. Extensionless paths resolve using exact path, `.html`, and nested `index.html` lookup in that order.
- Current ZIP handling strips the first path segment from every archive entry. The missing-home-page warning should be based on the final stored names returned by the upload session, not the original ZIP entry names.
- The backend summary contract is intentionally preferred over frontend-only logic so page discovery rules stay close to upload/session behavior and can be tested in Go.
