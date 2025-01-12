package usecase

import (
	"errors"
	"lntvan166/togo/internal/domain"
	e "lntvan166/togo/internal/entities"
	"lntvan166/togo/pkg"
)

type taskUsecase struct {
	taskRepo domain.TaskRepository
	userRepo domain.UserRepository
}

func NewTaskUsecase(repo domain.TaskRepository, userRepo domain.UserRepository) *taskUsecase {
	return &taskUsecase{
		taskRepo: repo,
		userRepo: userRepo,
	}
}

func (t *taskUsecase) CreateTask(task *e.Task, username string) (int, int, error) {
	id, err := t.userRepo.GetUserIDByUsername(username)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "get user id failed")
		return 0, 0, errors.New("user does not exist")
	}

	isLimit, err := t.CheckLimitTaskToday(id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "check limit task today failed")
		return 0, 0, errors.New("check limit task today failed")
	}

	if isLimit {
		// pkg.ERROR(w, http.StatusBadRequest, fmt.Errorf("you have reached the limit of task today"), "")
		return 0, 0, errors.New("you have reached the limit of task today")
	}

	task.CreatedAt = pkg.GetCurrentTime()
	task.UserID = id

	taskID, err := t.taskRepo.CreateTask(task)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "add task failed")
		return 0, 0, errors.New("add task failed")
	}

	numberTask, err := t.taskRepo.GetNumberOfTaskTodayByUserID(id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "get number of task today failed")
		return 0, 0, errors.New("get number of task today failed")
	}
	return taskID, numberTask, nil
}

func (t *taskUsecase) GetAllTask() (*[]e.Task, error) {
	return t.taskRepo.GetAllTask()
}

func (t *taskUsecase) GetTaskByID(id int, username string) (*e.Task, error) {
	user_id, err := t.GetUserIDByTaskID(id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "task does not exist!")
		return nil, err
	}

	err = t.CheckAccessPermission(username, user_id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "check access permission failed: ")
		return nil, err
	}

	task, err := t.taskRepo.GetTaskByID(id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "get task by id failed!")
		return nil, errors.New("task does not exist")
	}
	return task, nil
}

func (t *taskUsecase) GetTasksByUsername(username string) (*[]e.Task, error) {
	id, err := t.userRepo.GetUserIDByUsername(username)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "get user id failed")
		return nil, errors.New("user does not exist")
	}

	tasks, err := t.taskRepo.GetTasksByUserID(id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "get tasks by user id failed")
		return nil, errors.New("get tasks by user id failed")
	}
	return tasks, nil
}

func (t *taskUsecase) GetUserIDByTaskID(id int) (int, error) {
	task, err := t.taskRepo.GetTaskByID(id)
	if err != nil {
		return 0, errors.New("task does not exist")
	}
	return task.UserID, nil
}

func (t *taskUsecase) CheckLimitTaskToday(id int) (bool, error) {
	user, err := t.userRepo.GetUserByID(id)
	if err != nil {
		return false, errors.New("user does not exist")
	}
	numberTask, err := t.taskRepo.GetNumberOfTaskTodayByUserID(id)
	if err != nil {
		return false, errors.New("get number of task today failed")
	}
	if numberTask >= int(user.MaxTodo) {
		return true, nil
	}
	return false, nil
}

func (t *taskUsecase) UpdateTask(id int, username string, newTask e.Task) error {
	user_id, err := t.GetUserIDByTaskID(id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "task does not exist!")
		return err
	}

	err = t.CheckAccessPermission(username, user_id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "check access permission failed: ")
		return err
	}

	_, err = t.taskRepo.GetTaskByID(id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "get task by id failed!")
		return errors.New("task does not exist")
	}

	err = t.taskRepo.UpdateTask(&newTask)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "update task failed!")
		return errors.New("update task failed")
	}
	return nil
}

func (t *taskUsecase) CompleteTask(id int, username string) error {
	user_id, err := t.GetUserIDByTaskID(id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "task does not exist!")
		return errors.New("task does not exist")
	}

	err = t.CheckAccessPermission(username, user_id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "check access permission failed: ")
		return err
	}

	err = t.taskRepo.CompleteTask(id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "complete task failed!")
		return errors.New("complete task failed")
	}
	return nil
}

func (t *taskUsecase) CheckAccessPermission(username string, taskUserID int) error {
	userID, err := t.userRepo.GetUserIDByUsername(username)
	if err != nil {
		return errors.New("user does not exist")
	}
	if userID != taskUserID {
		return errors.New("you are not allowed to access this task")
	}

	return nil
}

func (t *taskUsecase) DeleteTask(id int, username string) error {
	user_id, err := t.GetUserIDByTaskID(id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "task does not exist!")
		return err
	}

	err = t.CheckAccessPermission(username, user_id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "check access permission failed: ")
		return err
	}

	err = t.taskRepo.DeleteTask(id)
	if err != nil {
		// pkg.ERROR(w, http.StatusInternalServerError, err, "delete task failed!")
		return errors.New("delete task failed")
	}
	return nil
}

// for testing
func (t *taskUsecase) RollbackFromDelete(task *e.Task) error {
	err := t.taskRepo.RollbackFromDelete(task)
	if err != nil {
		return errors.New("rollback from delete task failed")
	}
	return nil
}
