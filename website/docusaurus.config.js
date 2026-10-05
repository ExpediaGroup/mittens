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
          exclude: [
            '**/_*.{js,jsx,ts,tsx,md,mdx}',
            '**/_*/**',
            '**/*.test.{js,jsx,ts,tsx}',
            '**/__tests__/**',
            'superpowers/**',
          ],
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      },
    ],
  ],

  themes: [
    [
      '@easyops-cn/docusaurus-search-local',
      {
        hashed: true,
        language: ['en'],
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
