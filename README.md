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
        token: ${{ secrets.GITHUB_TOKEN }}
        scope: ${{ env.GITHUB_REPOSITORY }}
        secret: KUBECONFIG
        file: kubeconfig.yml
```

## Input variables

See [action.yml](./action.yml) for more detailed information.

* `token` - github token for repositories or personal access token for organisations
* `scope` - name of the organisation or username and repository
* `secret` - name of the secret
* `file` - file to be stored in secret

## Contributing

We would ❤️ for you to contribute to `nicklasfrahm/file-secret-action`, pull requests are welcome!

## License

This project is licensed under the [MIT license](./LICENSE.md).
