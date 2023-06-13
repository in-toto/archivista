var graphqlConfig;
try {
  graphqlConfig = require("./.cache/typegen/graphql.config.json");
} catch (error) {
  graphqlConfig = null;
}

module.exports = graphqlConfig;
