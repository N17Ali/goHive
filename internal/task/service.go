package task

import (
	"context"
	"time"

	"github.com/google/uuid"
	pb "github.com/n17ali/gohive/internal/pb"
)

type TaskServiceServer struct {
	pb.UnimplementedTaskServiceServer
}

func (s *TaskServiceServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskRresponse, error) {
	taskID := uuid.New().String()

	task := Task{
		ID:          taskID,
		Title:       req.Task.Title,
		Description: req.Task.Description,
		Interval:    time.Duration(req.Task.Interval),
	}

	if err := SaveTask(ctx, task); err != nil {
		return nil, err
	}

	return &pb.CreateTaskRresponse{Id: taskID}, nil
}

func (s *TaskServiceServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskRresponse, error) {
	task, err := GetTask(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetTaskRresponse{
		Task: &pb.Task{Id: task.ID,
			Title:       task.Title,
			Description: task.Description,
			Interval:    int64(task.Interval),
		}}, nil
}

func (s *TaskServiceServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskRresponse, error) {
	updatedTask := Task{
		ID:          req.Id,
		Title:       req.Task.Title,
		Description: req.Task.Description,
		Interval:    time.Duration(req.Task.Interval),
	}

	err := UpdateTask(ctx, req.Id, updatedTask)
	if err != nil {
		return &pb.UpdateTaskRresponse{Success: false}, err
	}

	return &pb.UpdateTaskRresponse{Success: true}, nil
}

func (s *TaskServiceServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	if err := DeleteTask(ctx, req.Id); err != nil {
		return &pb.DeleteTaskResponse{Success: false}, err
	}

	return &pb.DeleteTaskResponse{Success: true}, nil
}
