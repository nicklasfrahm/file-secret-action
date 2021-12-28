# 🔑 File to Secret for GitHub Actions (deprecated)

[![Go Report Card](https://goreportcard.com/badge/github.com/nicklasfrahm/file-secret-action)](https://goreportcard.com/report/github.com/nicklasfrahm/file-secret-action)
[![container](https://github.com/nicklasfrahm/file-secret-action/actions/workflows/container.yml/badge.svg?branch=main)](https://github.com/nicklasfrahm/file-secret-action/actions/workflows/container.yml)

A [GitHub Action](https://github.com/features/actions) to upload a file as a GitHub Actions Secret.

## Deprecation notice

This repository is deprecated. I recommend to use the Github CLI instead and set up the following step:

```yaml
...
      - name: Upload kubeconfig as pipeline secret
        run: |
          echo ${{ secrets.PERSONAL_ACCESS_TOKEN }} | gh auth login --with-token
          gh secret set KUBECONFIG -R ${{ github.repository }} -e ${{ env.ENVIRONMENT }} < kubeconfig.yml
...
```

## Usage

```yaml
name: file-secret

on:
  - push

jobs:
  file-secret:
    name: File Secret
    runs-on: ubuntu-latest
    steps:
      - name: Check out repository
        uses: actions/checkout@master

      - name: Upload kubeconfig
        uses: nicklasfrahm/file-secret-action@main
        with:
          token: ${{ secrets.PAT }}
          scope: ${{ github.repository }}
          secret: KUBECONFIG
          file: kubeconfig.yml
```

## Input variables

See [action.yml](./action.yml) for more detailed information.

* `token` - personal access token with the `workflow` permission, the standard `GITHUB_TOKEN` does not work
* `scope` - name of the organisation or username and repository
* `secret` - name of the secret
* `file` - file to be stored in secret
* `visibility` - visibility of the secret within the organization, currently only _private_ and _all_ are supported

## Contributing

This repository is archived and therefore frozen. Please maintain your own copies and forks.

## License

This project is licensed under the [MIT license](./LICENSE.md).
