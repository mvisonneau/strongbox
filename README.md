# mvisonneau/strongbox

[![GoDoc](https://godoc.org/github.com/mvisonneau/strongbox?status.svg)](https://godoc.org/github.com/mvisonneau/strongbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/strongbox)](https://goreportcard.com/report/github.com/mvisonneau/strongbox)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/strongbox.svg)](https://hub.docker.com/r/mvisonneau/strongbox/)
[![Build Status](https://travis-ci.org/mvisonneau/strongbox.svg?branch=master)](https://travis-ci.org/mvisonneau/strongbox)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/strongbox/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/strongbox?branch=master)

Securely store secrets at rest using [Hashicorp Vault](https://www.vaultproject.io/).

## Concept

Vault is really good at safely storing data. Allowing us to query an HTTP endpoint in order to perform actions against our sensitive values. The goal of this project is to give us an additional abstraction layer, allowing us to easily store and versionize those secrets outside of it.

The idea is to leverage the [Vault Transit Secret backend](https://www.vaultproject.io/docs/secrets/transit/) in order to cipher/decipher our secrets and store them securely. The goal being, end up with a file that can easily be used in order to recover a lost secret/key. As well as storing it safely into a **git repository** for instance.

## Installation

- Go : `go get -u github.com/mvisonneau/strongbox`
- Docker : `docker fetch mvisonneau/strongbox`

## Usage

```bash
~$ strongbox
NAME:
   strongbox - Securely store secrets at rest with Hashicorp Vault

USAGE:
   strongbox [global options] command [command options] [arguments...]

VERSION:
   <devel>

COMMANDS:
     transit  perform actions on transit key/backend
     secret   perform actions on secrets (locally)
     init     Create a empty state file at configured location
     status   display current status
     plan     compare local version with vault cluster
     apply    synchronize vault managed secrets
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --state FILE, -s FILE  load state from FILE (default: "~/.strongbox_state.yml") [$STRONGBOX_STATE]
   --vault-addr value     vault endpoint [$VAULT_ADDR]
   --vault-token value    vault token [$VAULT_TOKEN]
   --log-level value      log level (debug,info,warn,fatal,panic) (default: "info") [$STRONGBOX_LOG_LEVEL]
   --log-format value     log format (json,text) (default: "text") [$STRONGBOX_LOG_FORMAT]
   --help, -h             show help
   --version, -v          print the version
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

### Run it!

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
[STATE]
+------------+---------+
| TransitKey |         |
| SecretPath | secret/ |
| Secrets #  |       0 |
+------------+---------+
[VAULT]
+-----------------+--------------------------------------+
| Sealed          | false                                |
| Cluster Version | 0.9.1                                |
| Cluster ID      | 420572b9-af2f-e0a6-2b40-c4dd449dd29a |
| Secrets #       |                                    0 |
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

#### Secret Path

The **secret_path** value is where you actually want to store the secrets onto Vault. This is only required when you're planning on keeping your locally configuration in sync with Vault. If you only want to leverage the Transit encryption capabilities you can skip this part.

By default, it manages the root of the `secret/` mountpoint, it is advised to use a more specific location at scale as `strongbox` would by default remove the values it doesn't manage in the **secret-path**.

```bash
~$ strongbox get-secret-path
secret/
~$ strongbox set-secret-path secret/test/
~$ strongbox get-secret-path
secret/test/
```

#### Manage Secrets (the whole point!)

You are now all set to start managing secrets. Lets start by adding a few of them:

```bash
~$ strongbox secret write foo key sensitive
~$ strongbox secret write bar key sensitive
~$ strongbox secret write foo key2 sensitive2
~$ strongbox secret write foo key3 sensitive3
```

You can now list all your secrets to see what they look like:

```bash
~$ strongbox secret list
[bar]
+-----+---------------------------------------------------------------+
| key | vault:v1:zlU7fluN7E1/6qrjGG620KGhzE36SWyBeaNOU151eS9rkNfN1w== |
+-----+---------------------------------------------------------------+
[foo]
+------+---------------------------------------------------------------+
| key  | vault:v1:zl2idnXPPwzD/zI2GSc+wVbxCjit5jI6W+f/ps/8hpNsaJf06g== |
| key2 | vault:v1:Gil6RwgToO7ID/Xgewfvzu1Q/dnVH85mKu5XAEvIhUZGW1X+lzM= |
| key3 | vault:v1:ISeYexNfD0gFXF2qoEoQfSqzZUlH5DvQ/DO86YfRNhW8D24uw0Q= |
+------+---------------------------------------------------------------+
```

If you want you can also take a look at what your state file looks like :

```bash
~$ cat /tmp/state.yml
vault:
  transitkey: test
  secretpath: secret/test/
secrets:
  bar:
    key: vault:v1:zlU7fluN7E1/6qrjGG620KGhzE36SWyBeaNOU151eS9rkNfN1w==
  foo:
    key: vault:v1:zl2idnXPPwzD/zI2GSc+wVbxCjit5jI6W+f/ps/8hpNsaJf06g==
    key2: vault:v1:Gil6RwgToO7ID/Xgewfvzu1Q/dnVH85mKu5XAEvIhUZGW1X+lzM=
    key3: vault:v1:ISeYexNfD0gFXF2qoEoQfSqzZUlH5DvQ/DO86YfRNhW8D24uw0Q=
```

A choice has been made to keep the secrets and keys readable in order to be able to review changes in PR/MRs. As you can see otherwise, you now have a perfectly shareable/commitable.

#### Read secrets

In order to read the secrets, you can use this function:

```bash
~$ strongbox secret read foo key
sensitive
~$ strongbox secret read foo key2
sensitive2
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
Key             	Value
---             	-----
refresh_interval	768h0m0s
key             	sensitive
key2            	sensitive2
key3            	sensitive3
```

## Develop

If you use docker, you can easily get started using :

```bash
~$ make dev-env
# You should then be able to use go commands to work onto the project
```

If you also need a development Vault endpoint to play with, you can spin a working one in a few seconds :

```bash
# Start container
~$ docker run -d --name vault vault
# Fetch its IP
~$ docker vault inspiring_kirch | jq -r '.[0].NetworkSettings.IPAddress'
172.17.0.3
# Fetch the root token
~$ docker logs vault 2>/dev/null | grep 'Root Token' | cut -d' ' -f3
6c53eac6-8c67-caa0-2838-ae206bd83095
```

You should then be able to export those value in the `strongbox` devel container :

```
export VAULT_ADDR=http://172.17.0.3:8200
export VAULT_TOKEN=6c53eac6-8c67-caa0-2838-ae206bd83095
```

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/strongbox).
