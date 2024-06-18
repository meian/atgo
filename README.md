# atgo

`atgo` is a tool for Gophers that allows them to perform AtCoder operations using commands.

- [日本語](./README-ja.md)

## Conformance to rules for generated AI

[AtCoder生成AI対策ルール - 20240607版](https://info.atcoder.jp/entry/llm-abc-rules-ja)  
(mearning of `AtCoder AI Policy`)

AtCoder has established the rules described in the above document to restrict the use of generated AI in the ABC.  
`atgo` has been modified to conform to these rules in v0.0.4, so please use v0.0.4 or later to participate in ABC.

## Installation

`atgo` is installed with the following command.

```
curl -sSfL https://raw.githubusercontent.com/meian/atgo/main/install | bash
```

To use specify the build version, use the following command.

```
curl -sSfL https://raw.githubusercontent.com/meian/atgo/main/install | bash -s -- --tag v0.0.1
```

If a pre-built binary is available at [release](https://github.com/meian/atgo/releases), download the binary.  
For OS/architectures that do not have pre-built binaries, the installer will build them internally using `go install`, which requires Go 1.22 or higher.

## Operating Environment

- Linux
  - Development is working on Debian in Docker

Mac and Windows have not been tested, but they are expected to work if installation is possible.

## Usage

### Setting Up Your Workspace

The current directory will be used as the workspace.  
All commands should be run from this directory.

### Authentication

Register your credential with the following command using your AtCoder username and password:

```bash
$ atgo auth
Username: kitamin
Password: *******
```

Once you save your credential with `atgo auth`, you will be automatically logged in when executing other commands.

The credential is stored locally with simple encryption.  
If you want to remove the local credential, run the following command:

```bash
$ atgo auth clear
```

### Loading Tasks

To output the contest information and the task list of associated with it, run the following command:

```bash
$ atgo contest [Contest ID]
````

To output task details, run the following command:

```bash
$ atgo task [task ID]
````

If you want to prepare the files locally to the task, run the following command:

```bash
$ atgo task local-init [task ID]
````

The above command will create the following new files.

- `main.go`
- `main_test.go`
- `go.mod`
- `go.sum`

Once files are created, they are cached, and the next time you run `atgo task local-init` for the same task, the cached file will be loaded locally instead of being created anew.

### Implementing the answer and Testing

Please implement the code to solve the task in `main.go`.  
And you can verify the sample input and output result using `go test`.

### Submission

To submit your code, run the following command:  
The target for submission is the task that was last initialized with `atgo task local-init`.

```bash
$ atgo submit
````

### Additional Commands

To output the list of past contests:

```bash
# load past archived contests information(by category)
$ atgo contest load [abc or arc or agc or ahc]

# output list of past archived contests
$ atgo contest list [abc or arc or agc or ahc]
```

To clean up managed database:

```bash
$ atgo workspace clean
```

DB deleted by the above command contains the following information.  
Caches prepared locally with `atgo task local-init` are not removed.

- contest information
- task information
- contest and task associations

## License

This project is licensed under the MIT License.
