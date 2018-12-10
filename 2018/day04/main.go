package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var rLog = regexp.MustCompile(`^\[(.+)\] (.+)$`)
var rGuard = regexp.MustCompile(`^Guard #(\d+) begins shift$`)

type LogType int

const (
	GoToSleep LogType = iota
	WakeUp
	TakeShift
)

type Log struct {
	TS      time.Time
	Action  string
	Type    LogType
	GuardID int
}

type SleepRecord struct {
	Start int
	End   int
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filename passed")
	}
	logs, err := getSortedLogs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	part1, part2 := findSleepingBeauty(logs)
	fmt.Println("Part1 solution is", part1, "and Part2 solution is", part2)
}

// parsing the logs and sorting them by date. Could be done without parsing the
// date completely as the can sort the log strings alphabetically and just
// extract the minutes, but works fine this way for the input size
func getSortedLogs(fileName string) ([]Log, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	logs := make([]Log, 0)
	for scanner.Scan() {
		line := scanner.Text()
		extracts := rLog.FindAllStringSubmatch(line, -1)
		if len(extracts) != 1 || len(extracts[0]) != 3 {
			return nil, errors.New("Can't parse log")
		}
		extract := extracts[0]

		// Parsing date
		ts, err := time.Parse("2006-01-02 15:04", extract[1])
		if err != nil {
			return nil, err
		}

		log := Log{TS: ts, Action: extract[2]}

		// getting action type and guard id if any
		switch extract[2] {
		case "wakes up":
			log.Type = WakeUp
		case "falls asleep":
			log.Type = GoToSleep
		default:
			log.Type = TakeShift
			guardExtracts := rGuard.FindAllStringSubmatch(extract[2], -1)
			if len(guardExtracts) != 1 || len(guardExtracts[0]) != 2 {
				return nil, errors.New("Can't extract guard ID")
			}
			guardID, err := strconv.Atoi(guardExtracts[0][1])
			if err != nil {
				return nil, errors.New("Guard ID is not an int")
			}
			log.GuardID = guardID
		}
		logs = append(logs, log)
	}

	// sorting the logs by date
	sort.Slice(logs, func(i int, j int) bool {
		return logs[i].TS.Before(logs[j].TS)
	})
	return logs, nil
}

func findSleepingBeauty(logs []Log) (int, int) {
	countMap := make(map[int]int)
	sleepRecords := make(map[int][]SleepRecord)

	var currentGuard int
	var bedTime time.Time

	// counting sleep time for each guard, and recording the start and end for
	// each sleep periods, usefull for part 2
	for _, log := range logs {
		switch log.Type {
		case GoToSleep:
			bedTime = log.TS
		case TakeShift:
			currentGuard = log.GuardID
		case WakeUp:
			diff := log.TS.Sub(bedTime)
			countMap[currentGuard] += int(diff.Minutes())
			sleepRecords[currentGuard] = append(sleepRecords[currentGuard], SleepRecord{Start: bedTime.Minute(), End: log.TS.Minute()})
		}
	}
	// finding the guard that slept the most. Could be done above as we count the
	// sleep also to optimize
	var maxSleep int
	var sleepier int
	for guard, sleepTime := range countMap {
		if sleepTime > maxSleep {
			maxSleep = sleepTime
			sleepier = guard
		}
	}

	// we now get the minute where the sleepier slept the most
	maxMinute, _ := minuteFrequency(sleepRecords[sleepier])

	var maxGuard, maxGuardMinute, maxGuardMinuteCount int

	// for each guard, we get the minute where he sleeps the most, and keep the
	// one with the higher frequency
	for guard, records := range sleepRecords {
		guardMinute, guardMinuteCount := minuteFrequency(records)
		if guardMinuteCount > maxGuardMinuteCount {
			maxGuard = guard
			maxGuardMinute = guardMinute
			maxGuardMinuteCount = guardMinuteCount
		}
	}

	return sleepier * maxMinute, maxGuard * maxGuardMinute
}

// simply counting for each sleep period the frequency of each minute, and
// returning the most frequent and how frequent it is
func minuteFrequency(records []SleepRecord) (int, int) {
	minuteCount := make(map[int]int, 0)
	for _, record := range records {
		for i := record.Start; i < record.End; i++ {
			minuteCount[i]++
		}
	}

	var maxMinuteCount, maxMinute int
	for minute, count := range minuteCount {
		if count > maxMinuteCount {
			maxMinuteCount = count
			maxMinute = minute
		}
	}
	return maxMinute, maxMinuteCount
}
