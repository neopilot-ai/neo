export const database = new neo.aws.Dynamo("Database", {
  fields: {
    PK: "string",
    SK: "string",
  },
  primaryIndex: {
    hashKey: "PK",
    rangeKey: "SK",
  },
});
