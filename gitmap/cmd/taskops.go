package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
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
		err := apperror.NewWithDetails(
			"cmd.task.validateName",
			"E1092",
			constants.ErrTaskNameRequired,
			"cmd.task",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
	}
	if src == "" {
		err := apperror.NewWithDetails(
			"cmd.task.validateSrc",
			"E1093",
			constants.ErrTaskSrcRequired,
			"cmd.task",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
	}
	if dest == "" {
		err := apperror.NewWithDetails(
			"cmd.task.validateDest",
			"E1094",
			constants.ErrTaskDestRequired,
			"cmd.task",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
	}

	validateTaskSrcExists(src)
}

// validateTaskSrcExists ensures source directory exists.
func validateTaskSrcExists(src string) {
	_, err := os.Stat(src)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.task.validateSrcExists",
			"E1095",
			"source directory does not exist for task creation",
			"cmd.task",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"src": src},
		)
		cliexit.HandleError(appErr, 1)
	}
}

// checkTaskNotExists exits if a task with the given name already exists.
func checkTaskNotExists(tasks model.TaskFile, name string) {
	for _, t := range tasks.Tasks {
		if t.Name == name {
			err := apperror.NewWithDetails(
				"cmd.task.checkNotExists",
				"E1096",
				fmt.Sprintf("task '%s' already exists", name),
				"cmd.task",
				apperror.ErrorTypeValidation,
				apperror.SeverityError,
				map[string]any{"name": name},
			)
			cliexit.HandleError(err, 1)
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
	entry, err := findTaskByName(tasks, name)
	if err != nil {
		return err
	}

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
		err := apperror.NewWithDetails(
			"cmd.task.requireName",
			"E1097",
			constants.ErrTaskNameRequired,
			"cmd.task",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
	}

	return args[0]
}

// findTaskByName returns the task entry or exits with error.
func findTaskByName(tasks model.TaskFile, name string) (model.TaskEntry, error) {
	for _, t := range tasks.Tasks {
		if t.Name == name {
			return t, nil
		}
	}

	return model.TaskEntry{}, apperror.NewSimple("task not found: "+name, "E9023")
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
		err := apperror.NewWithDetails(
			"cmd.task.removeByName",
			"E1098",
			fmt.Sprintf("task '%s' not found for removal", name),
			"cmd.task",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"name": name},
		)
		cliexit.HandleError(err, 1)
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
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.task.loadJSON",
			"E1099",
			"failed to parse tasks JSON file",
			"cmd.task",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"path": path},
		)
		cliexit.HandleError(appErr, 1)
	}

	return tasks
}

// saveTaskFile writes the tasks.json file.
func saveTaskFile(tasks model.TaskFile) {
	path := constants.TasksFilePath

	err := os.MkdirAll(filepath.Dir(path), constants.DirPermission)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.task.saveMkdir",
			"E1100",
			"failed to create directory for tasks file",
			"cmd.task",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"path": path},
		)
		cliexit.HandleError(appErr, 1)
	}

	data, err := json.MarshalIndent(tasks, "", constants.JSONIndent)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.task.saveMarshal",
			"E1101",
			"failed to serialize tasks to JSON",
			"cmd.task",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(appErr, 1)
	}

	err = os.WriteFile(path, data, constants.FilePermission)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.task.saveWrite",
			"E1102",
			"failed to write tasks file",
			"cmd.task",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"path": path},
		)
		cliexit.HandleError(appErr, 1)
	}
}
