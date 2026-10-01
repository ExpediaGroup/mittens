# Docusaurus v1 → v3 Migration

## Context

`website/` runs Docusaurus **1.14.7** — the final v1 release (deprecated, no newer v1 version exists). Docs content lives in `docs/` at the repo root (5 markdown files across About/Installation/Contributing categories) and is referenced via `website/sidebars.json`. The site has a custom homepage (`pages/en/index.js`, `help.js`, `users.js`) built on v1's `CompLibrary`/`GridBlock` API, a custom `core/Footer.js`, Algolia DocSearch config in `siteConfig.js`, and is deployed to GitHub Pages via `.github/workflows/gh-pages.yml` (Node 16, `peaceiris/actions-gh-pages`).

No versioned docs, no blog, single locale (`en`).

## Goal

Migrate to the latest Docusaurus **v3**, via an intermediate **v2** step, so each major-version jump is independently buildable, reviewable, and revertable. Work happens locally only — no branch or PR created as part of this (confirmed with user: branch/PR creation requires separate explicit confirmation).

## Decisions (confirmed with user)

- **Search:** drop Algolia DocSearch entirely (the existing `apiKey`/`indexName` config has no `appId`, a sign it's on DocSearch's deprecated legacy shared infrastructure, and modern `@docsearch/react` requires a dedicated `appId` we don't have). Use **`@easyops-cn/docusaurus-search-local`** instead — builds a static local search index at build time, no external account/credentials, fine for a 5-page doc set. Added in the v2 step (installed as a plugin in `docusaurus.config.js`), carried through to v3 unchanged.
- **CI Node version:** bump `.github/workflows/gh-pages.yml` from Node 16 → **Node 20**, done alongside the v3 step (since v3 requires Node ≥18, this is when it becomes load-bearing).
- **Homepage/footer:** rebuild (not drop) — same hero (logo, tagline, "Get Started" + GitHub buttons), same Hotels.com/Expedia Group users showcase, same footer copyright/repo link — using v2's `Layout`/theme component API instead of the removed v1 `CompLibrary`. `pages/en/help.js` is dropped (it isn't linked from `headerLinks` today).
- **v1→v2 mechanism:** use the official `npx @docusaurus/migrate@latest v1 v2` CLI to scaffold + auto-convert config/docs, then hand-clean the output (homepage/footer rebuild, config review) rather than a fully from-scratch port.

## Plan

### Step 1: v1 → v2

1. Run the migrate CLI against the current `website/` + `docs/` into a scratch output dir, then replace `website/`'s contents with the migrated result (single swap, not an incremental diff) so the change is one bisectable unit.
2. Hand-clean `docusaurus.config.js`: title, tagline, url/baseUrl, organizationName/projectName, favicon, GitHub header link, `themeConfig.footer` (copyright + repo link, replacing custom `Footer.js`). Replace the Algolia block with the `@easyops-cn/docusaurus-search-local` plugin (added to `plugins: []`, no credentials needed).
3. Confirm `sidebars.js` (migrated from `sidebars.json`) matches the original 3-category structure (About / Installation / Contributing).
4. Rebuild `src/pages/index.js`: hero with logo/tagline/buttons, users showcase section (Hotels.com, Expedia Group logos/links). Drop `help.js`/`users.js` as separate page files — fold into `index.js` or theme config as v2 idioms dictate.
5. Copy `static/img/*` assets across unchanged.
6. Verify: `npm install && npm run build` succeeds; `npm run serve` — check homepage renders (hero, buttons, users section), docs nav (3 categories, 5 pages), and the local search box returns results for a known doc term.

### Step 2: v2 → v3

1. Bump `@docusaurus/core`, `@docusaurus/preset-classic` (and any other `@docusaurus/*` deps) to latest v3, plus `@easyops-cn/docusaurus-search-local` to its v3-compatible version.
2. Work through the official v3 migration guide's breaking changes: React 18 bump, MDX v2→v3 stricter parsing (re-check all 5 docs files render — already spot-checked, no unsafe `<...>`/`{...}` patterns found, but re-verify post-bump), config field renames, CSS/infima changes.
3. Bump `.github/workflows/gh-pages.yml` Node 16 → 20.
4. Verify: clean `npm install && npm run build` + `npm run serve`, same checks as step 1.

## Out of scope

- Algolia DocSearch — dropped (see Decisions). Not revisited as part of this migration; could be reconsidered later if a dedicated `appId` is obtained.
- Adding versioned docs, blog, or additional locales.
- Creating a git branch, committing, or opening a PR — stays as local uncommitted (or user-reviewed) changes until the user explicitly asks for that step.

## Testing

For each step: `npm install`, `npm run build` (must exit 0), `npm run serve` with manual check of homepage + docs nav + local search rendering/results. No automated test suite exists for the docs site; verification is build-success + manual render check.
