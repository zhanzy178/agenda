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
  logout      Command logout your current account.
  meeting     Command list meeting table you specific during time interval.
  moveout     Command move out user from meeting participators.
  quit        Command to quit a meeting you participated in.
  register    Command register your account for agenda system.
  state       Command to list your current user state.
  user        Command to list all user informations in system.

Flags:
      --config string   config file (default is $HOME/.agenda.yaml)
  -h, --help            help for agenda

Use "agenda [command] --help" for more information about a command.

```



### agenda register

​	The `agenda register` command will help you to create a account in agenda system. You need to setting user information(like username, password, email,phone number) with command flag, or input late. Once you create a user, your account informations will be store in `user.json` file under `agendaDataRoot` directory. Please try `agenda register help` to get more helping.

​	**Note**: You username must be unique in agenda system, and the first letter must be capitalized, or command will throw error.



### agenda login

​	The `agenda login` command provides a way to maintain your account login state under current shell (It works well under linux system). If you leave your shell or create a new shell, then you will lost your login state.

​	**Note**: If you login in an account logined, the account will be forced to logout on other shell.



### agenda logout

​	The `agenda logout` command provides a way to logout your current account under current shell. If you leave your shell or create a new shell, then you will lost you login status.

​	Once you logout your current accout and stay in an unlogined state, then what commands you can use are only `agenda login` and `agenda register`



### agenda state

​	The `agenda state` command will list your current user state, which means you should  login before use this command.



### agenda user

​	The `agenda user` command will list all user informations, such as Id, Name, E-mail, Number, State and Last-Log time,   in system with a table .

```bash
------------------------------------ User Table -----------------------------------------
Id   |Name     |E-mail                  |Number        |State     |Last-Log     
-----------------------------------------------------------------------------------------
1    |Caiye    |test@mail2.sysu.edu.cn  |159xxxxxxxx   |Online    |2018-11-03 17:32:39 Sat
-----------------------------------------------------------------------------------------
```



### agenda delete

​	The `agenda delete` command will delete your current account. So be carefule to use it. You have a chance to comfirm during the operation.



### agenda cm

​	The `agenda cm` command will create a meeting with meeting's Title, Start time and End time and particapators. 

​	Note : The first letter of title must be capitalized, or command will throw error. Also, the participator should not be empty, otherwise the meeting will not be created successfully.



### agenda cancel

​	The `agenda cancel` command will cancel the certain meeting you have created with a certain meeting title. You have a chance to comfirm during the operation.



### agenda join

​	The `agenda join` command will join other user to a certain meeting you initate. Make sure that the user you invite exist. 



### agenda meeting

​	The `agenda meeting` command will list meeting table you specific during time interval.

​	**Note** :  When input time, make sure that you are using the English form instead of Chinese form. 

```bash
-------------------------------------- Meeting Table -------------------------------------------
Id   |Title   |Sponsor   |Since                     |To                        |Participators  
------------------------------------------------------------------------------------------------
1    |Caiye   |Caiye     |2018-11-03 20:00:00 Sat   |2018-11-03 21:00:00 Sat   |Yvonne,        
------------------------------------------------------------------------------------------------
```




### agenda moveout

​	The `agenda moveout` move out user from meeting participators. Make sure that the user you move out from meeting is in the participators list. 



### agenda quit

​	The `agenda quit` command quit a meeting you participated in.  Make sure that you are in the participators list instead of the sponsor.



### agenda clear

​	The `agenda clear` command will cancel all meetings. 

​	Note : All meetings will be clear, no matter whether you are the sponsor or participator of the meeting or nor, even the meeting has nothing to do with you.