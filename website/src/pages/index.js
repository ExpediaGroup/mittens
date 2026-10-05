import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import Layout from '@theme/Layout';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import {users} from '../data/users';
import styles from './index.module.css';

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

const runAsOptions = [
  {
    title: 'a standalone cmd tool',
    image: 'img/cmd.svg',
    content: 'Mittens can run as a standalone cmd tool to send requests over REST/gRPC.',
  },
  {
    title: 'a linked Docker Container',
    image: 'img/docker.png',
    content: 'You can also run it as a Docker container linked to your main app.',
  },
  {
    title: 'a Sidecar on Kubernetes',
    image: 'img/kubernetes.png',
    content: "Or even as a sidecar on Kubernetes, in the pod where your app's running.",
  },
];

function RunAsSection() {
  const baseUrl = useBaseUrl('/');
  return (
    <section className={styles.runAsSection}>
      <div className="container">
        <h2>Run it as:</h2>
        <div className={styles.runAsRow}>
          {runAsOptions.map((option) => (
            <div className={styles.runAsCard} key={option.title}>
              <img className={styles.runAsIcon} src={`${baseUrl}${option.image}`} alt={option.title} />
              <h3>{option.title}</h3>
              <p>{option.content}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

const featureList = [
  'Sends requests continuously for X seconds',
  'Supports REST and gRPC',
  'Supports HTTP/1.1, HTTP/2 and HTTP/2 Cleartext protocol for REST',
  'Supports custom HTTP headers',
  'Supports concurrent requests',
  'Supports placeholders for random elements in requests',
  'Provides files or/and endpoints that can be used as liveness/readiness probes in Kubernetes',
];

function FeaturesSection() {
  return (
    <section className={styles.featuresSection}>
      <div className="container">
        <h2>Features</h2>
        <ul className={styles.featuresList}>
          {featureList.map((feature) => (
            <li key={feature}>{feature}</li>
          ))}
        </ul>
      </div>
    </section>
  );
}

function UsersSection() {
  const baseUrl = useBaseUrl('/');
  return (
    <section className={styles.usersSection}>
      <div className="container">
        <h2>Who is Using This?</h2>
        <p>This project is used by:</p>
        <div className={styles.usersRow}>
          {users.map((user) => (
            <a href={user.infoLink} key={user.caption} target="_blank" rel="noreferrer">
              <img className={styles.userLogo} src={`${baseUrl}${user.image}`} alt={user.caption} height={60} />
            </a>
          ))}
        </div>
        <div className={styles.moreUsers}>
          <Link className="button button--secondary" to="/users">
            More Mittens Users
          </Link>
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
        <RunAsSection />
        <FeaturesSection />
        <UsersSection />
      </main>
    </Layout>
  );
}
