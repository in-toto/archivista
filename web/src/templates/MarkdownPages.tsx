import * as React from 'react';

import { Container, Typography } from '@mui/material';

import { Helmet } from 'react-helmet';
import Layout from '../layouts/DefaultLayout';
import { graphql } from 'gatsby';

export type MarkdownPageProps = {
  data: {
    markdownRemark: {
      frontmatter: {
        title: string;
        date: string;
      };
      html: string;
    };
  };
};

const MarkdownPage: React.FC<MarkdownPageProps> = ({ data: { markdownRemark } }: MarkdownPageProps) => {
  return (
    <Layout>
      <>
        <Helmet>
          <title>{markdownRemark.frontmatter.title}</title>
        </Helmet>
        <Container maxWidth="md">
          <Typography variant="h1" gutterBottom>
            {markdownRemark.frontmatter.title}
          </Typography>
          <Typography variant="h2" gutterBottom>
            {markdownRemark.frontmatter.date}
          </Typography>
          <div dangerouslySetInnerHTML={{ __html: markdownRemark.html }} />
        </Container>
      </>
    </Layout>
  );
};

export default MarkdownPage;

export const pageQuery = graphql`
  query MarkdownRemark($id: String!) {
    markdownRemark(id: { eq: $id }) {
      html
      frontmatter {
        date(formatString: "MMMM DD, YYYY")
        slug
        title
      }
    }
  }
`;
