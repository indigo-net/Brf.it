import React, {useState} from 'react';
import Translate from '@docusaurus/Translate';

const installCommands = {
  macos: {
    label: 'macOS',
    icon: '🍎',
    commands: [
      { cmd: 'brew tap indigo-net/tap', desc: 'Add tap' },
      { cmd: 'brew install brfit', desc: 'Install' },
    ],
  },
  linux: {
    label: 'Linux',
    icon: '🐧',
    commands: [
      { cmd: 'curl -sSL https://github.com/indigo-net/Brf.it/releases/latest/download/brfit-linux-amd64 -o brfit', desc: 'Download' },
      { cmd: 'chmod +x brfit && sudo mv brfit /usr/local/bin/', desc: 'Install' },
    ],
  },
  windows: {
    label: 'Windows',
    icon: '🪟',
    commands: [
      { cmd: 'scoop bucket add indigo-net https://github.com/indigo-net/scoop-bucket', desc: 'Add bucket' },
      { cmd: 'scoop install brfit', desc: 'Install' },
    ],
  },
  npm: {
    label: 'npm',
    icon: '📦',
    commands: [
      { cmd: 'npm install -g @indigo-net/brfit', desc: 'Install globally' },
    ],
  },
};

type Platform = keyof typeof installCommands;

export default function InstallSection(): JSX.Element {
  const [platform, setPlatform] = useState<Platform>('macos');
  const current = installCommands[platform];

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  return (
    <section className="section">
      <div className="container">
        <h2 className="section-title">
          <Translate id="install.title">Quick Install</Translate>
        </h2>
        <p className="section-subtitle">
          <Translate id="install.subtitle">
            Get started in seconds with your preferred package manager
          </Translate>
        </p>

        {/* Platform Tabs */}
        <div style={{ display: 'flex', justifyContent: 'center', gap: '0.5rem', marginBottom: '2rem', flexWrap: 'wrap' }}>
          {(Object.entries(installCommands) as [Platform, typeof installCommands[Platform]][]).map(([key, value]) => (
            <button
              key={key}
              onClick={() => setPlatform(key)}
              style={{
                padding: '0.75rem 1.25rem',
                border: platform === key ? '2px solid var(--ifm-color-primary)' : '1px solid var(--ifm-color-emphasis-400)',
                borderRadius: '8px',
                background: platform === key ? 'var(--ifm-color-primary)' : 'transparent',
                color: platform === key ? 'white' : 'var(--ifm-color-emphasis-700)',
                cursor: 'pointer',
                fontWeight: 500,
                display: 'flex',
                alignItems: 'center',
                gap: '0.5rem',
              }}
            >
              <span>{value.icon}</span>
              <span>{value.label}</span>
            </button>
          ))}
        </div>

        {/* Commands */}
        <div style={{ maxWidth: '600px', margin: '0 auto' }}>
          {current.commands.map((item, idx) => (
            <div key={idx} className="install-command" style={{ marginBottom: '1rem' }}>
              <code style={{ flex: 1 }}>{item.cmd}</code>
              <button
                onClick={() => copyToClipboard(item.cmd)}
                className="copy-button"
                title="Copy to clipboard"
              >
                📋
              </button>
            </div>
          ))}
        </div>

        {/* Quick Example */}
        <div style={{ marginTop: '2rem', textAlign: 'center' }}>
          <p style={{ color: 'var(--ifm-color-emphasis-600)', marginBottom: '1rem' }}>
            <Translate id="install.example">Then run:</Translate>
          </p>
          <div className="install-command" style={{ maxWidth: '500px', margin: '0 auto' }}>
            <code style={{ flex: 1 }}>brfit . -f md --include-imports</code>
            <button
              onClick={() => copyToClipboard('brfit . -f md --include-imports')}
              className="copy-button"
              title="Copy to clipboard"
            >
              📋
            </button>
          </div>
        </div>
      </div>
    </section>
  );
}
