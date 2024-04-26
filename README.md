# atgo

`atgo` is a tool for Gophers that allows them to perform AtCoder operations using commands.  

## Installation

To install the latest version of `atgo`, run the following command:

```bash
go install github.com/yourusername/atgo@latest
```

## Operating Environment

- Linux
  - Development is working on Debian in Docker

Mac and Windows have not been tested, but are expected to work if they can be installed.

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

Once you save credential with `atgo auth`, you will be logged in when executing other commands.

Authentication information is stored locally with simple encryption.  
If you want to remove the local credentials, run the following command:

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

Once files created, their is cached, and the next time you run `atgo task local-init` for the same task, the cached file will be loaded locally instead of being created anew.

### Implementing the answer and Testing

Please implement the code that submit the task on `main.go`.  
And you can verify the sample input and output result using `go test`.

### Submission

To submit your code, run the following command:  
The target of the submission is the task that last executed `atgo task local-init`.

```bash
$ atgo submission
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

- contest informations
- task informations
- contest and task associations

## License

This project is licensed under the MIT License.
