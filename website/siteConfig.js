/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

// See https://docusaurus.io/docs/site-config for all the possible
// site configuration options.

// List of projects/orgs using your project for the users page.
const users = [
  {
    caption: 'Hotels.com',

    // You will need to prepend the image path with your baseUrl
    // if it is not '/', like: '/test-site/img/image.jpg'.
    image: 'img/hotels_logo.svg',
    infoLink: 'https://www.hotels.com',
    pinned: true,
  },
];

const siteConfig = {
  title: 'Mittens', // Title for your website.
  tagline: 'Warm-up routine for http applications',

  url: 'https://expediagroup.github.io', // Your website URL
  baseUrl: '/mittens/', // Base URL for your project */

  // Used for publishing and more
  projectName: 'Mittens',
  organizationName: 'ExpediaGroup',

  // For no header links in the top nav bar -> headerLinks: [],
  headerLinks: [
    { search: true },
    {doc: 'about/introduction', label: 'Docs'},
    {href: 'https://github.com/ExpediaGroup/mittens', label: 'GitHub'}
  ],

  users,

  /* path to images for header/footer */
  headerIcon: 'img/mittens_intuit.svg',
  footerIcon: 'img/mittens_intuit.svg',
  favicon: 'img/mittens_intuit.svg',

  /* Colors for website */
  colors: {
    primaryColor: '#000099',
    secondaryColor: '#01325A',
  },

  // This copyright info is used in /core/Footer.js and blog RSS/Atom feeds.
  copyright: `Copyright © ${new Date().getFullYear()} Expedia, Inc.`,

  highlight: {
    // Highlight.js theme to use for syntax highlighting in code blocks.
    theme: 'github',
  },

  algolia: {
    apiKey: 'b384c10816b317ab3062e38860f2e98b',
    indexName: 'expediagroup_mittens',
  },

  // Add custom scripts here that would be placed in <script> tags.
  scripts: ['https://buttons.github.io/buttons.js'],

  // On page navigation for the current documentation page.
  onPageNav: 'separate',
  // No .html extensions for paths.
  cleanUrl: true,

  // Open Graph and Twitter card images.
  ogImage: 'img/undraw_online.svg',
  twitterImage: 'img/undraw_tweetstorm.svg',

  // For sites with a sizable amount of content, set collapsible to true.
  // Expand/collapse the links and subcategories under categories.
  docsSideNavCollapsible: true,

  // Show documentation's last contributor's name.
  // enableUpdateBy: true,

  // Show documentation's last update time.
  // enableUpdateTime: true,

  // You may provide arbitrary config keys to be used as needed by your
  // template. For example, if you need your repo's URL...
  repoUrl: 'https://github.com/ExpediaGroup/mittens',
};

module.exports = siteConfig;
