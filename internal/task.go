package pkg

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
)

/* type Task struct {
	ID          int
	Title       string
	CreatedAt   string
	CompletedAt string
} */

func Execute(filename string) {
	if len(os.Args) < 2 {
		ShowHelp()
		return
	}

	cmd := os.Args[1:2][0]
	switch cmd {

	case "add":
		Add(os.Args, filename)

	case "list":
		List(os.Args, filename)

	case "complete":
		SetAsCompleted(os.Args, filename)

	case "delete":
		Delete(os.Args, filename)

	case "help":
		ShowHelp()

	default:
		ShowHelp()
	}
}

func ShowHelp() {
	fmt.Println("How to use:")
	wr := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(wr, "tasks add <task title>\tadd a task\ttasks add 'Have fun'")
	fmt.Fprintln(wr, "tasks list\tview list of tasks\t")
	fmt.Fprintln(wr, "tasks complete <taskId>\tset task <taskId> as completed\ttasks complete 1")
	fmt.Fprintln(wr, "tasks delete <taskId>\tdelete task <taskId>\ttasks delete 1")
	wr.Flush()
}

func Add(args []string, filepath string) {
	if len(args) < 3 {
		fmt.Println("error: please provide task title")
		os.Exit(1)
	}

	f := LoadFile(filepath)
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	title := args[2]
	wr := csv.NewWriter(f)
	var newRecordId string

	if len(records) < 2 {
		newRecordId = "1"
		wr.Write([]string{"ID", "Title", "CreatedAt", "CompletedAt"})
	} else {
		i, _ := strconv.Atoi(records[len(records)-1][0])
		newRecordId = strconv.Itoa(i + 1)
	}
	wr.Write([]string{newRecordId, title, time.Now().Format(time.DateTime), ""})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer wr.Flush()

	fmt.Println("✅ Added:", title)
	CloseFile(f)
}

func List(args []string, filepath string) {
	f := LoadFile(filepath)
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if len(records) < 2 {
		fmt.Println("No tasks in list yet")
		return
	}

	wr := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(wr, "ID\tTitle\tCreated\tCompleted")
	for _, r := range records[1:] {
		created, _ := time.Parse(time.DateTime, r[2])
		completed, _ := time.Parse(time.DateTime, r[3])
		completedDiff := "nil"
		if !completed.IsZero() {
			completedDiff = timediff.TimeDiff(completed)
		}

		fmt.Fprintf(
			wr,
			"%s\t%s\t%s\t%s\n",
			r[0], r[1], timediff.TimeDiff(created), completedDiff,
		)
	}
	wr.Flush()
	CloseFile(f)
}

func SetAsCompleted(args []string, filepath string) {
	if len(args) < 3 {
		fmt.Println("error: please provide task ID")
		os.Exit(1)
	}
	f := LoadFile(filepath)
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
		if r[0] == args[2] {
			r[3] = time.Now().Format(time.DateTime)
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

	f = LoadFile(filepath)
	wr := csv.NewWriter(f)
	err = wr.WriteAll(new)
	if err != nil {
		log.Fatal(err)
	}

	wr.Flush()
	fmt.Println("✅ Set task:", args[2], "as completed")
	CloseFile(f)
}

func Delete(args []string, filepath string) {
	if len(args) < 3 {
		fmt.Println("error: please provide task ID")
		os.Exit(1)
	}

	f := LoadFile(filepath)
	old, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if len(old) < 2 {
		fmt.Println("No tasks in list yet")
	}

	headers := old[0]
	var new [][]string
	new = append(new, headers)

	for _, r := range old[1:] {
		if r[0] != args[2] {
			new = append(new, r)
		}
	}
	f.Close()

	err = os.Remove(f.Name())
	if err != nil {
		log.Fatal(err)
	}

	f = LoadFile(filepath)
	wr := csv.NewWriter(f)
	err = wr.WriteAll(new)
	if err != nil {
		log.Fatal(err)
	}

	wr.Flush()
	fmt.Println("✅ Deleted task", args[2])
	CloseFile(f)
}
