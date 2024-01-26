# Version Control System Usage

## Note:-

-   This can only track files and not directories.
-   When the project is run for the first time, it creates `./vcs` and `./vcs/commits` directories, which stores necessary information about the version control system such as author, commit history, different file versions, etc.
-   The project already have `./vcs` directory and some files (`one.txt` & `two.txt`) for demonstration purposes. Feel free to delete them and start from the beginning.

## Let's start:-

To start using it, first of all build the binaries by executing `go build ./main.go` in the terminal. \
Then, run the `main.go` file using command-line (for Linux - `./main` & for Windows - `./main.exe`) by passing one of the following arguments:-

-   `config` : Gets the current username. `config example` : Sets the username to "example". The username is stored in `./vcs/config.txt`.
-   `add filename.txt` : Stages "filename.txt" for tracking/versioning. The staged/tracked files are stored in `./vcs/index.txt`. You can also provide the path of the file, if it's not in the project directory.
-   `log` : Logs the commit history to the console. The latest commit is displayed first.
-   `commit "message of commit"` : Commits the changes made in the tracked files. It uses `crypto/sha256` package to hash files and generate a commit hash. This hash is unique and changes when the contents of the tracked files changes. A directory with the name equal to the hash is created inside `./vcs/commits/` and the tracked files are copied to it. The commit is then logged to `./vcs/log.txt` with the unique hash, commit message and author.
-   `checkout commit_hash` : Restores the tracked files to a state corresponding to the commit. The contents of files inside `./vcs/commits/commit_hash` is copied to the original tracked files.

## Examples:-

```
./main config Ayush
```

```
./main add file1.txt
```

```
./main commit "file1.txt changed"
```

```
./main log
```

```
./main checkout 3a2d3879852d917ea9741c83e2f88d06970428a311511d452353b62bb06a4047
```
