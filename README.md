# mvisonneau/strongbox

[![PkgGoDev](https://pkg.go.dev/badge/github.com/mvisonneau/strongbox)](https://pkg.go.dev/mod/github.com/mvisonneau/strongbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/strongbox)](https://goreportcard.com/report/github.com/mvisonneau/strongbox)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/strongbox.svg)](https://hub.docker.com/r/mvisonneau/strongbox/)
[![test](https://github.com/mvisonneau/strongbox/actions/workflows/test.yml/badge.svg)](https://github.com/mvisonneau/strongbox/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/strongbox/badge.svg?branch=main)](https://coveralls.io/github/mvisonneau/strongbox?branch=main)
[![release](https://github.com/mvisonneau/strongbox/actions/workflows/release.yml/badge.svg)](https://github.com/mvisonneau/strongbox/actions/workflows/release.yml)
[![strongbox](https://snapcraft.io/strongbox/badge.svg)](https://snapcraft.io/strongbox)

Securely store secrets at rest using [Hashicorp Vault](https://www.vaultproject.io/).

## Concept

Vault is really good at safely storing data. Allowing us to query an HTTP endpoint in order to perform actions against our sensitive values. The goal of this project is to give us an additional abstraction layer, allowing us to easily store and versionize those secrets outside of it.

The idea is to leverage the [Vault Transit Secret backend](https://www.vaultproject.io/docs/secrets/transit/) in order to cipher/decipher our secrets and store them securely. The goal being, end up with a file that can easily be used in order to recover a lost secret/key. As well as storing it safely into a **git repository** for instance.

## Compatibility

`strongbox` supports **both version 1 and 2** of the [Vault K/V](https://www.vaultproject.io/api/secret/kv/kv-v1.html)

## Installation

Have a look onto the [latest release page](https://github.com/mvisonneau/strongbox/releases/latest) and pick your flavor. The exhaustive list of os/archs binaries we are releasing can be found in [here](https://github.com/mvisonneau/strongbox/blob/main/.goreleaser.yml#L8-16).

### Go

```bash
~$ go install github.com/mvisonneau/strongbox/cmd/strongbox@latest
```

### Homebrew

```bash
~$ brew install mvisonneau/tap/strongbox
```

### Snapcraft

```bash
~$ snap install strongbox
```

### Docker

```bash
~$ docker run -it --rm docker.io/mvisonneau/strongbox
~$ docker run -it --rm ghcr.io/mvisonneau/strongbox
~$ docker run -it --rm quay.io/mvisonneau/strongbox
```

### Scoop

```bash
~$ scoop bucket add https://github.com/mvisonneau/scoops
~$ scoop install strongbox
```

### Binaries, DEB and RPM packages

For the following ones, you need to know which version you want to install, to fetch the latest available :

```bash
~$ export STRONGBOX_VERSION=$(curl -s "https://api.github.com/repos/mvisonneau/strongbox/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
```

```bash
# Binary (eg: freebsd/amd64)
~$ wget https://github.com/mvisonneau/strongbox/releases/download/${STRONGBOX_VERSION}/strongbox_${STRONGBOX_VERSION}_freebsd_arm64.tar.gz
~$ tar zxvf strongbox_${STRONGBOX_VERSION}_freebsd_amd64.tar.gz -C /usr/local/bin

# DEB package (eg: linux/386)
~$ wget https://github.com/mvisonneau/strongbox/releases/download/${STRONGBOX_VERSION}/strongbox_${STRONGBOX_VERSION}_linux_386.deb
~$ dpkg -i strongbox_${STRONGBOX_VERSION}_linux_386.deb

# RPM package (eg: linux/arm64)
~$ wget https://github.com/mvisonneau/strongbox/releases/download/${STRONGBOX_VERSION}/strongbox_${STRONGBOX_VERSION}_linux_arm64.rpm
~$ rpm -ivh strongbox_${STRONGBOX_VERSION}_linux_arm64.rpm
```

## TL;DR - Get started and try it out in less than a minute

- Prereqs : **git**, **make** and **docker**

If you want to have a quick look and see how it works and/or you don't already have am operational **Vault cluster**, you can easily spin up a complete test environment:

```bash
# Install
~$ git clone git@github.com:mvisonneau/strongbox.git
~$ make dev-env
~$ make install

# Example commands to start with
~$ strongbox init
~$ strongbox transit create test
~$ strongbox secret write mysecret -k mykey -v sensitive_value
~$ strongbox status

# You can input longer strings using the '-' keyword for stdin input
~$ strongbox secret write mysecret -k verylong -v - <<EOF
THIS
IS
A
VERY
LONG
STRING
EOF
```

## Usage

```bash
~$ strongbox
NAME:
   strongbox - Manage Hashicorp Vault secrets at rest

USAGE:
   strongbox [global options] command [command options] [arguments...]

COMMANDS:
   transit  perform actions on transit key/backend
   secret   perform actions on secrets (locally)
   init     Create a empty state file at configured location
   status   display current status
   plan     compare local version with vault cluster
   apply    synchronize vault managed secrets
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --state FILE, -s FILE    load state from FILE (default: ".strongbox_state.yml") [$STRONGBOX_STATE]
   --vault-addr value       vault address (default: "http://vault.example.lan:8200") [$VAULT_ADDR]
   --vault-token value      vault token [$VAULT_TOKEN]
   --vault-role-id value    vault role id [$VAULT_ROLE_ID]
   --vault-secret-id value  vault secret id [$VAULT_SECRET_ID]
   --log-level value        log level (debug,info,warn,fatal,panic) (default: "info") [$STRONGBOX_LOG_LEVEL]
   --log-format value       log format (json,text) (default: "text") [$STRONGBOX_LOG_FORMAT]
   --help, -h               show help (default: false)
```

## Use case

The project was initially led in order to answer the following use case. I was willing to store sensitive secrets into Vault but also keep track on who was doing what onto them as well as being able to easily roll back the entire state of my set of secrets.

### Prerequisites

There are 3 mandatory configuration flags/environment variables to set to get started:

```bash
~$ strongbox | grep -E '_state|vault-addr|vault-token'
   --state FILE, -s FILE  load state from FILE (default: "~/.strongbox_state.yml") [$STRONGBOX_STATE]
   --vault-addr value     vault endpoint [$VAULT_ADDR]
   --vault-token value    vault token [$VAULT_TOKEN
```

Let's configure them:

```bash
~$ export VAULT_ADDR=https://vault.example.com:8200/
~$ export VAULT_TOKEN=9c9a9877-65e6-acea-8bdf-c1f0e959117f
~$ export STRONGBOX_STATE=/tmp/state.yml
```

### Run it

#### Init

In order to check you configuration, you can use this handy command : `strongbox status`

```bash
~$ strongbox status
State file not found at location: /tmp/state.yml, use 'strongbox init' to generate an empty one.
```

As you can see, the file doesn't exist yet, we can ask `strongbox` to create it for us :

```bash
~$ strongbox init
Creating an empty state file at /tmp/state.yml
```

The `status` command should be a bit more verbose now :

```bash
[STRONGBOX STATE]
+-------------+---------+
| Transit Key | default |
| KV Path     | secret/ |
| KV Version  |       2 |
| Secrets #   |       3 |
+-------------+---------+
[VAULT]
+-----------------+--------------------------------------+
| Sealed          | false                                |
| Cluster Version | 1.7.2                                |
| Cluster ID      | 5198332c-893c-ebbd-fdcf-82d3cdb47e4a |
| Secrets #       |                                    3 |
+-----------------+--------------------------------------+
```

#### Transit Key

If you want to reuse an existing key, you can use the following commands:

```bash
# List the available keys from the Vault endpoint
~$ strongbox transit list
+-------+
|  KEY  |
+-------+
| foo   |
| bar   |
+-------+

# Pick one of them
~$ strongbox transit use foo
```

Otherwise, `strongbox` can generate and use a new one for you:

```bash
~$ strongbox transit create test
Transit key created successfully
```

#### KV Path & Version

The **KV path** value is where you actually want to store the secrets onto Vault. This is only required when you're planning on keeping your locally configuration in sync with Vault. If you only want to leverage the Transit encryption capabilities you can skip this part.

By default, it manages the root of the `secret/` mountpoint, it is advised to use a more specific location at scale as `strongbox` would by default remove the values it doesn't manage in the **KV path**.

```bash
~$ strongbox kv get-path
secret/
~$ strongbox kv set-path secret/test/
~$ strongbox kv get-path
secret/test/
```

It is also paramount to set the correct version your KV is configured with. Default will be `2`

```bash
~$ strongbox kv get-version
2
~$ strongbox kv set-version 1
~$ strongbox kv get-version
1
```

#### Manage Secrets (the whole point!)

You are now all set to start managing secrets. Lets start by adding a few of them:

```bash
# Add defined values
~$ strongbox secret write foo -k key -v sensitive
~$ strongbox secret write foo -k key2 -v sensitive2

# Use a masked input
~$ strongbox secret write foo -k key3 -V
Sensitive
Enter a value: ***************

# Or generate random ones
~$ strongbox secret write bar -k key -r 8
```

You can now list all your secrets to see what they look like:

```bash
~$ strongbox secret list
[bar]
+-----+-------------------------------------------------------------+
| key | {{s5:zlU7fluN7E1/6qrjGG620KGhzE36SWyBeaNOU151eS9rkNfN1w==}} |
+-----+-------------------------------------------------------------+
[foo]
+------+-------------------------------------------------------------+
| key  | {{s5:zl2idnXPPwzD/zI2GSc+wVbxCjit5jI6W+f/ps/8hpNsaJf06g==}} |
| key2 | {{s5:Gil6RwgToO7ID/Xgewfvzu1Q/dnVH85mKu5XAEvIhUZGW1X+lzM=}} |
| key3 | {{s5:ISeYexNfD0gFXF2qoEoQfSqzZUlH5DvQ/DO86YfRNhW8D24uw0Q=}} |
+------+-------------------------------------------------------------+
```

If you want you can also take a look at what your state file looks like:

```bash
~$ cat /tmp/state.yml
vault:
  transitkey: test
  secretpath: secret/test/
secrets:
  bar:
    key: {{s5:zlU7fluN7E1/6qrjGG620KGhzE36SWyBeaNOU151eS9rkNfN1w==}}
  foo:
    key: {{s5:zl2idnXPPwzD/zI2GSc+wVbxCjit5jI6W+f/ps/8hpNsaJf06g==}}
    key2: {{s5:Gil6RwgToO7ID/Xgewfvzu1Q/dnVH85mKu5XAEvIhUZGW1X+lzM=}}
    key3: {{s5:ISeYexNfD0gFXF2qoEoQfSqzZUlH5DvQ/DO86YfRNhW8D24uw0Q=}}
```

A choice has been made to keep the secrets and keys readable in order to be able to review changes in PR/MRs. As you can see otherwise, you now have a perfectly shareable/commitable.

#### Read secrets

In order to read the secrets, you can use this function:

```bash
~$ strongbox secret read foo -k key
sensitive
~$ strongbox secret read bar -k key
4Qco_ndx
```

#### Syncing with Vault

If you want to be able to access those secrets directly from Vault, `strongbox` allows you to easily sync them with your cluster and maintain their state. You also have the capability of knowing what actions `strongbox` is going to perform before actually running the changes.

```bash
~$ strongbox plan
Add/Update: 2 secret(s) and 4 key(s)
=> foo:key
=> foo:key2
=> foo:key3
=> bar:key
```

```bash
~$ strongbox apply
=> Added/Updated secret 'foo' and managed keys
=> Added/Updated secret 'bar' and managed keys
```

FYI, the values that we store in Vault are deciphered. You can check that they have been correctly created using the Vault API or the Vault client :

```bash
~$ vault list secret/test
Keys
----
bar
foo

~$ vault read secret/test/foo
Key               Value
---               -----
refresh_interval  768h0m0s
key               sensitive
key2              sensitive2
key3              sensitive3
```

#### Rotate secrets

If you feel that you need to rotate the encryption of your state file or that the transit you are using might have been compromised, `strongbox` allows you to easily do it.

```bash
# Check your current key name
~ strongbox transit info | grep name | cut -d'|' -f 3 | xargs
old

# Create a new key
~$ strongbox transit create new
Transit key created successfully

# Rotate!
~$ strongbox state rotate-from old
Rotated secrets from 'old' to 'new'
```

## Develop / Test

If you use docker, you can easily get started using :

```bash
~$ make dev-env
# You should then be able to use go commands to work onto the project, eg:
~$ make install
~$ strongbox
```

This command will spin up a Vault container and a build one with everything required in terms of go dependencies in order to get started.

## Build / Release

If you want to build and/or release your own version of `strongbox`, you need the following prerequisites :

- [git](https://git-scm.com/)
- [golang](https://golang.org/)
- [make](https://www.gnu.org/software/make/)
- [goreleaser](https://goreleaser.com/)

```bash
~$ git clone git@github.com:mvisonneau/strongbox.git && cd strongbox

# Build the binaries locally
~$ make build

# Build the binaries and release them (you will need a GITHUB_TOKEN and to reconfigure .goreleaser.yml)
~$ make release
```

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/strongbox/pulls).
