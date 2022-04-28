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

### SSH

READ THIS. Onboarder generates a new ssh keypair and uploads the public
key to github for you. It will also add it to the ssh-agent, and it modifies
your `~/.ssh/config` file (creates if needed) to use the key.

Run onboarder, passing in flags for the output directory where repositories
should be cloned to, and which team you are on.

Current teams are `cet`, and `tdbx`.

```sh
onboarder -t tdbx -o ~/work
```

The above will fork repos appropriate for the *tdbx* team and then clone
them to the *~/work* directory.
