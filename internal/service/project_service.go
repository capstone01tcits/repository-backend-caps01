package service

import (
	"errors"

	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"

	"github.com/google/uuid"
)

type ProjectService interface {
	CreateProject(userID string, req *model.CreateProjectRequest) (*model.Project, error)
	GetProject(userID, projectID string) (*model.Project, error)
	GetProjects(userID string) ([]model.Project, error)
	UpdateProject(userID, projectID string, req *model.UpdateProjectRequest) (*model.Project, error)
	DeleteProject(userID, projectID string) error
}

type projectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{projectRepo}
}

func (s *projectService) CreateProject(userID string, req *model.CreateProjectRequest) (*model.Project, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	project := &model.Project{
		UserID:      uid,
		Name:        req.Name,
		Description: req.Description,
		Status:      "draft",
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, errors.New("failed to create project")
	}

	return project, nil
}

func (s *projectService) GetProject(userID, projectID string) (*model.Project, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this project")
	}

	return project, nil
}

func (s *projectService) GetProjects(userID string) ([]model.Project, error) {
	return s.projectRepo.FindByUserID(userID)
}

func (s *projectService) UpdateProject(userID, projectID string, req *model.UpdateProjectRequest) (*model.Project, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this project")
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	if req.Status != nil {
		project.Status = *req.Status
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, errors.New("failed to update project")
	}

	return project, nil
}

func (s *projectService) DeleteProject(userID, projectID string) error {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return errors.New("project not found")
	}

	if project.UserID.String() != userID {
		return errors.New("unauthorized access to this project")
	}

	return s.projectRepo.Delete(projectID)
}
