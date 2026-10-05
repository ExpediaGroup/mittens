import React from 'react';
import Link from '@docusaurus/Link';
import useBaseUrl from '@docusaurus/useBaseUrl';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import {useThemeConfig} from '@docusaurus/theme-common';

export default function Footer() {
  const {siteConfig} = useDocusaurusContext();
  const {footer} = useThemeConfig();
  const footerIconUrl = useBaseUrl('img/mittens_intuit.svg');
  const ossLogoUrl = useBaseUrl('img/expedia_oss_logo.png');

  return (
    <footer className={`footer footer--${footer?.style ?? 'dark'}`}>
      <div className="container container-fluid">
        <div className="footer__links">
          <div className="footer__col">
            <Link to="/" className="footer__link-item">
              <img src={footerIconUrl} alt={siteConfig.title} width="66" height="58" />
            </Link>
          </div>
          <div className="footer__col">
            <div className="footer__title">Docs</div>
            <ul className="footer__items">
              <li className="footer__item">
                <Link className="footer__link-item" to="/docs/about/introduction">
                  Introduction
                </Link>
              </li>
            </ul>
          </div>
          <div className="footer__col">
            <div className="footer__title">Community</div>
            <ul className="footer__items">
              <li className="footer__item">
                <Link className="footer__link-item" to="/users">
                  User Showcase
                </Link>
              </li>
            </ul>
          </div>
          <div className="footer__col">
            <div className="footer__title">More</div>
            <ul className="footer__items">
              <li className="footer__item">
                <Link className="footer__link-item" href="https://github.com/ExpediaGroup/mittens">
                  GitHub
                </Link>
              </li>
              <li className="footer__item">
                <a
                  className="github-button"
                  href="https://github.com/ExpediaGroup/mittens"
                  data-icon="octicon-star"
                  data-count-href="/ExpediaGroup/mittens/stargazers"
                  data-show-count="true"
                  data-count-aria-label="# stargazers on GitHub"
                  aria-label="Star this project on GitHub">
                  Star
                </a>
              </li>
            </ul>
          </div>
        </div>
        <a
          href="https://expediagroup.github.io"
          target="_blank"
          rel="noreferrer noopener"
          className="egOpenSource">
          <img src={ossLogoUrl} alt="Expedia Group Open Source" width="170" height="45" />
        </a>
        {footer?.copyright && (
          <div className="footer__copyright">{footer.copyright}</div>
        )}
      </div>
    </footer>
  );
}
