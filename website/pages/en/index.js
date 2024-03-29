/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

const React = require('react');

const CompLibrary = require('../../core/CompLibrary.js');

const Container = CompLibrary.Container;
const GridBlock = CompLibrary.GridBlock;

class HomeSplash extends React.Component {
    render() {
        const {siteConfig, language = ''} = this.props;
        const {baseUrl, docsUrl} = siteConfig;
        const docsPart = `${docsUrl ? `${docsUrl}/` : ''}`;
        const langPart = `${language ? `${language}/` : ''}`;
        const docUrl = doc => `${baseUrl}${docsPart}${langPart}${doc}`;

        const SplashContainer = props => (
            <div className="homeContainer">
                <div className="homeSplashFade">
                    <div className="wrapper homeWrapper">{props.children}</div>
                </div>
            </div>
        );

        const Logo = props => (
            <div className="projectLogo">
                <img src={props.img_src} alt="Project Logo"/>
            </div>
        );

        const ProjectTitle = () => (
            <h2 className="projectTitle">
                {siteConfig.title}
                <small>{siteConfig.tagline}</small>
            </h2>
        );

        const PromoSection = props => (
            <div className="section promoSection">
                <div className="promoRow">
                    <div className="pluginRowBlock">{props.children}</div>
                </div>
            </div>
        );

        const Button = props => (
            <div className="pluginWrapper buttonWrapper">
                <a className="button" href={props.href} target={props.target}>
                    {props.children}
                </a>
            </div>
        );

        return (
            <SplashContainer>
                <Logo img_src={`${baseUrl}img/mittens_intuit.svg`}/>
                <div className="inner">
                    <ProjectTitle siteConfig={siteConfig}/>
                    <PromoSection>
                        <Button href={docUrl('about/introduction.html')}>Get Started</Button>
                        <Button href="https://github.com/ExpediaGroup/mittens">GitHub</Button>
                    </PromoSection>
                </div>
            </SplashContainer>
        );
    }
}

class Index extends React.Component {
    render() {
        const {config: siteConfig, language = ''} = this.props;
        const {baseUrl} = siteConfig;

        const Block = props => (
            <Container
                padding={['bottom', 'top']}
                id={props.id}
                background={props.background}>
                <GridBlock
                    align="center"
                    contents={props.children}
                    layout={props.layout}
                />
            </Container>
        );

        const FeatureCallout = () => (
            <div className="productShowcaseSection paddingTop paddingBottom" style={{textAlign: 'center'}}>
                <h2>Features</h2>
                <div className="feature-callout-list" style={{textAlign: 'left'}}>
                    <ul>
                        <li>Sends requests continuously for X seconds</li>
                        <li>Supports REST and gRPC</li>
                        <li>Supports HTTP/1.1, HTTP/2 and HTTP/2 Cleartext protocol for REST</li>
                        <li>Supports custom HTTP headers</li>
                        <li>Supports concurrent requests</li>
                        <li>Supports placeholders for random elements in requests</li>
                        <li>Provides files or/and endpoints that can be used as liveness/readiness probes in Kubernetes</li>
                    </ul>
                </div>
            </div>
        );

        const Features = () => (
            <Block layout="fourColumn">
                {[
                    {
                        content: 'Mittens can run as a standalone cmd tool to send requests over REST/gRPC.',
                        image: `${baseUrl}img/cmd.svg`,
                        imageAlign: 'top',
                        title: 'a standalone cmd tool',
                    },
                    {
                        content: 'You can also run it as a Docker container linked to your main app.',
                        image: `${baseUrl}img/docker.png`,
                        imageAlign: 'top',
                        title: 'a linked Docker Container',
                    },
                    {
                        content: 'Or even as a sidecar on Kubernetes, in the pod where your app\'s running.',
                        image: `${baseUrl}img/kubernetes.png`,
                        imageAlign: 'top',
                        title: 'a Sidecar on Kubernetes',
                    },
                ]}
            </Block>
        );

        const Showcase = () => {
            if ((siteConfig.users || []).length === 0) {
                return null;
            }

            const showcase = siteConfig.users
                .filter(user => user.pinned)
                .map(user => (
                    <a href={user.infoLink} key={user.infoLink}>
                        <img src={user.image} alt={user.caption} title={user.caption}/>
                    </a>
                ));

            const pageUrl = page => baseUrl + (language ? `${language}/` : '') + page;

            return (
                <div className="productShowcaseSection paddingTop paddingBottom">
                    <h2>Who is Using This?</h2>
                    <p>This project is used by:</p>
                    <div className="logos">{showcase}</div>
                    <div className="more-users">
                        <a className="button" href={pageUrl('users.html')}>
                            More {siteConfig.title} Users
                        </a>
                    </div>
                </div>
            );
        };

        return (<div>
                <HomeSplash siteConfig={siteConfig} language={language}/>
                <div className="mainContainer homePage" style={{textAlign: 'center'}}>
                    <div className="running lightBackground">
                        <h1>Run it as:</h1>
                        <Features/>
                    </div>
                    <FeatureCallout/>
                    <div className="lightBackground">
                        <Showcase/>
                    </div>
                </div>
            </div>
        );
    }
}

module.exports = Index;
