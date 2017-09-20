1) Verify that all tests are passing.
1) Update `CHANGELOG.md` and commit.
1) Create release branch:
    ```
    git checkout -b release/v0.1
    ```
1) Update `version.go` and commit.
    * `0.1.0-dev` -> `0.1.0`
1) Tag release.
    ```
    git tag -a v0.1.0 -m "Release 0.1.0"
    ```
1) Build release binaries.
    ```
    make dist
    ```
1) Update `version.go` and commit.
    * `0.1.0` -> `0.1.1-dev`
1) Push commits and tags.
    ```
    git push origin release/v0.1
    git push origin v0.1.0
    ```
1) Upload release binaries to GitHub.
    * https://github.com/sjansen/hoggle/releases