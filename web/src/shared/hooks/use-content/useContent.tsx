/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import { graphql, useStaticQuery } from 'gatsby';

export type ContentYaml = {
  signInButtonText: string;
};

export type IndexYaml = {
  noAuthHeading: string;
  noAuthSubheading: string;
};

export type ContentModels = {
  contentYaml: ContentYaml;
  contentEsMxYaml: ContentYaml;
  noAuthYaml: IndexYaml;
  noAuthEsMxYaml: IndexYaml;
};
/**
 * This hook provides a way for us to access all of our app content globally from graphql
 * The pattern also keeps all of our content away from our application code
 * // TODO: as this grows, it might be better to have the graphql queries at the page level instead of entirely global.
 *
 * @returns
 */
const useContent = () => {
  const { contentYaml, contentEsMxYaml, noAuthYaml, noAuthEsMxYaml }: ContentModels = useStaticQuery(
    graphql`
      query Content {
        contentEsMxYaml {
          signInButtonText
        }
        contentYaml {
          signInButtonText
        }
        noAuthYaml {
          noAuthHeading
          noAuthSubheading
        }
        noAuthEsMxYaml {
          noAuthSubheading
          noAuthHeading
        }
      }
    `
  );
  return { contentYaml, contentEsMxYaml, noAuthYaml, noAuthEsMxYaml } as ContentModels;
};

export default useContent;
