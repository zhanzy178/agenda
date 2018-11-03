# Agenda Meeting Manager
​	Agenda is a CLI command tool, which will help team to manage their meetings on bash. It is the best way for you to cooperate with other user on shell.



## Installation

[Go installation](https://golang.org/doc/install) required!

Then, run the following command to install (It may be taken a while)

```bash
go get github.com/zhanzongyuan/agenda
```

Check install, then command will output helping for agenda system

```bash
agenda
```

You can setting disk file path `agendaDataRoot` to sync your system data in file `$HOME/.agenda.yaml`:

```yaml
agendaDataRoot: /path/to/your/home/.agenda/
```

Please read `.agenda.example.yaml` file for more system settings.



## Usage

Input `agenda` to read the helping information.

```bash
Agenda is an useful CLI program for everyone to manage meeting.

Usage:
  agenda [flags]
  agenda [command]

Available Commands:
  cancel      Command to cancel meeting you initial.
  clear       Command to cancel all meeting you initial.
  cm          Command to create meeting.
  delete      Command will delete your current account.
  help        Help about any command
  join        Command join other user to a certain meeting you initate.
  login       Command login your account on Agenda system.
  logout      Command logout your current account
  meeting     Command list meeting table you specific during time interval.
  moveout     Command move out user from meeting participators
  quit        Command to quit a meeting you participated in.
  register    Command register your account for agenda system.
  state       Command to list your current user state
  user        Command to list all user informations in system

Flags:
      --config string   config file (default is $HOME/.agenda.yaml)
  -h, --help            help for agenda

Use "agenda [command] --help" for more information about a command.

```



### agenda register

​	The `agenda register` command will help you to create a account in agenda system. You need to setting user information(like username, password, email..) with command flag, or input late. Once you create a user, your account informations will be store in `user.json` file under `agendaDataRoot` directory. Please try `agenda register help` to get more helping.

​	**Note**: You username must be unique in agenda system, or command will throw duplicate name error.



### agenda login

​	The `agenda login` command provides a way to maintain your account login state under current shell (It works well under linux system). If you leave your shell or create a new shell, then you will lost your login state.

​	**Note**: If you login in an account logined, the account will be forced to logout on other shell.



