# 🔑 File to Secret for GitHub Actions

[![Go Report Card](https://goreportcard.com/badge/github.com/nicklasfrahm/file-secret-action)](https://goreportcard.com/report/github.com/nicklasfrahm/file-secret-action)

A [GitHub Action](https://github.com/features/actions) to upload a file as a GitHub Actions Secret.

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

* `token` - personal access token with the `repo` permission, the standard `GITHUB_TOKEN` does not work
* `scope` - name of the organisation or username and repository
* `secret` - name of the secret
* `file` - file to be stored in secret
* `visibility` - visibility of the secret within the organization, currently only _private_ and _all_ are supported

## Contributing

We would ❤️ for you to contribute to `nicklasfrahm/file-secret-action`, pull requests are welcome!

## License

This project is licensed under the [MIT license](./LICENSE.md).
