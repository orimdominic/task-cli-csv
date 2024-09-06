package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"
)

type Task struct {
	ID          int
	Title       string
	CreatedAt   string
	CompletedAt string
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
	fmt.Println("How to use:")
	wr := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(wr, "tasks add <task title>\tadd a task\ttasks add 'Have fun'")
	fmt.Fprintln(wr, "tasks list\tview list of tasks\t")
	fmt.Fprintln(wr, "tasks complete <taskId>\tset task <taskId> as completed\ttasks complete 1")
	fmt.Fprintln(wr, "tasks delete <taskId>\tdelete task <taskId>\ttasks delete 1")
	wr.Flush()
}

func add(title string, f *os.File) {
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	wr := csv.NewWriter(f)

	var newRecordId string

	if len(records) < 2 {
		newRecordId = "1"
		wr.Write([]string{"ID", "Title", "CreatedAt", "CompletedAt"})
	} else {
		i, _ := strconv.Atoi(records[len(records)-1][0])
		newRecordId = strconv.Itoa(i + 1)
	}
	wr.Write([]string{newRecordId, title, time.Now().String(), ""})

	if err != nil {
		fmt.Println(err)
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

	if len(records) < 2 {
		fmt.Println("No tasks in task list")
		return
	}

	wr := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(wr, "ID\tTitle\tCreated\tCompleted")
	for _, r := range records[1:] {
		fmt.Fprintf(wr, "%s\t%s\t%s\t%s\n", r[0], r[1], r[2], r[3])
	}
	wr.Flush()
}

func setAsCompleted(ID string, f *os.File) {
	old, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if len(old) < 2 {
		fmt.Println("No tasks in task list")
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
	fmt.Println("✅ Set task:", ID, "as completed")
}

func delete(ID string, f *os.File) {
	old, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if len(old) < 2 {
		fmt.Println("No tasks in task list")
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

// TODO
// Format date output to use time differences
//
