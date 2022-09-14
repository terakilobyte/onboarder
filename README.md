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

Optionally, pass in the path to your public ssh key for uploading to github.

```sh
onboarder -s path_to_public_key -c path_to_config -o output_directory
```

The above will fork repos then clone
them to the *~/work* directory.
