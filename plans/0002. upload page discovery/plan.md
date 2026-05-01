# Plan: Uploaded Page Discovery

> Source PRD: [Uploaded Page Discovery](./prd.md)

## Architectural decisions

Durable decisions that apply across all phases:

- **Routes**:
  The existing upload session close endpoint remains the API boundary for post-upload results.
  The dashboard continues using the existing upload flow through the SvelteKit API proxy.
  No new endpoint is needed for the first version.
- **Schema**:
  No database schema change is required.
  The existing upload session file records remain the source of truth for the just-completed upload.
  No permanent page manifest is stored.
- **Key models**:
  Introduce an upload summary response model returned when an upload session closes.
  The summary contains uploaded files, skipped files, discovered pages, and a root home-page presence flag.
  A discovered page contains both display path and absolute URL so the frontend can render a regular link without reimplementing page-path rules.
- **Page discovery rules**:
  Only successfully uploaded files can become available pages.
  Only normalized `.html` filenames are considered pages.
  HTML detection follows the same case-insensitive extension behavior as existing file resolution.
  `index.html` is the home page and maps to the site root without a trailing slash.
  `folder/index.html` maps to `/folder` when that clean path resolves to it and does not satisfy the root home-page requirement.
  `page.html` maps to the clean extensionless path, for example `/about`.
  Page URL generation follows current serve-file precedence: extension paths resolve exactly, and extensionless paths try exact path, then `path.html`, then `path/index.html`.
  If two uploaded HTML files compete for the same clean URL, the file that current serving logic would resolve gets the clean URL and the other file gets an exact path URL.
  Skipped HTML files and non-HTML assets do not appear as available pages.
- **UI behavior**:
  Page discovery is shown only in the upload modal after that upload completes.
  Available pages render as ordinary links.
  Missing root `index.html` renders as a separate warning-colored alert card.
  Existing uploaded/skipped file details remain available.
- **Testing**:
  Backend automated tests carry the page-discovery rules and response contract.
  Frontend validation is manual because the webapp does not use unit tests.

---

## Phase 1: Return a backend upload summary

**User stories**: 4, 5, 6, 7, 9, 10, 11, 12, 13, 14, 18, 19, 20

### What to build

Change the upload session close flow so it returns a structured summary for the completed upload. The summary should preserve the existing raw uploaded and skipped file information while adding backend-derived available pages and a `hasHomePage`-style signal. This slice is complete when API clients can distinguish uploaded files, skipped files, page links, and missing-home-page state without doing filename-to-route conversion themselves.

### Acceptance criteria

- [ ] Closing an upload session returns uploaded files, skipped files, available pages, and home-page presence in one response.
- [ ] Available pages are derived only from successfully uploaded `.html` files.
- [ ] `index.html` maps to the site root and sets the home-page presence flag.
- [ ] `about.html` maps to `/about`.
- [ ] `about/index.html` maps to `/about` when no higher-precedence uploaded file owns that clean path and does not set the root home-page presence flag.
- [ ] When `about.html` and `about/index.html` are both uploaded, the generated links follow current serving precedence and do not produce two links that open the same file.
- [ ] Uppercase HTML extensions are handled consistently with existing file resolution.
- [ ] Non-HTML assets and skipped files do not appear in available pages.
- [ ] Existing upload success/error behavior and session authorization remain unchanged.
- [ ] Backend tests cover the page discovery table cases from the PRD.

---

## Phase 2: Render available pages after upload

**User stories**: 1, 2, 3, 7, 8, 9, 13, 14, 15, 16, 17, 19, 20

### What to build

Update the dashboard upload result to consume the new upload summary response and render a dedicated available-pages section after a successful upload. Users should see regular links for each discovered page, while the existing uploaded and skipped file details remain visible for debugging. This slice is complete when a successful upload with HTML files gives the user clear, clickable page URLs in the upload modal.

### Acceptance criteria

- [ ] The frontend API client understands the new upload summary response shape.
- [ ] A successful upload with `index.html` shows the site root as an available page link.
- [ ] A successful upload with `about.html` shows a regular link ending in `/about`.
- [ ] A successful upload with nested pages shows those pages in the available-pages section using paths supplied by the backend.
- [ ] Uploaded file details and skipped file details still render separately from available pages.
- [ ] Upload errors still render as errors and do not show a successful page discovery result.
- [ ] The page list disappears when the upload modal state is reset or closed.

---

## Phase 3: Add the missing-home-page warning

**User stories**: 4, 5, 6, 12, 15, 16, 17

### What to build

Add a warning visual state to the shared alert pattern and use it in the upload result when the backend summary says the root home page is missing. The warning should be separate from the success message, use warning color, and tell the user that `index.html` must be placed at the root of the uploaded site or ZIP archive for the main site URL to work.

### Acceptance criteria

- [ ] The shared alert UI supports a warning-colored card distinct from success, info, and error.
- [ ] A successful upload without root `index.html` shows a separate warning alert.
- [ ] A nested `index.html` such as `about/index.html` still triggers the missing-home-page warning when root `index.html` is absent.
- [ ] A successful upload with root `index.html` does not show the warning.
- [ ] The warning copy is in Russian and gives the user a concrete fix.
- [ ] The warning does not replace the upload success result or hide available nested pages.

---

## Phase 4: End-to-end validation and compatibility hardening

**User stories**: 7, 8, 13, 14, 15, 16, 17, 18, 19, 20

### What to build

Harden the full upload result path across backend, proxy, API client, and dashboard UI. Confirm existing upload scenarios still work while the new response shape is handled consistently. This phase is mostly about compatibility and release confidence: mixed uploads, ZIP extraction behavior, skipped files, no-HTML uploads, and upload failures should all produce predictable results.

### Acceptance criteria

- [ ] Backend tests cover mixed uploads with HTML, assets, skipped files, missing home page, and present home page.
- [ ] The API proxy passes the structured close-upload response through unchanged.
- [ ] Manual QA covers a ZIP with root `index.html`, a ZIP without root `index.html`, a ZIP with nested pages, individual HTML files, and asset-only uploads.
- [ ] Asset-only uploads show no available pages and still show the missing-home-page warning.
- [ ] Failed uploads do not produce a misleading available-pages or missing-home-page result.
- [ ] The dashboard copy and layout remain readable on desktop and mobile upload modal sizes.
