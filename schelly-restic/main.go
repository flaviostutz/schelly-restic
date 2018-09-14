package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/flaviostutz/schelly-webhook/schellyhook"
)

var sourcePath *string
var repoDir *string

//ResticBackuper sample backuper
type ResticBackuper struct {
}

func main() {
	logrus.Info("====Starting Restic REST server====")

	resticBackuper := ResticBackuper{}
	err := schellyhook.Initialize(resticBackuper)
	if err != nil {
		logrus.Errorf("Error initializating Schellyhook. err=%s", err)
		os.Exit(1)
	}
}

//RegisterFlags register command line flags
func (sb ResticBackuper) RegisterFlags() error {
	sourcePath = flag.String("source-path", "/backup-source", "Backup source path")
	repoDir = flag.String("repo-dir", "/backup-repo", "Restic repository of backups")
	return nil
}

//Init initialize
func (sb ResticBackuper) Init() error {
	err := mkDirs(baseIDDir())
	if err != nil {
		logrus.Errorf("Couldn't create id references. err=%s", err)
		return err
	}

	logrus.Debugf("Checking if Restic repo %s was already initialized", *repoDir)
	result, err := schellyhook.ExecShell("restic snapshots -r " + *repoDir)
	if err != nil {
		logrus.Debugf("Couldn't access Restic repo. Trying to create it. err=", err)
		_, err := schellyhook.ExecShell("restic init -r " + *repoDir)
		if err != nil {
			logrus.Debugf("Error creating Restic repo: %s %s", err, result)
			return err
		} else {
			logrus.Infof("Restic repo created successfuly")
		}
	} else {
		logrus.Infof("Restic repo already exists and is accessible")
	}
	return nil
}

//CreateNewBackup creates a new backup
func (sb ResticBackuper) CreateNewBackup(apiID string, timeout time.Duration, shellContext *schellyhook.ShellContext) error {
	logrus.Infof("CreateNewBackup() apiID=%s timeout=%d s", apiID, timeout.Seconds)

	logrus.Infof("Calling Restic...")
	result, err := schellyhook.ExecShell("restic backup /backup-source -r " + *repoDir)
	if err != nil {
		return err
	}
	logrus.Debugf("result: %s", result)
	rex, _ := regexp.Compile("snapshot ([0-9a-zA-z]+) saved")
	id := rex.FindStringSubmatch(result)
	success := (len(id) == 2)
	if !success {
		logrus.Warnf("Snapshot not created. result=%s", result)
	} else {
		dataID := id[1]
		err = saveDataID(apiID, dataID)
		if err != nil {
			return err
		}
		logrus.Infof("Backup finished")
	}

	return nil
}

//GetAllBackups returns all backups from underlaying backuper. optional for Schelly
func (sb ResticBackuper) GetAllBackups() ([]schellyhook.SchellyResponse, error) {
	logrus.Debugf("GetAllBackups")
	result, err := schellyhook.ExecShell("restic snapshots --quiet -r " + *repoDir)
	if err != nil {
		return nil, err
	}

	backups := make([]schellyhook.SchellyResponse, 0)
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		space := regexp.MustCompile(`\s+`)
		line = space.ReplaceAllString(line, " ")
		cols := strings.Split(line, " ")
		if i == 0 || len(cols) < 3 {
			continue
		}

		dataID := cols[0]

		sr := schellyhook.SchellyResponse{
			// ID:      getAPIID(dataID),
			DataID:  dataID,
			Status:  "available",
			Message: line,
			SizeMB:  -1,
		}
		backups = append(backups, sr)
	}

	return backups, nil
}

//GetBackup get an specific backup along with status
func (sb ResticBackuper) GetBackup(apiID string) (*schellyhook.SchellyResponse, error) {
	logrus.Debugf("GetBackup apiID=%s", apiID)

	dataID, err0 := getDataID(apiID)
	if err0 != nil {
		logrus.Debugf("BackyID not found for apiId %s. err=%s", apiID, err0)
		return nil, nil
	}

	res, err := findBackup(apiID, dataID)
	if err != nil {
		return nil, nil
	}

	return res, nil
}

//DeleteBackup removes current backup from underlaying backup storage
func (sb ResticBackuper) DeleteBackup(apiID string) error {
	logrus.Debugf("DeleteBackup apiID=%s", apiID)

	dataID, err0 := getDataID(apiID)
	if err0 != nil {
		logrus.Debugf("Restic backup not found for apiID %s dataID %s. err=%s", apiID, dataID, err0)
		return err0
	}

	_, err0 = findBackup(apiID, dataID)
	if err0 != nil {
		logrus.Debugf("Backup apiID %s, dataID %s not found for removal", apiID, dataID)
		return err0
	}

	logrus.Debugf("Backup apiID=%s dataID=%s found. Proceeding to deletion", apiID, dataID)
	result, err := schellyhook.ExecShell("restic forget " + dataID + " -r " + *repoDir)
	if err != nil {
		return err
	}
	logrus.Debugf("result: %s", result)

	rex, _ := regexp.Compile("removed snapshot ([0-9a-zA-z]+)")
	id := rex.FindStringSubmatch(result)
	if len(id) != 2 {
		return fmt.Errorf("Couldn't find returned id from response")
	}

	if id[1] != dataID {
		return fmt.Errorf("Returned id from forget is different from requested. %s != %s", id[1], dataID)
	}

	logrus.Debugf("Delete apiID %s dataID %s successful", apiID, dataID)
	return nil
}

func getDataID(apiID string) (string, error) {
	fn := baseIDDir() + apiID
	if _, err := os.Stat(fn); err == nil {
		logrus.Debugf("Found api id reference for %s", apiID)
		data, err0 := ioutil.ReadFile(fn)
		if err0 != nil {
			return "", err0
		} else {
			dataID := string(data)
			logrus.Debugf("apiID %s <-> dataID %s", apiID, dataID)
			return dataID, nil
		}
	} else {
		return "", fmt.Errorf("dataID for apiID %s not found", apiID)
	}
}

func saveDataID(apiID string, dataID string) error {
	logrus.Debugf("Setting apiID %s <-> dataID %s", apiID, dataID)
	fn := baseIDDir() + apiID
	if _, err := os.Stat(fn); err == nil {
		err = os.Remove(fn)
		if err != nil {
			return fmt.Errorf("Couldn't replace existing apiID file. err=%s", err)
		}
	}
	return ioutil.WriteFile(fn, []byte(dataID), 0644)
}

func mkDirs(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func findBackup(apiID string, dataID string) (*schellyhook.SchellyResponse, error) {
	result, err := schellyhook.ExecShell("restic snapshots " + dataID + " -r " + *repoDir)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("Query snapshots result: %s", result)

	rex, _ := regexp.Compile("-\n([0-9a-z]{4,16})")
	id0 := rex.FindStringSubmatch(result)
	if len(id0) != 2 {
		logrus.Debug("Couldn't find backup id in response '%'", id0, result)
		return nil, nil
	}

	logrus.Debugf("Snapshot %s found", id0)
	return &schellyhook.SchellyResponse{
		ID:      id0[1],
		DataID:  dataID,
		Status:  "available",
		Message: result,
	}, nil
}

func baseIDDir() string {
	return *repoDir + "/ids/"
}
