# cli-tasks-csv

This project was done in trying to get better at Golang. The specifications for the project can be found at [goprojects/01-todo-list](https://github.com/dreamsofcode-io/goprojects/tree/main/01-todo-list)

## How to build
> [!NOTE]
> You need to have Golang installed. Any version above 1.16 will do.

In the terminal of the cloned project, run the following to build the executable
```bash
go build . -o tasks
```

Now you can run `./tasks help` to see how to use the CLI app.

## What's Left?
- [ ] Tests
- [ ] List only uncompleted tasks unless `--all` flag is passed
- [ ] Ask user if they want to re-update completion date
- [ ] Coloured outputs on the terminal emulator