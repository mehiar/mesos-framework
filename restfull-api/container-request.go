package restapi

import "time"

type ContainerRequest struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Due       time.Time `json:"due"`
	Image     string    `json:"image"`
	Command   string    `json:"command"`
	Cpus      float64   `json:"cpus"`
	Mem       int	    `json:"mem"`
	Status    string    `json:"status"`
}

type ContainerRequests []ContainerRequest
