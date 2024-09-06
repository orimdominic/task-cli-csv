package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"
)

type Task struct {
	ID          int
	Title       string
	CreatedAt   string
	CompletedAt string // REFACTOR make this a date somehow
	// DueAt       int // REFACTOR make this a date somehow
}

var filename = "./tasks.csv"

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	cmd := os.Args[1:2][0]
	switch cmd {
	case "help":
		showHelp()

	case "add":
		if len(os.Args) < 3 {
			fmt.Println("error: please provide task title")
			os.Exit(1)
		}
		f := loadFile(filename)
		add(os.Args[2], f)
		closeFile(f)

	case "list":
		f := loadFile(filename)
		list(f)
		closeFile(f)

	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("error: please provide task ID")
			os.Exit(1)
		}
		f := loadFile(filename)
		setAsCompleted(os.Args[2], f)
		closeFile(f)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("error: please provide task ID")
			os.Exit(1)
		}
		f := loadFile(filename)
		delete(os.Args[2], f)
		closeFile(f)

	default:
		showHelp()
	}

}

func loadFile(filepath string) *os.File {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
	if err != nil {
		f.Close()
		log.Fatal(err)
	}

	return f
}

func closeFile(f *os.File) {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	f.Close()
}

func showHelp() {
	// TODO probably want to use tabwriter to make it neater
	fmt.Println(`How to use:
tasks add <task title> - add a task
tasks list - view list of tasks
tasks complete <taskId> - set task <taskId> as completed
tasks delete <taskId> - delete task <taskId>`)
}

func add(title string, f *os.File) {
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	wr := csv.NewWriter(f)
	err = wr.Write([]string{
		strconv.Itoa(len(records)),
		title,
		time.Now().String(),
		"",
	})
	if err != nil {
		fmt.Println("error writing record to csv:", err)
		os.Exit(1)
	}
	defer wr.Flush()

	fmt.Println("✅ Added:", title)
}

func list(f *os.File) {
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	records = records[1:]

	if len(records) == 0 {
		fmt.Println("No tasks in task list")
		return
	}

	// TODO probably want to use tabwriter to make it neater
	fmt.Println("ID", "Title", "Date Created", "Date Completed")
	for _, r := range records {
		for _, t := range r {
			fmt.Print(t, " ")
		}
		fmt.Println("")
	}
}

func setAsCompleted(ID string, f *os.File) {
	old, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	headers := old[0]
	var new [][]string
	new = append(new, headers)

	for _, r := range old[1:] {
		// TODO use stdin to ask user if to update already-updated value
		// when CompletedAt is set already
		if r[0] == ID {
			r[3] = time.Now().String()
			new = append(new, r)
			continue
		}
		new = append(new, r)
	}
	f.Close()

	err = os.Remove(f.Name())
	if err != nil {
		log.Fatal(err)
	}

	f = loadFile(filename)

	wr := csv.NewWriter(f)
	err = wr.WriteAll(new)
	if err != nil {
		log.Fatal(err)
	}

	wr.Flush()
	// fmt.Println("✅ Set task:", title, "as completed")
}

func delete(ID string, f *os.File) {
	old, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	headers := old[0]
	var new [][]string
	new = append(new, headers)

	for _, r := range old[1:] {
		if r[0] != ID {
			new = append(new, r)
		}
	}
	f.Close()

	err = os.Remove(f.Name())
	if err != nil {
		log.Fatal(err)
	}

	f = loadFile(filename)

	wr := csv.NewWriter(f)
	err = wr.WriteAll(new)
	if err != nil {
		log.Fatal(err)
	}

	wr.Flush()
	fmt.Println("✅ Deleted task", ID)
}
