package main

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int
	Title       string
	CreatedAt   string
	CompletedAt int // REFACTOR make this a date somehow
	// DueAt       int // REFACTOR make this a date somehow
}

func main() {
	fileName := "./tasks.csv"
	createStoreIfNotExist(fileName)

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
		add(os.Args[2])

	case "list":
		list(fileName)

	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("error: please provide task ID")
			os.Exit(1)
		}
		setAsCompleted(os.Args[2], fileName)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("error: please provide task ID")
			os.Exit(1)
		}
		delete(os.Args[2], fileName)

	default:
		showHelp()
	}
}

func createStoreIfNotExist(filepath string) fs.FileInfo {
	file, _ := os.Stat(filepath)
	if file != nil {
		return file
	}

	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal("Unable to create new file at", filepath, err)
	}
	defer f.Close()

	wr := csv.NewWriter(f)
	headers := []string{
		"ID",
		"Title",
		"CreatedAt",
		"CompletedAt",
	}
	err = wr.Write(headers)
	if err != nil || wr.Error() != nil {
		log.Fatal(err)
	}
	defer wr.Flush()

	file, err = os.Stat(filepath)
	if err != nil {
		log.Fatal("Unable to fetch file at", filepath, err)
	}

	return file
}

func showHelp() {
	// TODO probably want to use tabwriter to make it neater
	fmt.Println(`How to use:
tasks add <task title> - add a task
tasks list - view list of tasks
tasks complete <taskId> - set task <taskId> as completed
tasks delete <taskId> - delete task <taskId>`)
}

func add(title string) {
	f, err := os.OpenFile("tasks.csv", os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer f.Close()

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
	}
	defer wr.Flush()

	fmt.Println("✅ Added:", title)
}

func list(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

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

func setAsCompleted(ID, filepath string) {
	f, err := os.OpenFile(filepath, os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

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

	f, err = os.Create(filepath)
	if err != nil {
		log.Fatal("Unable to create new file at", filepath, err)
	}
	defer f.Close()

	wr := csv.NewWriter(f)
	err = wr.WriteAll(new)
	if err != nil {
		log.Fatal(err)
	}

	wr.Flush()
	// fmt.Println("✅ Set task:", title, "as completed")
}

func delete(ID, filepath string) {
	f, err := os.OpenFile(filepath, os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

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

	f, err = os.Create(filepath)
	if err != nil {
		log.Fatal("Unable to create new file at", filepath, err)
	}
	defer f.Close()

	wr := csv.NewWriter(f)
	err = wr.WriteAll(new)
	if err != nil {
		log.Fatal(err)
	}

	wr.Flush()
	fmt.Println("✅ Deleted task", ID)
}
