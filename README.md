<div align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./logo-dark.svg" style="max-width: 50%">
  <source media="(prefers-color-scheme: light)" srcset="./logo-light.svg" style="max-width: 50%">
  <img alt="enc logo" src="./logo-light.svg" style="max-width: 50%">
</picture>
    <h1>enc</h1>
    <p><b>a modern and friendly alternative to GnuPG</b></p>
</div>

# What is enc?

Enc is a CLI tool for encryption, a modern and friendly alternative to [GnuPG](https://gnupg.org/). It is easy to use, secure by default and can encrypt and decrypt files using password or encryption keys, manage and download keys, and sign data. Our goal was to make encryption available to all engineers without the need to learn a lot of new words, concepts, and commands. It is the most beginner-friendly CLI tool for encryption, and keeping it that way is our top priority.

## Features

+ **Easy installation**. Grab the binary, and you're ready to go.
+ **Friendly CLI**. We use well-isolated subcommands to group flags. There are no flags that can't be used together or must be used in a very specific combination.
+ **Well-documented**.
+ **Reliable**. Under the hood, enc uses [gopenpgp](https://github.com/ProtonMail/gopenpgp) library. The same library that powers ProtonMail.
+ **UNIX-way**. Enc does only one job and does it well. And it plays nicely with any other tools. It reads all possible input from stdin and writes all possible output into stdout.
+ **CI-friendly**. There is no interactive prompt. All input is strictly stdin or CLI flags.

A few drawbacks to keep in mind:

+ Not all encryption algorithms supported by GnuPG are supported by enc.
+ You'll still need to import keys into GnuPG to use tools that are integrated with GnuPG, like git.

## Install

If you have Go:

```bash
go install github.com/life4/enc@latest
```

If you don't have Go, [grab the binary for your OS](https://github.com/life4/enc/releases).

On Linux (and OS X, probably) that's how you can make the executable globally available:

1. Extract the binary: `tar -xf enc_*.tar.gz`
1. Make it executable: `chmod +x enc`
1. Place it in your PATH: `mv enc ~/.local/bin`
1. Check if it works: `env version`
1. If it says "command not found", run `echo $PATH` and check if `~/.local/bin` is there. If not, add into your `~/.bashrc` the following: `export PATH=$PATH:~/.local/bin`

## Encrypt

"To encrypt something" means making it unreadable for someone without a secret. Only the one who knows the secret can read an encrypted message. Let's encrypt a text message using a password:

```bash
echo 'my secret message' | enc encrypt --password 'very secret password' > encrypted.bin
```

## Decrypt

"To decrypt something" means to restore the encrypted message. If you look at the content of `encrypted.bin` from the previous step, you'll see that it's some binary gibberish. Let's decrypt that. And for that, you need to know the password that was used to encrypt the message.

```bash
cat encrypted.bin | enc decrypt --password 'very secret password'
```

And you should see the "my secret message" output. And if you pass an incorrect password, you'll see a "wrong password or malformed message" error instead.

## A note on secrets and shell history

It's not safe to just plainly put your passwords like this as an argument to a command. Or to use `echo` to write a secret message. All your input will be stored in the history of your terminal. For example, for bash, it will be saved in `~/.bash_history`. There are a few helpful tips on how to avoid that:

1. Start the command with a space. Then it will not be stored in the bash history. It should work for other shells as well.
1. Use [pass](https://www.passwordstore.org/) or another password manager: `enc encrypt --password=$(pass path/to/password)`.
1. Use `cat` without arguments as input: `enc encrypt --password=$(cat)`. It will read whatever you type in the terminal until you press `ctrl+d`.

## Armor/dearmor the message

Sometimes, you need to send the encrypted message as text, in a place where binary input isn't supported. For example, in a chat. For that, enc provides "armoring" that turns any binary input into text:

```bash
cat encrypted.bin | enc armor > encrypted.txt
```

Now, inside encrypted.txt you'll see something like this:

```text
-----BEGIN PGP MESSAGE-----
Version: enc 0.1.0
Comment: https://github.com/life4/enc

wy4ECQMIT0iy0Z6UgXHg6Zt9gwmLNWJ4Jx0aVE7K1CuFT03VoP7dmtAknap3+ioR
0kMB8dNyuHDE5mO27fu0GCJih60VSWcTbcFsSwanO8r462A0itZ68sDG5Tyv1b9C
y6LeJYJwgyGi8wemlqZVqdStggNM
=ArNH
-----END PGP MESSAGE-----
```

And to decrypt the armored message, you should dearmor it back into binary:

```bash
cat encrypted.txt | enc dearmor | enc decrypt --password 'very secret password'
```

**Tip**: you can omit `enc dearmor`. Enc will automatically detect if the input is armored and dearmor it.

## Generate a key

Passwords aren't that good for encrypting things. It's helpful when you want to send an encrypted file and then tell your friend the secret by phone (or shout it to him in the next room), but when you can get a bit fancier, it's better to use a secret key. A key is a file that can be used to encrypt or decrypt messages. It's longer (and so safer) than a typical password and has one more feature we'll cover later. For now, let's just generate a new key:

```bash
enc key generate > private.key
```

**Tip**: make sure to limit permissions for the keys you store locally (`chmod 600 *.key`).

The key has quite a bit of information inside: your name and email, when it was generated, and expiration date. Of course, you can have a look yourself:

```bash
cat private.key | enc key info
```

## Encrypt/decrypt with a key

Encrypting the message using the key is quite similar to encrypting it with a password. Just pass the path to the key to use:

```bash
echo 'hello world' | enc encrypt --key private.key > encrypted.bin
```

And similarly, decrypt:

```bash
cat encrypted.bin | enc decrypt --key private.key
```

## Use public key (generate and encrypt)

The "one more feature" of keys we mentioned before is that your private key actually contains 2 keys:

1. Public key is used to encrypt messages.
1. Private key is used to decrypt messages encrypted with the public key.

The idea is that you can make your public available for everyone on your website, chats, etc. Anyone can take that public key, use it to encrypt a message, and send the encrypted message to you. And despite the public key being public, nobody but you can decrypt the message. Neat!

Extract the public key from the private key:

```bash
cat private.key | enc key public > public.key
```

Encrypt the message with the public key:

```bash
echo 'hello world' | enc encrypt --key public.key > encrypted.bin
```

The message can be decrypted only using the private key:

```bash
$ cat encrypted.bin | enc decrypt --key private.key
hello world
$ cat encrypted.bin | enc decrypt --key public.key
Error: public key cannot be used to decrypt
```

**Tip**: keys can be armored using `enc key armor`.

## Protect private key with a password

If you use a private key to protect your files from evil hackers, the whole effort is in vain if the key lies in plain sight next to the files. It's like locking your door and then leaving the key in the keyhole. The solution is to encrypt ("lock") the private key itself with a password.

Lock the key with a password:

```bash
cat private.key | enc key lock --password 'my secret pass' > locked.key
```

You can always unlock it back if you change your mind:

```bash
cat locked.key | enc key unlock --password 'my secret pass' > unlocked.key
```

**Tip**: you can chain `enc key unlock` and `enc key lock` to change the password for the key. It's good to update your passwords time-to-time.

To use a locked key when using `encrypt` or `decrypt`, pass both `--key` and `--password` at the same time:

```bash
echo 'hello world' | enc encrypt --key locked.key --password 'my secret pass' > encrypted.bin
cat encrypted.bin | enc decrypt --key locked.key --password 'my secret pass'
```

## Sign

From the math perspective, there is no difference between private and public keys, they both can be used to encrypt messages that only can de be decrypted by the other. Most of the security tools, including enc, artificially forbid using the public key for decrypting messages because that's not how it should be used (encrypting messages that anyone can decrypt is pointless). But what if we bypass that limitation? Then we could calculate the hash from the message, encrypt it using our private key, and publish it alongside the message itself. Then anyone can take this "signature", decrypt it using the public key, and check if the hash matches the message. It will match only if the message is not altered by anyone and the signature was encrypted using your private key. In other words, anyone can validate that the message was sent by you and wasn't altered. This is what signing is.

Create a new signature:

```bash
cat encrypted.bin | enc sig create --key private.key > message.sig
```

**Tip**: signatures can be armored using `enc sig armor`.

The signature will contain the ID of the key that was used to generate it:

```bash
$ cat message.sig | enc sig id
91c1be98e13a8621
$ cat private.key | enc key info | jq .id
"91c1be98e13a8621"
```

## Verify signature

To verify the signature, you'll need the signed message, the signature, and the public key:

```bash
cat encrypted.bin | enc sig verify --key public.key --signature message.sig
```

## Download public key

Many services can host the public GPG keys of their users. And enc can search these services and download the key for you.

Supported providers:

1. `github`: get keys from [github.com](https://github.com/) by username.
1. `gitlab`: get keys from [gitlab.com](https://gitlab.com/) (or a self-hosted GitLab instance) by username.
1. `hkp`: get a key from a public GPG key server (by default, [keyserver.ubuntu.com](https://keyserver.ubuntu.com/)) by its fingerprint. Downloading keys by author's email is not supported by design. HKP servers do not verify user emails, so anyone can upload a key with any email address.
1. `keybase`: get keys from [keybase.io](https://keybase.io/) by username.
1. `protonmail`: get a key from [proton.me](https://proton.me/mail) by email address.

In the list above, "keys" means that the provider can return multiple keys, not just one.

Download a key of a proton mail user by their email:

```bash
enc remote get --provider=protonmail git@orsinium.dev
```

Search all providers and download a key by author's username:

```bash
enc remote get orsinium
```

## Publish public key

To publish a key in a supported provider, us the official tools provided by the provider:

+ Upload to github.com using [gh](https://cli.github.com/): `gh gpg-key add public.key`.
+ Upload to gitlab.com using [glab](https://gitlab.com/gitlab-org/cli): [not supported yet](https://gitlab.com/gitlab-org/cli/-/issues/1052).
+ Upload to keybase.io using [keybase](https://book.keybase.io/docs/cli): `keybase pgp import -i private.key`.

[![xkcd: Public Key](https://imgs.xkcd.com/comics/public_key_2x.png)](https://xkcd.com/1553/)

## Experimental: work with GnuPG keyring

Many great tools have integration with GnuPG. To name a few, git, some email clients, [pass](https://www.passwordstore.org/). Wouldn't it be great to integrate them with enc too? Well, that's not that easy. Many tools don't allow specifying a different path to GnuPG binary to use, so all we are left with is to integrate enc with GnuPG directly: import, export, and list keys. This is what this section is about. How you can work with GnuPG "keyring": the internal collection of keys that GnuPG knows about.

So far, we managed to only provide a few commands for public keys' keyring. The private keyring is a bit trickier, different versions of GnuPG store it differently.

List all keys that GnuPG knows about:

```bash
cat ~/.gnupg/pubring.gpg | enc keys list
```

Red keys are expired or revoked, green keys are locked (password-protected), and yellow keys aren't locked.

Get a key from the list (by ID or email):

```bash
cat ~/.gnupg/pubring.gpg | enc keys get 0123456789abcdef > public.key
cat ~/.gnupg/pubring.gpg | enc keys get mail@example.com > public.key
```

Add a key into the GnuPG keyring:

```bash
gpg --import private.key
```

## Type commands faster

1. Under the hood, enc uses [cobra](https://github.com/spf13/cobra) Go library for describing CLI. And [cobra provides shell completion support](https://github.com/spf13/cobra/blob/main/shell_completions.md). If you run `enc completion bash -h` (or another shell name you use instead of `bash`), it will show you how you can activate autocomplete for your shell depending on your OS.
1. Every command provides multiple aliases and shortcuts. For example, `enc key generate` can be abbreviated to `enc k g`. You can call the command with `-h` (`enc key generate -h`) to see what aliases it has.
1. Most of the flags can also be abbreviated to the first letter. For example, you can use `-p` instead of `--password` in all commands.
