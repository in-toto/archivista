import { gql, useQuery } from '@apollo/client';

import { ApiStatus } from '../use-api-status/useApiStatus';
import { Dsse } from '../../../generated/graphql';

export type ArchivistaProps = {
  apiStatus: ApiStatus;
  results: Dsse[];
};

/**
 * A hook for calling Archivista via graphql
 *
 */
const useArchivistaByRepos = (repos: string[]): [ArchivistaProps] => {
  // const [getEnvelopes] = useBySubjectDigestLazyQuery();

  const GET_ATTESTATIONS_FOR_SUBJECT = gql`
    query DsseBySubjectNames($subjectNames: [String!]) {
      dsses(
        where: {
          hasSignaturesWith: { hasTimestampsWith: { timestampGT: "2023-04-01T00:00:00Z", timestampLT: "2023-05-01T00:00:00Z" } }
          hasStatementWith: { hasSubjectsWith: { nameIn: $subjectNames } }
        }
        first: 200
      ) {
        edges {
          node {
            id
            gitoidSha256

            statement {
              id
              attestationCollections {
                name
                attestations {
                  attestationCollection {
                    name
                  }
                }
              }
              subjects {
                edges {
                  node {
                    name
                  }
                }
              }
            }
            signatures {
              timestamps {
                timestamp
              }
            }
          }
        }
      }
    }
  `;

  const { loading, error, data } = useQuery(GET_ATTESTATIONS_FOR_SUBJECT, {
    variables: {
      subjectNames: repos.map((r) => `https://witness.dev/attestations/github/v0.1/projecturl:${r}`),
    },
    context: { uri: '/archivista/query' },
  });

  // console.log(JSON.stringify(data, null, 2));
  // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-unsafe-return, @typescript-eslint/no-unsafe-call
  const results = data?.dsses?.edges?.map((edge: { node: any }) => edge?.node) || [];
  console.log(JSON.stringify(results, null, 2));
  const apiStatus: ApiStatus = { isLoading: loading, hasError: error !== undefined };
  const archivistaProps: ArchivistaProps = { apiStatus, results };

  return [archivistaProps];
};

export default useArchivistaByRepos;
