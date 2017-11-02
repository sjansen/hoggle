# Hoggle: Standalone Custom Transfer Agent for Git LFS

Hoggle enables Git LFS to store objects in S3 without server-side
support.

[![Go Report Card](https://goreportcard.com/badge/github.com/sjansen/hoggle)](https://goreportcard.com/report/github.com/sjansen/hoggle)

## Use Cases

### AWS CodeCommit

As of September 2017, CodeCommit does not have integrated support
for Git LFS.  While CodeCommit is able to handle very large repos,
Git itself is not good at handling large files.  Hoggle makes it
possible to combine Git & CodeCommit for source code, with Git LFS
& S3 for binary assets.

### Private Git Servers

Hoggle's minimal requirements make it a convenient solution for
private repositories. It is not a good choice for public repositories.

## Setup

1) Make sure your AWS credentials are properly configured.

   Hoggle loads your AWS credentials from the
   standard locations:
   * [environment variables](http://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html)
   * [config files](http://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html)
   * [instance metadata](http://docs.aws.amazon.com/cli/latest/userguide/cli-metadata.html)

2) Install [Git LFS](https://github.com/git-lfs/git-lfs/wiki/Installation).

3) Install Hoggle:
    ```
    $ go get github.com/sjansen/hoggle
    ```

4) Configure your local git repository to use Hoggle:
    ```
    $ hoggle init s3://bucket-name/prefix/
    ```
    Alternatively:
    ```
    $ git config lfs.customtransfer.hoggle.path hoggle
    $ git config lfs.customtransfer.hoggle.args s3://bucket-name/prefix/
    $ git config lfs.customtransfer.hoggle.concurrent false
    $ git config lfs.standalonetransferagent hoggle
    ```

5) If the project already uses Hoggle, download any files:
    ```
    $ git lfs fetch
    ```

6) If the project doesn't use Hoggle yet, configure Git LFS to manage large files:
    ```
    $ git lfs track 'assets/**'
    ```

## Status

Hoggle is only appropriate for early adopters interested in providing
feedback and patches. It has not had enough real world usage to ensure
it is robust. That said, I expect it to be reliable enough for personal
projects or use by small teams.

### Tested with:
 - Git LFS 2.3.0
 - Go 1.9

### Current Limitations

#### Concurrent Transfer

Hoggle has not been tested with concurrent transfer enabled. This is a
temporary situation.

#### Newly Cloned Repos

Cloning a repo that uses Hoggle results in the error:
`Clone succeeded, but checkout failed.` More specifically:

```
Receiving objects: 100% (29/29), 39.42 KiB | 858.00 KiB/s, done.
Downloading assets/icons/app-icon.png (3.0 KB)
Error downloading object: assets/icons/app-icon.png (d64a5b9): Smudge error: Error downloading assets/icons/app-icon.png
(snip)
error: external filter 'git-lfs filter-process' failed
fatal: assets/icons/app-icon.png: smudge filter lfs failed
warning: Clone succeeded, but checkout failed.
You can inspect what was checked out with 'git status'
and retry the checkout with 'git checkout -f HEAD'
```

This happens because standalone Git LFS tranfer agents must be configured
manually after the repository has been cloned, they can not be added to
`.lfsconfig` for security reasons. The work around it to clone the repo,
configure Hoggle, then finish checking out a branch.

```
$ git clone $REPO
$ hoggle init s3://bucket-name/prefix/
$ git checkout -f HEAD
```

I believe this limitation can be improved, but not completely eliminated
without reducing the safety of Git LFS.

## Roadmap

Hoggle is a personal project. I work on it when I feel like it.
The following is a list of things I might work on. It is not a
list of commitments.

Contributions are welcome.

### 0.2
 - improved logging and tracing
 - more robust error handling

### 0.3
 - more flexible configuration

### 0.4
 - support for additional cloud storage solutions
