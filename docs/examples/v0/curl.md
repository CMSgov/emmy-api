# Using `curl` with the Emmy API

`curl` is a command-line style tool used to perform HTTP requests. It's powerful and flexible and requires very little in the way of installation.

If you need help installing `curl`, view the [Installing `curl`](curl.md#installing-curl) section below.

## Installing `curl`

### On a Mac or Linux Distro

Assuming you have already [installed Homebrew](https://brew.sh/), you can install `curl` from a Terminal with the command:

```bash
brew install curl
```

Be sure to note the instructions after installation about adding `curl` to your path. Namely:

```bash
echo 'export PATH="/opt/homebrew/opt/curl/bin:$PATH"' >> ~/.zshrc
```

### In a Windows Environment

Windows environments will generally have `curl` available in the latest versions. You should be able to open a command prompt and use it out of the box.

Note that in Windows, you will need to use double quotes (") to specify paramters, rather than single quote as in MacOS or Linux.
