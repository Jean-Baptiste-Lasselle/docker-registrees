package settings

import (
	"errors"

	"github.com/Sirupsen/logrus"
)

// ReleaseVersion contains the git tag version
var ReleaseVersion string

func init() {
	sha, err := GetAppSHA()
	if err != nil {
		logrus.Error(err.Error())
		logrus.Error("Could not get git sha version")
	}
	ReleaseVersion = sha
}

// UpdateApp stops the current instance of the app, updates to the latest sha on the branch, and restarts the app
func UpdateApp() (bool, error) {

	isUpToDate, err := IsAppUpToDate()
	if err != nil {
		return false, err
	}

	if isUpToDate {
		return false, errors.New("Already up to date!")
	}

	return true, nil
}

// IsAppUpToDate checks to see if the local status of the git tree is up to date with the remote
func IsAppUpToDate() (bool, error) {

	// Fetch origin refs
	RemoteUpdate()

	// Get the local branch info
	localBranchName, err := GetAppBranch()
	if err != nil {
		return true, err
	}
	localBranchSHA, err := GetAppSHA()
	if err != nil {
		return true, err
	}

	// Get the remote info
	remoteBaseSHA, err := GetBaseSHA()
	if err != nil {
		return true, err
	}
	remoteBranchSHA, err := GetRemoteBranchSHA(localBranchName)
	if err != nil {
		return true, err
	}

	// Compare the local SHA and remote SHA, if they're the same we are up to date
	// http://stackoverflow.com/questions/3258243/check-if-pull-needed-in-git
	if localBranchSHA == remoteBranchSHA {
		return true, nil
	} else if localBranchSHA == remoteBaseSHA {
		// This means we need to update
		return false, nil
	}

	// If branch is diverged or you need to push
	return true, err
}
