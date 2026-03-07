import React from 'react';
import Translate from '@docusaurus/Translate';

interface Feature {
  title: string;
  description: string;
  icon: string;
}

const features: Feature[] = [
  {
    title: 'Tree-sitter AST',
    description: 'Accurate syntax parsing using Tree-sitter for precise signature extraction.',
    icon: '🌳',
  },
  {
    title: 'Multi-format Output',
    description: 'Export as XML or Markdown. Perfect for AI context windows.',
    icon: '📄',
  },
  {
    title: 'Token Counting',
    description: 'Automatic token estimation to help manage AI costs.',
    icon: '🔢',
  },
  {
    title: 'Import Detection',
    description: 'Optional import/export statement extraction for dependency context.',
    icon: '📥',
  },
  {
    title: '12 Languages',
    description: 'Go, TypeScript, Python, Java, Kotlin, Rust, Ruby, PHP, Swift, Scala, C/C++, JavaScript.',
    icon: '🌐',
  },
  {
    title: 'Zero Config',
    description: 'Works out of the box. Respects .gitignore automatically.',
    icon: '⚡',
  },
];

export default function FeaturesSection(): JSX.Element {
  return (
    <section className="section" style={{ background: 'var(--ifm-color-emphasis-100)' }}>
      <div className="container">
        <h2 className="section-title">
          <Translate id="features.title">Features</Translate>
        </h2>
        <p className="section-subtitle">
          <Translate id="features.subtitle">
            Everything you need to give AI the right context
          </Translate>
        </p>

        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
          gap: '1.5rem',
        }}>
          {features.map((feature, idx) => (
            <div
              key={idx}
              style={{
                background: 'var(--ifm-card-background-color)',
                padding: '1.5rem',
                borderRadius: '12px',
                boxShadow: 'var(--ifm-global-shadow-lw)',
              }}
            >
              <div style={{ fontSize: '2rem', marginBottom: '0.75rem' }}>{feature.icon}</div>
              <h3 style={{ fontSize: '1.125rem', fontWeight: 600, marginBottom: '0.5rem' }}>
                {feature.title}
              </h3>
              <p style={{ color: 'var(--ifm-color-emphasis-600)', fontSize: '0.875rem', margin: 0 }}>
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
