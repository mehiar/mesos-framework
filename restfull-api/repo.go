package restapi

import "fmt"

var currentId int

var containers ContainerRequests
var submission_requests ContainerRequests

func RepoFindRequest(id int) ContainerRequest {
	for _, t := range containers {
		if t.Id == id {
			return t
		}
	}
	// return empty ContainerRequest if not found
	return ContainerRequest{}
}

//this is bad, I don't think it passes race condtions
func RepoCreateRequest(c ContainerRequest) ContainerRequest {
	currentId += 1
	c.Id = currentId
	containers = append(containers, c)
	submission_requests = append(submission_requests,c)
	return c
}

func RepoRemoveRequest(id int) error {
	for i, c := range containers {
		if c.Id == id {
			containers = append(containers[:i], containers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find ContainerRequest with id of %d to delete", id)
}

func RepoGetSubmissionRequests() *ContainerRequests{
	return &submission_requests
}
