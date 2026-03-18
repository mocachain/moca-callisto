const { ApolloClient, InMemoryCache, gql } = require('@apollo/client');

// define GraphQL query
const BUCKETS_QUERY = gql`
  query Buckets($limit: Int = 10, $offset: Int = 0) {
    buckets: bucket(limit: $limit, offset: $offset, order_by: { create_at: desc }) {
      id
      bucket_name
      bucket_status
      charged_read_quota
      create_at
      global_virtual_group_family_id
      owner
      payment_address
      source_type
      sp_as_delegated_agent_disabled
      tags
      visibility
    }
  }
`;

// create Apollo client instance
const client = new ApolloClient({
  uri: 'http://localhost:8080/v1/graphql',
  cache: new InMemoryCache(),
});

// query test
async function fetchBuckets(limit = 2, offset = 0) {
  try {
    const result = await client.query({
      query: BUCKETS_QUERY,
      variables: { limit, offset },
    });

    return result.data.buckets;
  } catch (error) {
    console.error('fetchBuckets error: ', error);
    return [];
  }
}

const main = async () => {
  try {
    const buckets = await fetchBuckets();
    console.log(buckets);
  } catch (error) {
    console.error('error', error);
  }
};

main();
