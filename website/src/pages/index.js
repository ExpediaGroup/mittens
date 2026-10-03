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
