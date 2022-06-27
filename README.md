# Onboarder

Onboarder is an onboarding tool built for the Docs team at MongoDB (initially).

## Install

The easiest way to install is through *homebrew*

```sh
brew tap terakilobyte/tools
brew install onboarder
```

## Use

```sh
onboarder --help
```

Run onboarder, passing in flags for the output directory where repositories
should be cloned to.

You also need to specify your gpg key to enable signed commits. Follow the instructions on
[Generate a GPG key](https://docs.github.com/en/enterprise-cloud@latest/authentication/managing-commit-signature-verification/generating-a-new-gpg-key)
to create one.

Also, pass in the path to your public ssh key for uploading to github.

```sh
onboarder -g gid -s path_to_public_key -c path_to_config ~/work
```

The above will fork repos then clone
them to the *~/work* directory.
