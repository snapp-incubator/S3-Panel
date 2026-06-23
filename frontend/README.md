[![Latest Release](https://gitlab.snapp.ir/snappcloud/unified-panel/-/badges/release.svg)](https://gitlab.snapp.ir/snappcloud/unified-panel/-/releases)
[![coverage report](https://gitlab.snapp.ir/snappcloud/unified-panel/badges/main/coverage.svg)](https://gitlab.snapp.ir/snappcloud/unified-panel/-/commits/main)
[![pipeline status](https://gitlab.snapp.ir/snappcloud/unified-panel/badges/main/pipeline.svg)](https://gitlab.snapp.ir/snappcloud/unified-panel/-/commits/main)

# SnappCost PWA

> [!CAUTION]
> **This repository is DEPRECATED and no longer deployed.**
>
> The unified panel has been split up and its parts relocated:
>
> | Old feature | New home |
> | --- | --- |
> | **Cost panel** (usage, quota, billing, openstack) | Moved into the SnappCost backend repo and embedded in the binary — [`platform/snappcost`](https://gitlab.snapp.ir/platform/snappcost) (`web/`). Backend + frontend now ship in one image. |
> | **Home page service links** | Moved into the DevEx portal as the "Service Portals" page — [`platform/devexportal`](https://gitlab.snapp.ir/platform/devexportal) (`apps/frontend/src/docs/setup/service-portals.mdx`). |
> | **Object Storage (S3) browser** | Retired. |
>
> Do not develop here. The `cost.snappcloud.io` route now points at the `snappcost` service. This repo is kept only for history and will be archived.

SnappCost as SnappCloud billing solution use the cluster metrics to find out about
usage and provides bills.

## Configuration

Create a .env file based on the provided .env.example.

## Usage

Run the development server:

```sh
bun run dev
# or
npm run dev
# or
yarn dev
```

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react/README.md) uses [Babel](https://babeljs.io/) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type aware lint rules:

- Configure the top-level `parserOptions` property like this:

```js
export default {
  // other rules...
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
    project: ['./tsconfig.json', './tsconfig.node.json'],
    tsconfigRootDir: __dirname
  }
}
```

- Replace `plugin:@typescript-eslint/recommended` to `plugin:@typescript-eslint/recommended-type-checked` or `plugin:@typescript-eslint/strict-type-checked`
- Optionally add `plugin:@typescript-eslint/stylistic-type-checked`
- Install [eslint-plugin-react](https://github.com/jsx-eslint/eslint-plugin-react) and add `plugin:react/recommended` & `plugin:react/jsx-runtime` to the `extends` list
