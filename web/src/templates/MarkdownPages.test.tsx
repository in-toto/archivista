/* eslint-disable react/display-name */
import '@testing-library/jest-dom';

import * as React from 'react';

import Layout from '../layouts/DefaultLayout';
import MarkdownPage from './MarkdownPages';
import { graphql } from 'gatsby';
import { render } from '@testing-library/react';

jest.mock('../layouts/DefaultLayout', () => ({ children }: { children: React.ReactNode }) => <div id="default-layout">{children}</div>);

jest.mock('gatsby', () => {
  const graphqlMock = jest.fn();
  return {
    graphql: graphqlMock,
    useStaticQuery: jest.fn(),
  };
});

const mockData = {
  markdownRemark: {
    frontmatter: {
      title: 'Test Title',
      date: '2022-04-01',
    },
    html: '<p>Test content</p>',
  },
};

// jest.mock('../graphql/markdownRemarkQuery.graphql', () => '');

describe('MarkdownPage component', () => {
  beforeAll(() => {
    (graphql as jest.Mock).mockImplementation((query: string) => {
      if (query === '') {
        return Promise.resolve(mockData);
      } else {
        return Promise.reject('Invalid query');
      }
    });
  });

  afterAll(() => {
    jest.resetAllMocks();
  });

  it('should render the page title and content', () => {
    const { getByText } = render(<MarkdownPage data={mockData} />);

    expect(getByText(mockData.markdownRemark.frontmatter.title)).toBeInTheDocument();
    expect(getByText(mockData.markdownRemark.frontmatter.date)).toBeInTheDocument();
    expect(getByText('Test content')).toBeInTheDocument();
  });
});
