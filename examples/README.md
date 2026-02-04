# Examples

A collection of example NEO apps. You can also [view these in our docs](https://neo.khulnasoft.com/docs/examples/).

## Generated Docs

The comments in the `neo.config.ts` of these apps are used to generate the doc. Here's an [example comment block](/examples/aws-info/neo.config.ts).

```
/**
 * ## Current AWS account
 *
 * You can use the `aws.getXXXXOutput()` provider functions to get info about the current
 * AWS account.
 * Learn more about [provider functions](/docs/providers/#functions).
 */
```

The [generate script](/www/generate.ts) looks for examples that have a comment block like this in their `neo.config.ts`. It'll extract this add it as a section to the examples doc.

## Contributing

To contribute an example or to edit one, submit a PR. Make sure to document the `neo.config.ts` in your example if you want to add it to the docs.
