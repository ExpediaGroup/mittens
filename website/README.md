# Mittens Docs Site

This is the [Docusaurus](https://docusaurus.io/) v3 site that powers the Mittens documentation at https://opensource.expediagroup.com/mittens/.

## Install

```sh
npm install
```

## Run locally

```sh
npm start
```

This starts a local dev server and opens the site in your browser with hot reload.

## Build

```sh
npm run build
```

This generates static content into `website/build/`, which can be served by any static hosting service.

## Where content lives

- `docs/` (repo root, outside `website/`) — the actual documentation pages, as Markdown files with `id`/`title` front matter.
- `website/docusaurus.config.js` — site-wide configuration (title, navbar, footer, theme, plugins).
- `website/sidebars.js` — the docs sidebar/navigation structure, referencing doc IDs from `docs/`.
- `website/src/pages/` — standalone React pages, such as the homepage and the `/users` showcase page.
- `website/static/` — static assets (images, etc.) served as-is at the site root.
