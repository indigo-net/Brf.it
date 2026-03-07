import type {ReactNode} from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Hero from '@site/src/components/Hero';
import TokenComparison from '@site/src/components/TokenComparison';
import FeaturesSection from '@site/src/components/FeaturesSection';
import LanguageGrid from '@site/src/components/LanguageGrid';
import InstallSection from '@site/src/components/InstallSection';

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title="Brf.it - Give AI the gist, not the bloat"
      description="Extract function signatures from any codebase. Same context. 80% fewer tokens. Perfect for AI coding assistants.">
      <Hero />
      <main>
        <TokenComparison />
        <FeaturesSection />
        <LanguageGrid />
        <InstallSection />
      </main>
    </Layout>
  );
}
