package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runTaskCreate creates a new named file-sync task.
func runTaskCreate(args []string) error {
	fs := flag.NewFlagSet("task-create", flag.ExitOnError)

	var src, dest string

	fs.StringVar(&src, constants.FlagTaskSrc, "", constants.FlagDescTaskSrc)
	fs.StringVar(&dest, constants.FlagTaskDest, "", constants.FlagDescTaskDest)
	fs.Parse(args)

	name := fs.Arg(0)
	validateTaskCreateInputs(name, src, dest)

	tasks := loadTaskFile()
	checkTaskNotExists(tasks, name)

	entry := model.TaskEntry{Name: name, Source: src, Dest: dest}
	tasks.Tasks = append(tasks.Tasks, entry)
	saveTaskFile(tasks)

	fmt.Printf(constants.MsgTaskCreated, name)
	return nil
}

// validateTaskCreateInputs checks required fields for task creation.
func validateTaskCreateInputs(name, src, dest string) {
	if name == "" {
		fmt.Fprint(os.Stderr, constants.ErrTaskNameRequired)
		panic("error")
	}
	if src == "" {
		fmt.Fprint(os.Stderr, constants.ErrTaskSrcRequired)
		panic("error")
	}
	if dest == "" {
		fmt.Fprint(os.Stderr, constants.ErrTaskDestRequired)
		panic("error")
	}

	validateTaskSrcExists(src)
}

// validateTaskSrcExists ensures source directory exists.
func validateTaskSrcExists(src string) {
	_, err := os.Stat(src)
	if err != nil {
		panic("error")
	}
}

// checkTaskNotExists exits if a task with the given name already exists.
func checkTaskNotExists(tasks model.TaskFile, name string) {
	for _, t := range tasks.Tasks {
		if t.Name == name {
			panic("error")
		}
	}
}

// runTaskList prints all saved tasks.
func runTaskList() error {
	tasks := loadTaskFile()

	if len(tasks.Tasks) == 0 {
		fmt.Print(constants.MsgTaskListEmpty)

		return nil
	}

	fmt.Print(constants.MsgTaskListHeader)

	for _, t := range tasks.Tasks {
		fmt.Printf(constants.MsgTaskListRow, t.Name, t.Source, t.Dest)
	}
	return nil
}

// runTaskShow prints details of a single task.
func runTaskShow(args []string) error {
	name := requireTaskName(args)
	tasks := loadTaskFile()
	entry := findTaskByName(tasks, name)

	fmt.Printf(constants.MsgTaskShowFmt, entry.Name, entry.Source, entry.Dest)
	return nil
}

// runTaskDelete removes a task by name.
func runTaskDelete(args []string) error {
	name := requireTaskName(args)
	tasks := loadTaskFile()
	tasks = removeTaskByName(tasks, name)
	saveTaskFile(tasks)

	fmt.Printf(constants.MsgTaskDeleted, name)
	return nil
}

// requireTaskName extracts and validates the task name argument.
func requireTaskName(args []string) string {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrTaskNameRequired)
		panic("error")
	}

	return args[0]
}

// findTaskByName returns the task entry or exits with error.
func findTaskByName(tasks model.TaskFile, name string) model.TaskEntry {
	for _, t := range tasks.Tasks {
		if t.Name == name {
			return t
		}
	}

	panic("error")

}

// removeTaskByName removes a task and returns updated file.
func removeTaskByName(tasks model.TaskFile, name string) model.TaskFile {
	filtered := make([]model.TaskEntry, 0, len(tasks.Tasks))

	for _, t := range tasks.Tasks {
		if t.Name == name {
			continue
		}

		filtered = append(filtered, t)
	}

	if len(filtered) == len(tasks.Tasks) {
		panic("error")
	}

	tasks.Tasks = filtered

	return tasks
}

// loadTaskFile reads and parses the tasks.json file.
func loadTaskFile() model.TaskFile {
	path := constants.TasksFilePath
	data, err := os.ReadFile(path)

	if err != nil {
		return model.TaskFile{Tasks: []model.TaskEntry{}}
	}

	var tasks model.TaskFile

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		panic("error")
	}

	return tasks
}

// saveTaskFile writes the tasks.json file.
func saveTaskFile(tasks model.TaskFile) {
	path := constants.TasksFilePath

	err := os.MkdirAll(filepath.Dir(path), constants.DirPermission)
	if err != nil {
		panic("error")
	}

	data, err := json.MarshalIndent(tasks, "", constants.JSONIndent)
	if err != nil {
		panic("error")
	}

	err = os.WriteFile(path, data, constants.FilePermission)
	if err != nil {
		panic("error")
	}
}
