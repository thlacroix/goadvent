package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var rStep = regexp.MustCompile(`Step (\w) must be finished before step (\w) can begin.`)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if res, _, err := getOrder(fileName, 1, -1); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Result for one worker is", res)
	}

	if _, time, err := getOrder(fileName, 5, 60); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Time for 5 workers is", time)
	}
}

type Step struct {
	Name     string
	Before   []*Step // Steps that needs this step done before
	After    []*Step // Steps that are prerequisites of this one
	Done     bool
	Progress int
	Taken    bool
}

func (s Step) timeToProcess(baseTime int) int {
	if baseTime < 0 {
		return 0
	}
	return baseTime + int(s.Name[0]) - 64
}

type Worker struct {
	ID          int
	CurrentTask *Step
}

func getOrder(fileName string, workerCount int, baseTime int) (string, int, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	steps := make(map[string]*Step)
	var stepOrder []*Step

	scanner := bufio.NewScanner(file)
	// first building the step mapping (name -> Step object)
	for scanner.Scan() {
		line := scanner.Text()
		extract := rStep.FindAllStringSubmatch(line, -1)[0]
		before := extract[1]
		after := extract[2]

		var beforeStep, afterStep *Step

		if _, ok := steps[before]; !ok {
			beforeStep = &Step{
				Name: before,
			}
			steps[before] = beforeStep
		} else {
			beforeStep = steps[before]
		}
		if _, ok := steps[after]; !ok {
			afterStep = &Step{
				Name: after,
			}
			steps[after] = afterStep
		} else {
			afterStep = steps[after]
		}
		beforeStep.Before = append(beforeStep.Before, afterStep)
		afterStep.After = append(afterStep.After, beforeStep)
	}

	// sorting the step dependencies alphabetically
	for _, step := range steps {
		sort.Slice(step.Before, func(i, j int) bool {
			return step.Before[i].Name < step.Before[j].Name
		})
		sort.Slice(step.After, func(i, j int) bool {
			return step.After[i].Name < step.After[j].Name
		})
	}

	// maps can't be iterated alphabetically in Go by default, so building a list
	// and sorting it
	stepList := make([]*Step, 0, len(steps))
	for _, step := range steps {
		stepList = append(stepList, step)
	}
	sort.Slice(stepList, func(i, j int) bool {
		return stepList[i].Name < stepList[j].Name
	})

	// building worker list
	var workers []*Worker
	for i := 0; i < workerCount; i++ {
		workers = append(workers, &Worker{ID: i})
	}

	// building the initial list of possible steps, ie steps without prerequisites
	var possibleNextSteps []*Step

	// building first list of possible steps
	for _, step := range stepList {
		if len(step.After) == 0 {
			possibleNextSteps = append(possibleNextSteps, step)
		}
	}

	var totalTime int

	// stopping when all task have been done
	for len(stepOrder) != len(steps) {
		// each work work on their assigned tasks. Nothing to do in the first loop
		for _, worker := range workers {
			if worker.CurrentTask != nil {
				worker.CurrentTask.Progress++
				// if the task is finished, we mark it as done, and add its after steps
				// to the possible next step list (that we keep sorted)
				if worker.CurrentTask.Progress >= worker.CurrentTask.timeToProcess(baseTime) {
					worker.CurrentTask.Done = true
					stepOrder = append(stepOrder, worker.CurrentTask)
					possibleNextSteps = append(possibleNextSteps, worker.CurrentTask.Before...)
					// instead of sorting, could be inserted directly at the right place,
					// didn't want to bother here
					sort.Slice(possibleNextSteps, func(i, j int) bool {
						return possibleNextSteps[i].Name < possibleNextSteps[j].Name
					})
					worker.CurrentTask = nil
				}
			}
		}

		// each worker pick a step in the possible next steps
		for _, worker := range workers {
			if worker.CurrentTask == nil {
				// looking for a task for the worker
				for _, nextStep := range possibleNextSteps {
					if !nextStep.Taken {
						ready := true
						// checking that all requirements are done
						for _, reqStep := range nextStep.After {
							if !reqStep.Done {
								ready = false
								break
							}
						}
						if ready {
							worker.CurrentTask = nextStep
							nextStep.Taken = true
							// we could also remove the task from the possible next step list
							break
						}
					}
				}
			}
		}
		totalTime++
	}

	var res strings.Builder

	// building the step order string
	for _, step := range stepOrder {
		res.WriteString(step.Name)
	}
	return res.String(), totalTime - 1, nil
}
