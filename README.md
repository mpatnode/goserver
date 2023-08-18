## Intro
I took advantage of this project to teach myself Go.    The simplicity of setting up an http server and router paid off here, though it certainly slowed me down as I learned the language semantics.   I'm sure I've violated a couple Go best practices here and there, but it made the task much more enjoyable.  This was also my first use of SQLite.   Couple reasons; one to simply make the assignment more interesting and the other to avoid any database bootstrap setup in the deployment.  Furthermore with 1B+ users, it's probably worth knowing something about. Obviously, this wouldn't make a great production server architecture, but greatly simplified the task at a cost of learning a slightly new SQL api.

## Instructions

1. Unpack the zip file
1. Run "docker-compose up -d"
1. Run "docker-compose exec app go test -count=1 ./tests" to run the tests
1. Connect to localhost:8080/ to see API summary
1. Connect to http://localhost:8080/0.9/user/tkid@gmail.com to see the user created by the test

Note the tests and code are currently hardcoded to port 8080

## Summary Questions

### What were the most difficult tasks
- Figuring out that delve wouldn't run on WSL without upgrading the OS & Linux Kernel (syscall mapping vs real VM):  https://github.com/microsoft/vscode-go/pull/3167 but this also added docker support which was handy.
- braintree & SQLite go API docs were a little confusing.   Examples helped.
- Remembering not to type trailing semicolons and () in if statements

### Did you learn anything new while completing this assignment
- Go (I did write Towers of Hanio before starting this project)
    - Nice HTTP & router package
    - Cool JSON features
    - Was disappointed at the lack of references (at least for function argurments) weren't support, but the flexible . notation made up for that, as long as you're not trying to index into a slice.   (*slice)[] => :-(
    - strings.Index() is haystack, needle :sigh:
    - You can return a pointer to an item on the stack!
- SQLite parameter substitution handles the quotes around string values for you

### What did you not have time to add?
- Command line/config/env control of HTTP port number
- Secrets managment
- More/better unit tests
- Better namespace management via interfaces/packages
- User update

### What work took the up majority of your time?
Understanding the SQLite API in the Go world, and how it related to tables and column types.   The fun of mixing a new language with a new DB API.  A couple of Visual Studio/Go integration quirks tripped me up

### How could the API URI be improved?
- Allow for email key as well as ID (did that)
- Add PUT for Update