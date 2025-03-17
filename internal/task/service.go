package task

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/n17ali/gohive/api/taskpb"
	"github.com/n17ali/gohive/internal/logs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskServiceServer struct {
	taskpb.UnimplementedTaskServiceServer
	logger     logs.Logger
	repository *TaskRepository
}

func NewTaskService(logger logs.Logger, repo *TaskRepository) *TaskServiceServer {
	return &TaskServiceServer{
		logger:     logger,
		repository: repo,
	}
}

func (s *TaskServiceServer) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskRresponse, error) {
	taskID := uuid.New().String()

	task := Task{
		ID:          taskID,
		Title:       req.Task.Title,
		Description: req.Task.Description,
		Interval:    time.Duration(req.Task.Interval),
	}

	if err := s.repository.SaveTask(ctx, task); err != nil {
		return nil, err
	}

	return &taskpb.CreateTaskRresponse{Id: taskID}, nil
}

func (s *TaskServiceServer) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetTaskRresponse, error) {
	task, err := s.repository.GetTask(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &taskpb.GetTaskRresponse{
		Task: &taskpb.Task{Id: task.ID,
			Title:       task.Title,
			Description: task.Description,
			Interval:    int64(task.Interval),
		}}, nil
}

func (s *TaskServiceServer) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskRresponse, error) {
	updatedTask := Task{
		ID:          req.Id,
		Title:       req.Task.Title,
		Description: req.Task.Description,
		Interval:    time.Duration(req.Task.Interval),
	}

	err := s.repository.UpdateTask(ctx, req.Id, updatedTask)
	if err != nil {
		return &taskpb.UpdateTaskRresponse{Success: false}, err
	}

	return &taskpb.UpdateTaskRresponse{Success: true}, nil
}

func (s *TaskServiceServer) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	if err := s.repository.DeleteTask(ctx, req.Id); err != nil {
		return &taskpb.DeleteTaskResponse{Success: false}, err
	}

	return &taskpb.DeleteTaskResponse{Success: true}, nil
}

func (s *TaskServiceServer) GetTaskLogs(ctx context.Context, req *taskpb.GetTaskLogsRequest) (*taskpb.GetTaskLogsResponse, error) {
	logEntries, err := s.logger.GetTaskLogs(ctx, req.Id, req.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch logs: %v", err)
	}

	var resp taskpb.GetTaskLogsResponse
	for _, entry := range logEntries {
		resp.Logs = append(resp.Logs, &taskpb.TaskLog{
			Time:    entry["time"].(string),
			Status:  entry["status"].(string),
			Message: entry["message"].(string),
		})
	}

	return &resp, nil
}
