import React from 'react';
import Layout from '@theme/Layout';
import useBaseUrl from '@docusaurus/useBaseUrl';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import {users} from '../data/users';
import styles from './index.module.css';

export default function Users() {
  const {siteConfig} = useDocusaurusContext();
  const baseUrl = useBaseUrl('/');
  return (
    <Layout title={`${siteConfig.title} Users`} description="Who is using Mittens">
      <main className="container margin-vert--lg">
        <h1>Who is Using This?</h1>
        <p>This project is used by:</p>
        <div className={styles.usersRow}>
          {users.map((user) => (
            <a href={user.infoLink} key={user.caption} target="_blank" rel="noreferrer">
              <img className={styles.userLogo} src={`${baseUrl}${user.image}`} alt={user.caption} height={60} />
            </a>
          ))}
        </div>
        <p>Are you using this project?</p>
        <a
          className="button button--secondary"
          href="https://github.com/ExpediaGroup/mittens/edit/main/website/src/data/users.js"
        >
          Add your company
        </a>
      </main>
    </Layout>
  );
}
