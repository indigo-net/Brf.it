import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docsSidebar: [
    {
      type: 'doc',
      id: 'getting-started',
      label: 'Getting Started',
    },
    {
      type: 'doc',
      id: 'cli-reference',
      label: 'CLI Reference',
    },
    {
      type: 'category',
      label: 'Languages',
      items: [
        'languages/index',
        'languages/go',
        'languages/typescript',
        'languages/python',
        'languages/java',
        'languages/kotlin',
        'languages/rust',
        'languages/ruby',
        'languages/php',
        'languages/swift',
        'languages/scala',
        'languages/c-cpp',
      ],
    },
  ],
};

export default sidebars;
