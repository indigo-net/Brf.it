import React from 'react';
import Translate from '@docusaurus/Translate';

interface Language {
  name: string;
  icon: string;
}

const languages: Language[] = [
  { name: 'Go', icon: '🐹' },
  { name: 'TypeScript', icon: '📘' },
  { name: 'JavaScript', icon: '📒' },
  { name: 'Python', icon: '🐍' },
  { name: 'Java', icon: '☕' },
  { name: 'Kotlin', icon: '🟣' },
  { name: 'Rust', icon: '🦀' },
  { name: 'Ruby', icon: '💎' },
  { name: 'PHP', icon: '🐘' },
  { name: 'Swift', icon: '🍎' },
  { name: 'Scala', icon: '🔴' },
  { name: 'C/C++', icon: '⚙️' },
];

export default function LanguageGrid(): JSX.Element {
  return (
    <section className="section" style={{ background: 'var(--ifm-color-emphasis-100)' }}>
      <div className="container">
        <h2 className="section-title">
          <Translate id="languages.title">Supported Languages</Translate>
        </h2>
        <p className="section-subtitle">
          <Translate id="languages.subtitle">
            Tree-sitter powered parsing for accurate signature extraction
          </Translate>
        </p>

        <div className="language-grid">
          {languages.map(lang => (
            <div key={lang.name} className="language-card">
              <span className="language-icon">{lang.icon}</span>
              <span className="language-name">{lang.name}</span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
