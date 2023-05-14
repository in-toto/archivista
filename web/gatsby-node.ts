/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
import { GatsbyNode } from 'gatsby';
import path from 'path';

interface MarkdownPageNode {
  id: string;
  frontmatter: {
    title: string;
    slug: string;
  };
}

interface QueryResult {
  data: {
    allMarkdownRemark: {
      nodes: MarkdownPageNode[];
    };
  };
  errors: any;
}

export const createPages: GatsbyNode['createPages'] = async ({ actions, graphql, reporter }) => {
  const { createPage } = actions;

  const { data, errors } = (await graphql(`
    query AllMarkdown {
      allMarkdownRemark {
        nodes {
          id
          frontmatter {
            title
            slug
          }
        }
      }
    }
  `)) as QueryResult;

  if (errors) {
    reporter.panicOnBuild('Error while running GraphQL query.');
    return;
  }

  const markdownPages = data.allMarkdownRemark.nodes;

  markdownPages.forEach((node: any) => {
    createPage({
      path: node.frontmatter.slug,
      component: path.resolve('./src/templates/MarkdownPages.tsx'),
      context: {
        id: node.id,
        title: node.frontmatter.title,
        slug: node.frontmatter.slug,
      },
    });
  });
};
