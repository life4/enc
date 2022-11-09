# enc

Enc is a **modern and user-friendly alternative to gnupg**. It is easy to use, secure by default, and can encrypt and decrypt files using password or encryption keys, manage keys, and sign data. Our goal was to make encryption available to all engineers without the need to learn a lot of new words, concepts, and commands. It is the most beginner-friendly CLI tool for encryption, and keeping it that way is our top priority.

Features:

+ **Easy installation**. Grab the binary, and you're ready to go.
+ **Friendly CLI**. We use well isolated subcommands to group together flags. There are no flags that can't be used together or must be used in a very specific combination.
+ **Well-documented**.
+ **Reliable**. Under the hood, enc uses [gopenpgp](github.com/ProtonMail/gopenpgp) library. The same library that powers ProtonMail.
+ **UNIX-way**. Enc does only on job and does it well. And it plays nicely with any other tools. It reads all possible input from stdin and writes all possible output into stdout.
+ **CI-friendly**. There is no interactive prompt. All input is strictly stdin or CLI flags.

A few drawbacks to keep in mind:

+ Not all encryption algorithms supported by gnupg are supported by enc.
+ You'll still need import keys into gnupg to use tools that are integrated with gnupg, like git.

## Installation

If you have Go:

```bash
go install github.com/life4/enc@latest
```

If you don't have Go, grab the binary for your OS and put it anywhere in your PATH.

## Encrypt

"To encrypt something" means making it unreadable for someone without a secret. Only who knows the secret can read an encrypted message. Let's encrypt a text message using a password:

```bash
echo 'my secret message' | enc encrypt --password 'very secret password' > encrypted.bin
```

## Decrypt

"To decrypt something" means to restore the encrypted message. If you look at the content of `encrypted.bin` from the previous step, you'll see that it's some binary gibberish. Let's decrypt that. And for that, you need to know the password that was used to encrypt the message.

```bash
cat encrypted.bin | enc decrypt --password 'very secret password'
```

And you should see "my secret message" output. And if you pass an incorrect password, you'll see "wrong password or malformed message" error instead.

## A note on secrets and shell history

It's not safe to just plainly put your passwords like this as an argument to a command. Or use `echo` to write a secret message. All you input will be stored in the history of your terminal. For example, for bash it will be saved in `~/.bash_history`. There are a few helpful tips on how to avoid that:

1. Start the command with a space. Then it will not be stored in the bash history. It should work for other shells as well.
1. Use pass or another password manager: `enc encrypt --password=$(pass path/to/password)`.
1. Use `cat` without arguments as input: `enc encrypt --password=$(cat)`. It will read whatever you type in the terminal until you press `ctrl+d`.

## Armor/dearmor the message

Sometimes, you need to send the encrypted message as a text, in a place where binary input isn't supported. For example, in a chat. For that, enc provides "armoring" that turns any binary input into text:

```bash
cat encrypted.bin | enc armor > encrypted.txt
```

Now, inside of encrypted.txt you'll see something like this:

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

## Generate a key

Passwords aren't that good for encrypting things. It's helpful when you want to send an encrypted file and then tell your friend the secret by phone (or shout it to him in the next room), but when you can get a bit more fancy, it's better to use a secret key. A key is a file that can be used to encrypt or decrypt messages. It's longer (and so safer) than a typical password and has one more feature we'll cover later. For now, let's just generate a new key:

```bash
enc key generate > private.key
```

The key has quite a bit of information inside: your name and email, when it was generated, and expiration date. Of course, you can have a look yourself:

```bash
cat private.key | ./enc key info
```

## Encrypt with a key

Encrypting the message using the key is quite similar to encrypting it with a password. Just pass the path to the key to use:

```bash
echo 'hello world' | ./enc encrypt --key private.key > encrypted.bin
```

## Decrypt with a key

The "one more feature" of keys we mentioned before is that your private key actually contains 2 keys. Either of them can be used to decrypt what the other has encrypted. The full private key must be known only to you. And one of the two parts of the private key, called the "public key", is what you'll share with your friends. Your friends can use that key to decrypt messages you send to them or encrypt messages for you. And what's cool, nobody, even other friends with the same public key, can decrypt what your friends send to you. And nobody, even your friends, can generate a message that others will think is ncrypted by you. Pretty cool, huh?

Extract public key from private key:

```bash
cat private.key | ./enc key public > public.key
```

Decrypt the message:

```bash
cat encrypted.bin | ./enc decrypt --key public.key
```

## Sign

...

## Verify signature

...
