# Elixir CI Email (GitHub action)

This is a simple GitHub action that allows to send emails when the GitHub Check Suite completes.

## Getting started

The action expects this environment variables:

  * `APP_NAME` - Name of an application for which to send emails for.
  * `MAIL_FROM` - email address to send emails from.
  * `MAIL_HOST` - SMTP host to send emails to.
  * `MAIL_USERNAME` and `MAIL_PASSWORD` - username and password to authorize with the SMTP server.
  * `GITHUB_TOKEN` - is standard environment variable for GitHub actions.
  * optional `IGNORED_CONCLUSIONS` to specify conclusions to report. By default only `success` and `neutral` checks are ignored.

Now your action can look liker this in your `.github/workflows/ci_email.yml` workflow file:

```yaml
on: check_suite
name: CI email
jobs:
  sendEmail:
    name: Send email
    runs-on: ubuntu-latest
    steps:
    - name: Send email
      # Source: https://github.com/elixir-lang/elixir-ci
      uses: docker://fertapric/elixir-ci-email:latest
      env:
        APP_NAME: Cirrus CI
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        MAIL_FROM: ci@elixir-lang.org
        MAIL_HOST: smtp.sendgrid.net
        MAIL_PASSWORD: ${{ secrets.MAIL_PASSWORD }}
        MAIL_USERNAME: ${{ secrets.MAIL_USERNAME }}
```

## Releasing a new version

These are the steps to release a new version of the action:

```shell
docker build -t fertapric/elixir-ci-email:latest .
docker push fertapric/elixir-ci-email:latest
```
