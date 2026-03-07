import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Brf.it',
  tagline: 'Give AI the gist, not the bloat',
  favicon: 'img/favicon.ico',

  future: {
    v4: true,
  },

  url: 'https://indigo-net.github.io',
  baseUrl: '/Brf.it/',

  organizationName: 'indigo-net',
  projectName: 'Brf.it',
  trailingSlash: false,

  onBrokenLinks: 'throw',

  i18n: {
    defaultLocale: 'en',
    locales: ['en', 'ko', 'ja', 'de', 'hi'],
    localeConfigs: {
      en: { label: 'English', direction: 'ltr' },
      ko: { label: '한국어', direction: 'ltr' },
      ja: { label: '日本語', direction: 'ltr' },
      de: { label: 'Deutsch', direction: 'ltr' },
      hi: { label: 'हिन्दी', direction: 'ltr' },
    },
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/indigo-net/Brf.it/tree/main/website/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/brfit-social-card.png',
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'Brf.it',
      logo: {
        alt: 'Brf.it Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          type: 'localeDropdown',
          position: 'right',
          dropdownItemsBefore: [],
          dropdownItemsAfter: [],
        },
        {
          href: 'https://github.com/indigo-net/Brf.it',
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
          items: [
            {
              label: 'Getting Started',
              to: '/docs/',
            },
            {
              label: 'CLI Reference',
              to: '/docs/cli-reference',
            },
            {
              label: 'Supported Languages',
              to: '/docs/languages/',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'GitHub Issues',
              href: 'https://github.com/indigo-net/Brf.it/issues',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/indigo-net/Brf.it',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} indigo-net. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['go', 'typescript', 'python', 'java', 'kotlin', 'rust', 'ruby', 'php', 'swift', 'scala'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
