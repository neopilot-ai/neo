<p align="center">
  <a href="https://neo.khulnasoft.com/">
    <img alt="NEO" src="https://raw.githubusercontent.com/neo/identity/main/variants/neo-full.svg" width="300" />
  </a>
</p>

<p align="center">
  <a href="https://neo.khulnasoft.com/discord"><img alt="Discord" src="https://img.shields.io/discord/983865673656705025?style=flat-square&label=Discord" /></a>
  <a href="https://www.npmjs.com/package/neo"><img alt="npm" src="https://img.shields.io/npm/v/neo.svg?style=flat-square" /></a>
  <a href="https://github.com/neopilot-ai/neo/actions/workflows/build.yml"><img alt="Build status" src="https://img.shields.io/github/actions/workflow/status/neo/neo/build.yml?style=flat-square&branch=dev" /></a>
</p>

---

Build full-stack apps on your own infrastructure.

NEO v3 uses a new engine for deploying NEO apps. It uses Pulumi and Terraform, as opposed to CDK and CloudFormation. [Read the full announcement here](https://neo.khulnasoft.com/blog/neo-v3).

## Installation

If you are using NEO as a part of your Node project, we recommend installing it locally.

```bash
npm install neo
```

If you are not using Node, you can install the CLI globally.

```bash
curl -fsSL https://neo.khulnasoft.com/install | bash
```

To install a specific version.

```bash
curl -fsSL https://neo.khulnasoft.com/install | VERSION=0.0.403 bash
```

To use a package manager, [check out our docs](https://neo.khulnasoft.com/docs/reference/cli/).

#### Manually

Download the pre-compiled binaries from the [releases](https://github.com/neopilot-ai/neo/releases/latest) page and copy to the desired location.

## Get Started

Get started with your favorite framework:

- [Next.js](https://neo.khulnasoft.com/docs/start/aws/nextjs)
- [Remix](https://neo.khulnasoft.com/docs/start/aws/remix)
- [Astro](https://neo.khulnasoft.com/docs/start/aws/astro)
- [API](https://neo.khulnasoft.com/docs/start/aws/api)

## Learn More

Learn more about some of the key concepts:

- [Live](https://neo.khulnasoft.com/docs/live)
- [Linking](https://neo.khulnasoft.com/docs/linking)
- [Console](https://neo.khulnasoft.com/docs/console)
- [Components](https://neo.khulnasoft.com/docs/components)

## Contributing

Here's how you can contribute:

- Help us improve our docs
- Find a bug? Open an issue
- Feature request? Submit a PR 

## Running Locally

1. Clone the repo
2. `bun install`
3. `go mod tidy`
4. `cd platform && bun run build`

Now you can run the CLI locally on any of the `examples/` apps.

```bash
cd examples/aws-api
go run ../../cmd/neo <command>
```

If you want to build the CLI, you can run `go build ./cmd/neo` from the root. This will create a
`neo` binary that you can use.

For building the docs, you need to run `bun generate` and `bun dev` inside the `www` directory.

---

**Join our community** [Discord](https://neo.khulnasoft.com/discord) | [YouTube](https://www.youtube.com/c/neo-dev) | [X.com](https://x.com/NEO_dev)
