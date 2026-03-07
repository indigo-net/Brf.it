import React from 'react';
import Link from '@docusaurus/Link';
import Translate, {translate} from '@docusaurus/Translate';

export default function Hero(): JSX.Element {
  return (
    <header className="hero hero--primary">
      <div className="container">
        <h1 className="hero__title">
          <Translate id="hero.title">Brf.it</Translate>
        </h1>
        <p className="hero__subtitle">
          <Translate id="hero.tagline">
            Give AI the gist, not the bloat. Extract function signatures from any codebase.
            Same context. 80% fewer tokens.
          </Translate>
        </p>
        <div className="cta-buttons">
          <Link
            to="/docs/"
            className="cta-button primary">
            <Translate id="hero.getStarted">Get Started</Translate>
          </Link>
          <Link
            to="https://github.com/indigo-net/Brf.it"
            className="cta-button secondary">
            <Translate id="hero.github">View on GitHub</Translate>
          </Link>
        </div>
      </div>
    </header>
  );
}
