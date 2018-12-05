# Gitea Command Line Tool for Go

This project acts as a command line tool for operating one or multiple Gitea instances. It depends on [code.gitea.io/sdk](https://code.gitea.io/sdk) client SDK implementation written in Go to interact with
the Gitea API implementation.

## Installation

```
go get github.com/go-gitea/tea
go install github.com/go-gitea/tea
```

## Usage

First of all, you have to create a token on your personal settings -> application.

```
git clone git@try.gitea.io:gitea/gitea.git
cd gitea
tea login add --name=try --url=https://try.gitea.io --token=xxxxxx
tea issues
tea releases
```

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

* [Maintainers](https://github.com/orgs/go-gitea/people)
* [Contributors](https://github.com/go-gitea/tea/graphs/contributors)

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the
full license text.
