/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
import { ConfigConstants } from './src/shared/constants';
import { GatsbyConfig } from 'gatsby';
import { createProxyMiddleware } from 'http-proxy-middleware';

const config: GatsbyConfig = {
  developMiddleware: (app) => {
    let archivistaUrl = 'https://archivista.testifysec.io';
    // let archivistaUrl = 'https://archivista.testifysec.localhost';
    let judgeApiUrl = 'https://judge-api.testifysec.localhost';
    let kratosUrl = 'https://kratos.testifysec.localhost';
    let loginUrl = 'https://login.testifysec.localhost';

    // Check if --remote-proxy flag is passed
    // TODO: these should be nonprod
    if (process.env.GATSBY_REMOTE_PROXY === 'true') {
      archivistaUrl = 'https://archivista.testifysec.io';
      judgeApiUrl = 'https://judge-api.testifysec.io';
      kratosUrl = 'https://kratos.testifysec.io';
      loginUrl = 'https://login.testifysec.io';
    }

    app.use(
      createProxyMiddleware('/archivista', {
        target: archivistaUrl,
        changeOrigin: true,
        secure: false,
        pathRewrite: {
          '^/archivista': '',
        },
      }),
      createProxyMiddleware('/judge-api', {
        target: judgeApiUrl,
        changeOrigin: true,
        secure: false,
        pathRewrite: {
          '^/judge-api': '',
        },
      }),
      createProxyMiddleware('/kratos', {
        target: kratosUrl,
        changeOrigin: true,
        secure: false,
        pathRewrite: {
          '^/kratos': '',
        },
      }),
      createProxyMiddleware('/login', {
        target: loginUrl,
        changeOrigin: true,
        secure: false,
        pathRewrite: {
          '^/login': '',
        },
      })
    );
  },
  siteMetadata: {
    title: 'Judge Platform',
    siteUrl: 'https://judge.testifysec.com',
  },
  graphqlTypegen: true,
  plugins: [
    {
      resolve: 'gatsby-plugin-apollo',
      options: {
        uri: `${ConfigConstants.archivista_svc}`,
      },
    },
    {
      resolve: 'gatsby-plugin-svgr-loader',
      options: {
        inlineSvgOptions: {
          /* options here */
          jsx: true,
        },
        rule: {
          include: /\.inline\.svg$/,
        },
      },
    },
    {
      resolve: 'gatsby-plugin-google-fonts-v2',
      options: {
        fonts: [
          {
            family: 'Roboto',
            weights: ['100', '400'],
          },
        ],
      },
    },
    'gatsby-plugin-material-ui',
    'gatsby-plugin-image',
    'gatsby-transformer-yaml',
    'gatsby-plugin-react-helmet',
    'gatsby-plugin-testing',
    'gatsby-plugin-sitemap',
    'gatsby-theme-material-ui',
    {
      resolve: 'gatsby-plugin-manifest',
      options: {
        icon: 'src/images/icon.png',
      },
    },
    'gatsby-transformer-remark',
    'gatsby-plugin-sharp',
    'gatsby-transformer-sharp',
    {
      resolve: 'gatsby-source-filesystem',
      options: {
        name: 'images',
        path: './src/images/',
      },
      __key: 'images',
    },
    {
      resolve: 'gatsby-source-filesystem',
      options: {
        name: 'pages',
        path: './src/pages/',
      },
      __key: 'pages',
    },
    {
      resolve: 'gatsby-source-filesystem',
      options: {
        name: 'images',
        path: './src/images/',
      },
      __key: 'images',
    },
    {
      resolve: 'gatsby-source-filesystem',
      options: {
        name: 'markdown-pages',
        path: `${__dirname}/src/markdown-pages`,
      },
    },
    {
      resolve: 'gatsby-source-filesystem',
      options: {
        path: './src/content/',
      },
    },
  ],
};

export default config;
