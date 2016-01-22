/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package scheduler

import (
	"github.com/gogo/protobuf/proto"
	"strconv"
	"math"
	"fmt"
	"sort"

	log "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
	. "github.com/mehiar/mesos-framework/restfull-api"
)

type ExampleScheduler struct {
	executor      *mesos.ExecutorInfo
	tasksLaunched int
	tasksFinished int
	totalTasks    int
	cpuPerTask    float64
	memPerTask    float64
}

func NewExampleScheduler(exec *mesos.ExecutorInfo, taskCount int, cpuPerTask float64, memPerTask float64) *ExampleScheduler {
	return &ExampleScheduler{
		executor:      exec,
		tasksLaunched: 0,
		tasksFinished: 0,
		totalTasks:    taskCount,
		cpuPerTask:    cpuPerTask,
		memPerTask:    memPerTask,
	}
}

func (sched *ExampleScheduler) Registered(driver sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	log.Infoln("Scheduler Registered with Master ", masterInfo)
}

func (sched *ExampleScheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	log.Infoln("Scheduler Re-Registered with Master ", masterInfo)
}

func (sched *ExampleScheduler) Disconnected(sched.SchedulerDriver) {
	log.Infoln("Scheduler Disconnected")
}

func (sched *ExampleScheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {
	logOffers(offers)
	var maxCpu float64 = 0
	var maxMem float64 = 0
	var resourceTable [][]float64
	var metric []float64 
	for _, offer := range offers{
		remainingCpus := getOfferCpu(offer)
		remainingMems := getOfferMem(offer)
		temp :=[]float64{remainingCpus, remainingMems, 0}
		resourceTable = append(resourceTable,temp)
		maxCpu = math.Max(maxCpu,remainingCpus)
		maxMem = math.Max(maxMem,remainingMems)
	}

	//Normalization
	for itr:=0; itr<len(resourceTable); itr++{
		resourceTable[itr][0] = resourceTable[itr][0] / maxCpu
		resourceTable[itr][1] = resourceTable[itr][1] /maxMem
		resourceTable[itr][2] = resourceTable[itr][0] * resourceTable[itr][1]
		metric = append( metric, resourceTable[itr][2])
		log.Infof("MEHIAR: %v, %v, %v \n",resourceTable[itr][0],resourceTable[itr][1],resourceTable[itr][2])
	}

	//Sorting
	s := NewSlice(metric)
	sort.Sort(s)
	fmt.Println("MEHIAR: ", s.Float64Slice, s.idx)
	submission_requests_ptr := RepoGetSubmissionRequests()
	fmt.Println("MEHIAR: ", *submission_requests_ptr)
	
	for _, offer := range offers {
		remainingCpus := getOfferCpu(offer)
		remainingMems := getOfferMem(offer)
		var tasks []*mesos.TaskInfo
		for sched.cpuPerTask <= remainingCpus &&
			sched.memPerTask <= remainingMems {

			sched.tasksLaunched++

			taskId := &mesos.TaskID{
				Value: proto.String(strconv.Itoa(sched.tasksLaunched)),
			}
			containerType := mesos.ContainerInfo_DOCKER
			task := &mesos.TaskInfo{
				Name:     proto.String("go-task-" + taskId.GetValue()),
				TaskId:   taskId,
				SlaveId:  offer.SlaveId,
//				Executor: sched.executor,
				Container: &mesos.ContainerInfo {
      					Type: &containerType,
        				Docker: &mesos.ContainerInfo_DockerInfo {
			            		Image: proto.String("ubuntu"),
        				},
    				},
				Command: &mesos.CommandInfo {
					Value: proto.String("sleep 1000000"),
				},
				Resources: []*mesos.Resource{
					util.NewScalarResource("cpus", sched.cpuPerTask),
					util.NewScalarResource("mem", sched.memPerTask),
				},
			}
			log.Infof("Prepared task: %s with offer %s for launch\n", task.GetName(), offer.Id.GetValue())

			tasks = append(tasks, task)
			remainingCpus -= sched.cpuPerTask
			remainingMems -= sched.memPerTask
		}
		log.Infoln("Launching ", len(tasks), "tasks for offer", offer.Id.GetValue())
		driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, &mesos.Filters{RefuseSeconds: proto.Float64(1)})
	}

}

func (sched *ExampleScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	log.Infoln("Status update: task", status.TaskId.GetValue(), " is in state ", status.State.Enum().String())
}

func (sched *ExampleScheduler) OfferRescinded(s sched.SchedulerDriver, id *mesos.OfferID) {
	log.Infof("Offer '%v' rescinded.\n", *id)
}

func (sched *ExampleScheduler) FrameworkMessage(s sched.SchedulerDriver, exId *mesos.ExecutorID, slvId *mesos.SlaveID, msg string) {
	log.Infof("Received framework message from executor '%v' on slave '%v': %s.\n", *exId, *slvId, msg)
}

func (sched *ExampleScheduler) SlaveLost(s sched.SchedulerDriver, id *mesos.SlaveID) {
	log.Infof("Slave '%v' lost.\n", *id)
}

func (sched *ExampleScheduler) ExecutorLost(s sched.SchedulerDriver, exId *mesos.ExecutorID, slvId *mesos.SlaveID, i int) {
	log.Infof("Executor '%v' lost on slave '%v' with exit code: %v.\n", *exId, *slvId, i)
}

func (sched *ExampleScheduler) Error(driver sched.SchedulerDriver, err string) {
	log.Infoln("Scheduler received error:", err)
}
