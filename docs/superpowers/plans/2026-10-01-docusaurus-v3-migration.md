# Docusaurus v1 → v3 Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate `website/` from Docusaurus 1.14.7 to the latest v3, via an intermediate v2 step, with local-only changes (no branch/commit/PR created as part of this plan).

**Architecture:** Two sequential stages. Stage A uses the official `@docusaurus/migrate` CLI to scaffold a v2 site from the current v1 config/docs, then hand-cleans the output (config, homepage, footer, search). Stage B upgrades that v2 site's dependencies to v3 and works through the documented breaking changes. Each stage ends in a verified, buildable site.

**Tech Stack:** Docusaurus (v1.14.7 → v2.x → v3.x), React, Node.js (≥18 required by v3; CI bumped to Node 20), `@easyops-cn/docusaurus-search-local` for search, GitHub Actions + `peaceiris/actions-gh-pages` for deploy (unchanged).

**Spec:** [docs/superpowers/specs/2026-10-01-docusaurus-v3-migration-design.md](../specs/2026-10-01-docusaurus-v3-migration-design.md)

## Global Constraints

- No git branch, commit, or PR is created by this plan — all work stays as local uncommitted changes until the user explicitly asks for that step.
- Algolia DocSearch is dropped entirely, not carried forward — replaced by `@easyops-cn/docusaurus-search-local`.
- CI Node version bumps from 16 → 20, applied in Stage B (the point at which it becomes load-bearing, since v3 requires Node ≥18).
- No versioned docs, blog, or additional locales are added — scope stays at today's single-locale (`en`), unversioned, blog-less site.
- `docs/` content (the 5 markdown files) is not rewritten — only moved/relinked as the migration tooling requires.

---

## Stage A: v1 → v2

### Task 1: Scaffold v2 site via migrate CLI

**Files:**
- Create (scratch, outside `website/`): `/tmp/mittens-docusaurus-v2-migration/` — CLI output lands here first so we can inspect before swapping into place
- Read: `website/siteConfig.js`, `website/sidebars.json`, `website/pages/en/*.js`, `website/core/Footer.js`, `website/i18n/en.json`, `website/static/**`

**Interfaces:**
- Produces: a scratch v2 site directory whose `docusaurus.config.js`, `sidebars.js`, and copied `docs/` content become the input to Task 2.

- [ ] **Step 1: Check the migrate CLI's exact invocation**

The CLI's subcommand name has varied across releases — confirm before relying on memory.

Run: `npx @docusaurus/migrate@latest --help`

Expected: output lists a subcommand for v1→v2 migration (commonly `migrate-v1-to-v2` or `v1-to-v2`, taking `[siteDir] [newDir]` or similar positional args). Note the exact subcommand and argument order shown — use it verbatim in Step 2.

- [ ] **Step 2: Run the migration against a scratch output directory**

From the repo root:

```bash
mkdir -p /tmp/mittens-docusaurus-v2-migration
npx @docusaurus/migrate@latest <subcommand-from-step-1> website /tmp/mittens-docusaurus-v2-migration
```

(Substitute the real subcommand name found in Step 1 — do not guess; if the CLI prompts interactively, answer with the site's existing values: name `Mittens`, organization `ExpediaGroup`.)

Expected: command exits 0 and `/tmp/mittens-docusaurus-v2-migration` now contains a `docusaurus.config.js`, `sidebars.js`, `docs/`, `src/`, `package.json`.

- [ ] **Step 3: Inspect what the CLI produced**

Run: `find /tmp/mittens-docusaurus-v2-migration -maxdepth 2`

Confirm:
- `docs/` contains the 5 migrated markdown files (same `id`/`title` front matter as `../../docs/about/introduction.md` etc. — check at least one file's front matter survived)
- `sidebars.js` exists and references the same 3 categories (About, Installation, Contributing)
- `docusaurus.config.js` exists (content gets hand-cleaned in Task 2, so don't worry about its exact contents yet)

If any of the 5 docs files are missing or empty, stop and re-run Step 2 after checking the CLI's docs-path argument (it may need an explicit `--docs` flag pointing at the repo-root `docs/` directory rather than `website/docs`).

- [ ] **Step 4: No commit this task** — this task only produces scratch output; nothing in the repo has changed yet, so there is nothing to commit.

---

### Task 2: Hand-clean `docusaurus.config.js` and swap into `website/`

**Files:**
- Create: `website/docusaurus.config.js`
- Create: `website/sidebars.js`
- Delete: `website/siteConfig.js`, `website/sidebars.json`, `website/core/Footer.js`, `website/pages/en/help.js`, `website/pages/en/users.js` (content folded into config/homepage — see Task 3)
- Modify: `website/package.json` (replace v1 `docusaurus` dependency + scripts with v2 deps/scripts)

**Interfaces:**
- Consumes: scratch output from Task 1 (`/tmp/mittens-docusaurus-v2-migration/docusaurus.config.js`, `sidebars.js`) as a starting point to hand-edit
- Produces: `website/docusaurus.config.js` exporting the config object Task 3 (homepage) and Task 4 (search plugin) both add to

- [ ] **Step 1: Copy scratch config as starting point**

```bash
cp /tmp/mittens-docusaurus-v2-migration/docusaurus.config.js website/docusaurus.config.js
cp /tmp/mittens-docusaurus-v2-migration/sidebars.js website/sidebars.js
```

- [ ] **Step 2: Hand-edit `website/docusaurus.config.js` to match the original site's values**

Replace its contents with:

```js
// @ts-check

const config = {
  title: 'Mittens',
  tagline: 'Warm-up routine for http applications',
  favicon: 'img/mittens_intuit.svg',

  url: 'https://expediagroup.github.io',
  baseUrl: '/mittens/',

  organizationName: 'ExpediaGroup',
  projectName: 'Mittens',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          path: '../docs',
          routeBasePath: 'docs',
          sidebarPath: './sidebars.js',
          editUrl: 'https://github.com/ExpediaGroup/mittens/edit/main/website/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      },
    ],
  ],

  themeConfig: {
    navbar: {
      title: 'Mittens',
      logo: {
        alt: 'Mittens Logo',
        src: 'img/mittens_intuit.svg',
      },
      items: [
        {to: '/docs/about/introduction', label: 'Docs', position: 'left'},
        {
          href: 'https://github.com/ExpediaGroup/mittens',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [{label: 'Get Started', to: '/docs/about/introduction'}],
        },
        {
          title: 'More',
          items: [
            {label: 'GitHub', href: 'https://github.com/ExpediaGroup/mittens'},
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Expedia, Inc.`,
    },
    prism: {
      theme: require('prism-react-renderer/themes/github'),
    },
  },
};

module.exports = config;
```

Note: `docs.path: '../docs'` keeps content at the existing repo-root `docs/` directory rather than duplicating it under `website/docs/` — confirm in Task 5's build step that this resolves correctly relative to `website/`.

- [ ] **Step 3: Confirm `website/sidebars.js` matches the original structure**

Open `website/sidebars.js` and verify it defines the same 3 categories as `website/sidebars.json` did:

```js
const sidebars = {
  docs: [
    {
      type: 'category',
      label: 'About',
      items: ['about/introduction', 'about/getting-started'],
    },
    {
      type: 'category',
      label: 'Installation',
      items: ['installation/building', 'installation/running'],
    },
    {
      type: 'category',
      label: 'Contributing',
      items: ['contributing/contributing'],
    },
  ],
};

module.exports = sidebars;
```

If the CLI output differs in structure (e.g. flat array instead of categories), replace it with the above verbatim — this is the known-correct structure matching the original `sidebars.json`.

- [ ] **Step 4: Remove v1-only files**

```bash
rm -f website/siteConfig.js website/sidebars.json website/core/Footer.js
rm -f website/pages/en/help.js website/pages/en/users.js
```

(`index.js` is rebuilt in Task 3, not removed here.)

- [ ] **Step 5: Replace `website/package.json`**

```json
{
  "name": "mittensn-docs",
  "license": "Apache-2.0",
  "scripts": {
    "docusaurus": "docusaurus",
    "start": "docusaurus start",
    "build": "docusaurus build",
    "swizzle": "docusaurus swizzle",
    "deploy": "docusaurus deploy",
    "clear": "docusaurus clear",
    "serve": "docusaurus serve",
    "write-translations": "docusaurus write-translations",
    "write-heading-ids": "docusaurus write-heading-ids"
  },
  "dependencies": {
    "@docusaurus/core": "^2.4.3",
    "@docusaurus/preset-classic": "^2.4.3",
    "@mdx-js/react": "^1.6.22",
    "clsx": "^2.0.0",
    "prism-react-renderer": "^1.3.5",
    "react": "^17.0.2",
    "react-dom": "^17.0.2"
  },
  "browserslist": {
    "production": [">0.5%", "not dead", "not op_mini all"],
    "development": ["last 1 chrome version", "last 1 firefox version", "last 1 safari version"]
  },
  "engines": {
    "node": ">=16.14"
  }
}
```

- [ ] **Step 6: No build verification yet** — the homepage (Task 3) and search plugin (Task 4) are still missing; building now would fail on a missing `src/pages/index.js`. Verification happens at the end of Task 4.

---

### Task 3: Rebuild homepage and footer content

**Files:**
- Create: `website/src/pages/index.js`
- Create: `website/src/pages/index.module.css`
- Create: `website/src/css/custom.css`

**Interfaces:**
- Consumes: `website/docusaurus.config.js`'s exported `themeConfig`/site metadata (via `useDocusaurusContext()`), same shape as Task 2 Step 2's config
- Produces: homepage route `/` rendering the original hero + users showcase content

- [ ] **Step 1: Create `website/src/css/custom.css`**

```css
:root {
  --ifm-color-primary: #000099;
  --ifm-color-primary-dark: #000089;
  --ifm-color-primary-darker: #000080;
  --ifm-color-primary-darkest: #000069;
  --ifm-color-primary-light: #0000a9;
  --ifm-color-primary-lighter: #0000b2;
  --ifm-color-primary-lightest: #0000c9;
}
```

(Carries forward the original `siteConfig.colors.primaryColor: '#000099'`.)

- [ ] **Step 2: Create `website/src/pages/index.module.css`**

```css
.heroBanner {
  padding: 4rem 0;
  text-align: center;
  position: relative;
  overflow: hidden;
}

.buttons {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}

.usersSection {
  padding: 2rem 0;
  text-align: center;
}

.usersRow {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 2rem;
  flex-wrap: wrap;
  margin-top: 1rem;
}

.userLogo {
  max-height: 60px;
}
```

- [ ] **Step 3: Create `website/src/pages/index.js`**

```jsx
import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import Layout from '@theme/Layout';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import styles from './index.module.css';

const users = [
  {
    caption: 'Hotels.com',
    image: 'img/hotels_logo.svg',
    infoLink: 'https://www.hotels.com',
  },
  {
    caption: 'Expedia Group',
    image: 'img/expedia_group_logo.png',
    infoLink: 'https://www.expediagroup.com/',
  },
];

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  const logoUrl = useBaseUrl('img/mittens_intuit.svg');
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <img src={logoUrl} alt="Mittens Logo" width={120} />
        <h1 className="hero__title">{siteConfig.title}</h1>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link className="button button--secondary button--lg" to="/docs/about/introduction">
            Get Started
          </Link>
          <Link className="button button--secondary button--lg" to="https://github.com/ExpediaGroup/mittens">
            GitHub
          </Link>
        </div>
      </div>
    </header>
  );
}

function UsersSection() {
  return (
    <section className={styles.usersSection}>
      <div className="container">
        <h2>Who's using Mittens?</h2>
        <div className={styles.usersRow}>
          {users.map((user) => (
            <a href={user.infoLink} key={user.caption} target="_blank" rel="noreferrer">
              <img
                className={styles.userLogo}
                src={useBaseUrl(user.image)}
                alt={user.caption}
              />
            </a>
          ))}
        </div>
      </div>
    </section>
  );
}

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout title={siteConfig.title} description={siteConfig.tagline}>
      <HomepageHeader />
      <main>
        <UsersSection />
      </main>
    </Layout>
  );
}
```

Note: `useBaseUrl` is called inside the `.map()` callback above, which violates the rules-of-hooks lint rule (hooks can't run in a loop/callback conditionally). Since `useBaseUrl` here is a deterministic pure function of its argument and `users` is a static module-level array (not state), this won't cause a runtime bug, but if `npm run build`'s lint step (Task 5) flags it, replace the `useBaseUrl(user.image)` call with a plain string concatenation: `` `${useBaseUrl('')}${user.image}` `` computed once outside the map via a helper, e.g.:

```jsx
function UsersSection() {
  const baseUrl = useBaseUrl('/');
  return (
    <section className={styles.usersSection}>
      <div className="container">
        <h2>Who's using Mittens?</h2>
        <div className={styles.usersRow}>
          {users.map((user) => (
            <a href={user.infoLink} key={user.caption} target="_blank" rel="noreferrer">
              <img className={styles.userLogo} src={`${baseUrl}${user.image}`} alt={user.caption} />
            </a>
          ))}
        </div>
      </div>
    </section>
  );
}
```

Use this second version — it calls the hook once at the top level, which is correct.

- [ ] **Step 4: No commit this task** — build verification (next task) must pass first; commits happen in Task 5 Step 4 for the whole Stage A swap.

### Task 3 addendum (found post-hoc, after Task 6): homepage was missing two sections and a page

The brainstorming step that produced this plan only read the first ~80 lines of the original `website/pages/en/index.js` (via `head -60`/`head -80` in exploration) and missed the rest of the file. Comparing the migrated site against the live production site (https://opensource.expediagroup.com/mittens/) after Task 6 was complete surfaced the gap: the original homepage has two more sections after the hero — a "Run it as:" section (3 deployment options: standalone cmd tool, linked Docker container, Kubernetes sidecar, each with an icon from `static/img/{cmd.svg,docker.png,kubernetes.png}`) and a "Features" section (7-item bullet list) — neither of which Task 3 recreated. Additionally, the original site has a separate `/users` page (`website/pages/en/users.js`, deleted by Task 2 with no v2/v3 replacement) listing the same user logos plus an "Add your company" edit link, linked from the homepage's showcase section and the site footer.

- [ ] **Addendum Step 1: Extract shared user data**

Create `website/src/data/users.js`:

```js
export const users = [
  {
    caption: 'Hotels.com',
    image: 'img/hotels_logo.svg',
    infoLink: 'https://www.hotels.com',
  },
  {
    caption: 'Expedia Group',
    image: 'img/expedia_group_logo.png',
    infoLink: 'https://www.expediagroup.com/',
  },
];
```

Update `website/src/pages/index.js` to `import {users} from '../data/users';` instead of defining the array locally.

- [ ] **Addendum Step 2: Add "Run it as:" and "Features" sections to `website/src/pages/index.js`**

Add these two components (between `HomepageHeader` and `UsersSection` in the file, and render them in that order inside `<main>`):

```jsx
const runAsOptions = [
  {
    title: 'a standalone cmd tool',
    image: 'img/cmd.svg',
    content: 'Mittens can run as a standalone cmd tool to send requests over REST/gRPC.',
  },
  {
    title: 'a linked Docker Container',
    image: 'img/docker.png',
    content: 'You can also run it as a Docker container linked to your main app.',
  },
  {
    title: 'a Sidecar on Kubernetes',
    image: 'img/kubernetes.png',
    content: "Or even as a sidecar on Kubernetes, in the pod where your app's running.",
  },
];

function RunAsSection() {
  const baseUrl = useBaseUrl('/');
  return (
    <section className={styles.runAsSection}>
      <div className="container">
        <h2>Run it as:</h2>
        <div className={styles.runAsRow}>
          {runAsOptions.map((option) => (
            <div className={styles.runAsCard} key={option.title}>
              <img className={styles.runAsIcon} src={`${baseUrl}${option.image}`} alt={option.title} />
              <h3>{option.title}</h3>
              <p>{option.content}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

const featureList = [
  'Sends requests continuously for X seconds',
  'Supports REST and gRPC',
  'Supports HTTP/1.1, HTTP/2 and HTTP/2 Cleartext protocol for REST',
  'Supports custom HTTP headers',
  'Supports concurrent requests',
  'Supports placeholders for random elements in requests',
  'Provides files or/and endpoints that can be used as liveness/readiness probes in Kubernetes',
];

function FeaturesSection() {
  return (
    <section className={styles.featuresSection}>
      <div className="container">
        <h2>Features</h2>
        <ul className={styles.featuresList}>
          {featureList.map((feature) => (
            <li key={feature}>{feature}</li>
          ))}
        </ul>
      </div>
    </section>
  );
}
```

Update the `UsersSection`'s heading area to add a "More Mittens Users" link to the new `/users` page, right after the `usersRow` div:

```jsx
        <div className={styles.usersRow}>
          {users.map((user) => (
            <a href={user.infoLink} key={user.caption} target="_blank" rel="noreferrer">
              <img className={styles.userLogo} src={`${baseUrl}${user.image}`} alt={user.caption} />
            </a>
          ))}
        </div>
        <div className={styles.moreUsers}>
          <Link className="button button--secondary" to="/users">
            More Mittens Users
          </Link>
        </div>
```

Update the default-exported `Home` component to render all three sections in order:

```jsx
export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout title={siteConfig.title} description={siteConfig.tagline}>
      <HomepageHeader />
      <main>
        <RunAsSection />
        <FeaturesSection />
        <UsersSection />
      </main>
    </Layout>
  );
}
```

- [ ] **Addendum Step 3: Add CSS for the new sections to `website/src/pages/index.module.css`**

```css
.runAsSection {
  padding: 2rem 0;
  text-align: center;
}

.runAsRow {
  display: flex;
  justify-content: center;
  gap: 2rem;
  flex-wrap: wrap;
  margin-top: 1rem;
}

.runAsCard {
  max-width: 280px;
}

.runAsIcon {
  max-height: 80px;
  margin-bottom: 0.5rem;
}

.featuresSection {
  padding: 2rem 0;
  text-align: center;
}

.featuresList {
  display: inline-block;
  text-align: left;
  margin-top: 1rem;
}

.moreUsers {
  margin-top: 1rem;
}
```

- [ ] **Addendum Step 4: Recreate the `/users` page at `website/src/pages/users.js`**

```jsx
import React from 'react';
import Layout from '@theme/Layout';
import useBaseUrl from '@docusaurus/useBaseUrl';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import {users} from '../data/users';
import styles from './index.module.css';

export default function Users() {
  const {siteConfig} = useDocusaurusContext();
  const baseUrl = useBaseUrl('/');
  return (
    <Layout title={`${siteConfig.title} Users`} description="Who is using Mittens">
      <main className="container margin-vert--lg">
        <h1>Who is Using This?</h1>
        <p>This project is used by:</p>
        <div className={styles.usersRow}>
          {users.map((user) => (
            <a href={user.infoLink} key={user.caption} target="_blank" rel="noreferrer">
              <img className={styles.userLogo} src={`${baseUrl}${user.image}`} alt={user.caption} />
            </a>
          ))}
        </div>
        <p>Are you using this project?</p>
        <a
          className="button button--secondary"
          href="https://github.com/ExpediaGroup/mittens/edit/main/website/src/data/users.js"
        >
          Add your company
        </a>
      </main>
    </Layout>
  );
}
```

- [ ] **Addendum Step 5: Add a footer link to the new `/users` page**

In `website/docusaurus.config.js`'s `themeConfig.footer.links`, add a "User Showcase" entry alongside the existing "Get Started" link, e.g. under a `Community` column:

```js
        {
          title: 'Community',
          items: [{label: 'User Showcase', to: '/users'}],
        },
```

- [ ] **Addendum Step 6: Verify and commit**

`npm run build` (from `website/`) must still exit 0. `npm run serve` and re-check: homepage now shows hero → "Run it as:" (3 cards) → "Features" (7-item list) → "Who's using Mittens?" (2 logos + "More Mittens Users" link) → footer (with new "User Showcase" link); `/users` page renders the same 2 logos plus "Add your company". Commit this addendum's changes.

### Footer addendum (found post-hoc, after the homepage addendum): swizzle Footer for full parity

Comparing further against the live site and the original `website/core/Footer.js` (`git show main:website/core/Footer.js`) found the footer is still missing, beyond what config alone can express: a footer logo (mittens_intuit.svg) linking home, a GitHub "Star" button widget (stargazer count via the `buttons.github.io` script), and an "Expedia Group Open Source" logo+link below the link columns. Standard Docusaurus `themeConfig.footer` config can't render a third-party button widget or an extra logo outside the column grid — this requires ejecting (swizzling) the classic theme's `Footer` component.

- [ ] **Footer addendum Step 1: Add the GitHub buttons script**

In `website/docusaurus.config.js`, add a top-level `scripts` array (sibling to `presets`/`themeConfig`):

```js
  scripts: ['https://buttons.github.io/buttons.js'],
```

- [ ] **Footer addendum Step 2: Eject the Footer theme component**

Create `website/src/theme/Footer/index.js`:

```jsx
import React from 'react';
import Link from '@docusaurus/Link';
import useBaseUrl from '@docusaurus/useBaseUrl';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import {useThemeConfig} from '@docusaurus/theme-common';

export default function Footer() {
  const {siteConfig} = useDocusaurusContext();
  const {footer} = useThemeConfig();
  const footerIconUrl = useBaseUrl('img/mittens_intuit.svg');
  const ossLogoUrl = useBaseUrl('img/expedia_oss_logo.png');

  return (
    <footer className={`footer footer--${footer?.style ?? 'dark'}`}>
      <div className="container container-fluid">
        <div className="footer__links">
          <div className="footer__col">
            <Link to="/" className="footer__link-item">
              <img src={footerIconUrl} alt={siteConfig.title} width="66" height="58" />
            </Link>
          </div>
          <div className="footer__col">
            <div className="footer__title">Docs</div>
            <ul className="footer__items">
              <li className="footer__item">
                <Link className="footer__link-item" to="/docs/about/introduction">
                  Introduction
                </Link>
              </li>
            </ul>
          </div>
          <div className="footer__col">
            <div className="footer__title">Community</div>
            <ul className="footer__items">
              <li className="footer__item">
                <Link className="footer__link-item" to="/users">
                  User Showcase
                </Link>
              </li>
            </ul>
          </div>
          <div className="footer__col">
            <div className="footer__title">More</div>
            <ul className="footer__items">
              <li className="footer__item">
                <Link className="footer__link-item" href="https://github.com/ExpediaGroup/mittens">
                  GitHub
                </Link>
              </li>
              <li className="footer__item">
                <a
                  className="github-button"
                  href="https://github.com/ExpediaGroup/mittens"
                  data-icon="octicon-star"
                  data-count-href="/ExpediaGroup/mittens/stargazers"
                  data-show-count="true"
                  data-count-aria-label="# stargazers on GitHub"
                  aria-label="Star this project on GitHub">
                  Star
                </a>
              </li>
            </ul>
          </div>
        </div>
        <a
          href="https://expediagroup.github.io"
          target="_blank"
          rel="noreferrer noopener"
          className="egOpenSource">
          <img src={ossLogoUrl} alt="Expedia Group Open Source" width="170" height="45" />
        </a>
        {footer?.copyright && (
          <div className="footer__copyright">{footer.copyright}</div>
        )}
      </div>
    </footer>
  );
}
```

This reuses `themeConfig.footer.style`/`footer.copyright` (already set in `docusaurus.config.js` from Task 2) via `useThemeConfig()`, so the copyright line stays config-driven — only the link columns and extra logos are now hardcoded in this component (matching the original v1 `Footer.js`, which was also hardcoded, not data-driven).

- [ ] **Footer addendum Step 3: Remove the now-unused `themeConfig.footer.links` array**

In `website/docusaurus.config.js`, the `footer.links` array (Docs/Community/More columns) is no longer read by anything — the ejected `Footer` component above hardcodes its own columns instead. Remove the `links: [...]` key from `themeConfig.footer`, leaving just `style` and `copyright`.

- [ ] **Footer addendum Step 4: Verify and commit**

`npm run build` must still exit 0. `npm run serve` and check the footer HTML: a footer logo image linking to `/`, three columns (Docs→Introduction, Community→User Showcase, More→GitHub + a `github-button` anchor), an Expedia Group Open Source logo/link below the columns, and the copyright line. Commit this addendum's changes.

---

### Task 4: Add local search plugin and verify the v2 build

**Files:**
- Modify: `website/docusaurus.config.js` (add `themes: ['@easyops-cn/docusaurus-search-local']`)
- Modify: `website/package.json` (add dependency)

**Interfaces:**
- Consumes: `website/docusaurus.config.js` from Task 2
- Produces: a working, buildable v2 site — the deliverable this whole stage verifies

- [ ] **Step 1: Add the search plugin dependency**

Edit `website/package.json`'s `dependencies` block, adding:

```json
"@easyops-cn/docusaurus-search-local": "^0.40.1",
```

- [ ] **Step 2: Register the theme in `docusaurus.config.js`**

Add a top-level `themes` key to the config object from Task 2 Step 2 (alongside `presets`, `themeConfig`, etc.):

```js
  themes: [
    [
      '@easyops-cn/docusaurus-search-local',
      {
        hashed: true,
        language: ['en'],
      },
    ],
  ],
```

- [ ] **Step 3: Copy static assets across**

```bash
mkdir -p website/static
cp -r /tmp/mittens-docusaurus-v2-migration/static/* website/static/ 2>/dev/null || true
# If the CLI didn't carry static/img forward, the originals are still present on disk — confirm:
ls website/static/img/
```

Expected: `website/static/img/` still contains `mittens_intuit.svg`, `mittens_logo.svg`, `hotels_logo.svg`, `expedia_group_logo.png` (these were untouched by Task 2's `rm` commands, which only targeted `core/`, `pages/en/help.js`, `pages/en/users.js`, and config files — not `static/`).

- [ ] **Step 4: Install dependencies and build**

```bash
cd website
rm -rf node_modules package-lock.json yarn.lock
npm install
npm run build
```

Expected: exits 0, producing `website/build/`. If it fails, read the error — common causes at this stage:
- `docs.path: '../docs'` resolving incorrectly → confirm `website/../docs` (i.e. repo-root `docs/`) exists and run `ls ../docs` from inside `website/` to double check the relative path is correct
- missing `sidebars.js` category IDs not matching actual doc `id` front matter → cross-check against Task 2 Step 3's sidebar content vs. the 5 docs files' `id:` front matter values (`introduction`, `getting-started`, `building`, `running`, `contributing`)

- [ ] **Step 5: Serve and manually verify**

```bash
npm run serve
```

Open the printed local URL (typically `http://localhost:3000/mittens/`) and confirm:
- Homepage shows the Mittens logo, title "Mittens", tagline "Warm-up routine for http applications", "Get Started" and "GitHub" buttons, and a "Who's using Mittens?" section with Hotels.com and Expedia Group logos
- Clicking "Get Started" navigates to the Introduction doc page
- The left sidebar shows 3 categories: About, Installation, Contributing, each expandable to their pages
- A search box appears in the navbar; typing a term that appears in the docs (e.g. "warm-up") returns the Introduction page as a result

- [ ] **Step 6: Commit the whole Stage A swap**

```bash
cd /Users/sanvekar/Documents/repositories/mittens
git add website/
git status
```

Review `git status` output, confirm it shows the v1 files deleted and v2 files added/modified as expected, then stop — **do not run `git commit`**. Per the Global Constraints, this plan does not create commits; present the reviewed `git status` output to the user and wait for them to commit themselves or explicitly ask you to.

---

## Stage B: v2 → v3

### Task 5: Bump dependencies to v3 and resolve breaking changes

**Files:**
- Modify: `website/package.json`
- Modify: `website/docusaurus.config.js`

**Interfaces:**
- Consumes: the verified v2 site from Task 4
- Produces: `website/package.json`/`docusaurus.config.js` on v3, ready for Task 6's verification

- [ ] **Step 1: Bump core dependencies**

Edit `website/package.json`'s `dependencies`:

```json
"@docusaurus/core": "^3.6.0",
"@docusaurus/preset-classic": "^3.6.0",
"@mdx-js/react": "^3.0.0",
"@easyops-cn/docusaurus-search-local": "^0.44.0",
"clsx": "^2.0.0",
"prism-react-renderer": "^2.3.0",
"react": "^18.3.0",
"react-dom": "^18.3.0"
```

Update `engines` to:

```json
"engines": {
  "node": ">=18.0"
}
```

- [ ] **Step 2: Update `docusaurus.config.js` for v3's prism theme import shape**

Bumping `prism-react-renderer` from the v1.x line (used in Task 2, pinned `^1.3.5`, subpath-per-theme imports) to `^2.3.0` (Step 1 above) changes how themes are imported. Replace in `themeConfig.prism`:

```js
    prism: {
      theme: require('prism-react-renderer/themes/github'),
    },
```

with:

```js
    prism: {
      theme: require('prism-react-renderer').themes.github,
    },
```

(v1.x required the per-theme subpath import; v2.x+ exports a `themes` object instead. Confirm which shape the installed version actually exports by running `node -e "console.log(Object.keys(require('prism-react-renderer')))"` from `website/` after Task 6 Step 1's install — adjust to whichever shape matches if this differs.)

- [ ] **Step 3: No commit this task** — verification (Task 6) must pass first.

---

### Task 6: Verify the v3 build and bump CI Node version

**Files:**
- Modify: `.github/workflows/gh-pages.yml`

**Interfaces:**
- Consumes: the v3-bumped config/package.json from Task 5
- Produces: a working v3 site build, matching CI environment

- [ ] **Step 1: Clean install and build**

```bash
cd website
rm -rf node_modules package-lock.json yarn.lock build .docusaurus
npm install
npm run build
```

Expected: exits 0. If MDX-related errors appear (v3 uses MDX 3, stricter than v2's MDX 1), they'll point at a specific docs file and line — re-check that file's content against [the Context section's note](../specs/2026-10-01-docusaurus-v3-migration-design.md) that no unsafe `<...>`/`{...}` patterns were found; if one surfaces now, escape the offending character (e.g. wrap in backticks or escape `<` as `\<`, matching the pattern already used in `docs/about/getting-started.md`).

- [ ] **Step 2: Serve and manually re-verify**

```bash
npm run serve
```

Re-check the same 4 items from Task 4 Step 5 (homepage content, Get Started link, sidebar categories, search returning results).

- [ ] **Step 3: Bump CI Node version**

In `.github/workflows/gh-pages.yml`, change:

```yaml
      - name: Use Node.js 16
        uses: actions/setup-node@v4.0.2
        with:
          node-version: 16
```

to:

```yaml
      - name: Use Node.js 20
        uses: actions/setup-node@v4.0.2
        with:
          node-version: 20
```

- [ ] **Step 4: Review full diff, do not commit**

```bash
cd /Users/sanvekar/Documents/repositories/mittens
git add website/ .github/workflows/gh-pages.yml
git status
git diff --stat --cached
```

Present this to the user. Per the Global Constraints, stop here — no `git commit`, no branch, no PR unless the user explicitly asks for it.
