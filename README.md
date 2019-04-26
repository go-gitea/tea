# Gitea Command Line Tool for Go

This project acts as a command line tool for operating one or multiple Gitea instances. It depends on [code.gitea.io/sdk](https://code.gitea.io/sdk) client SDK implementation written in Go to interact with
the Gitea API implementation.

## Installation

Currently no prebuilt binaries are provided.
To install, a Go installation is needed.

```sh
go get code.gitea.io/tea
go install code.gitea.io/tea
```

If the `tea` executable is not found, you might need to set up your `$GOPATH` and `$PATH` variables first:

```sh
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

## Usage

First of all, you have to create a token on your `personal settings -> application` page of your gitea instance.
Use this token to login with `tea`:

```sh
tea login add --name=try --url=https://try.gitea.io --token=xxxxxx
```

Now you can use the `tea` commands:

```sh
tea issues
tea releases
```

> If you are inside a git repository hosted on a gitea instance, you don't need to specify the `--login` and `--repo` flags!

## Compilation

To compile the sources yourself run the following:

```sh
go get code.gitea.io/tea
cd "${GOPATH}/src/code.gitea.io/tea"
go build
```

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

* [Maintainers](https://github.com/orgs/go-gitea/people)
* [Contributors](https://github.com/go-gitea/tea/graphs/contributors)

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the
full license text.
