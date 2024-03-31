package updater

import (
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/AlvareYN/auto-updater/cmd"
	"github.com/gin-gonic/gin"
	"golang.org/x/mod/semver"
)

const (
	Version = "v0.0.12"
	AppName = "auto-updater"
)

func CheckUpdates(c *gin.Context) {
	latestVersion, err := cmd.GetRemoteVersionMetadata()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message":         "Check updates successful",
		"latest_version":  latestVersion,
		"current_version": Version,
		"semver_status":   semver.Compare(Version, latestVersion["current"].(string)),
	})
}

func Update(c *gin.Context) {
	latestVersion, err := cmd.GetRemoteVersionMetadata()

	if err != nil {
		return
	}
	log.Println("Latest version: ", latestVersion["current"].(string), "Current version: ", Version)
	log.Println("Current version: ", semver.Compare(Version, latestVersion["current"].(string)))
	if semver.Compare(Version, latestVersion["current"].(string)) >= 0 {
		c.JSON(200, gin.H{
			"message": "No updates available",
		})
		return
	}

	build := latestVersion["build"].(map[string]interface{})

	err = cmd.FetchRemoteVersion(latestVersion["current"].(string), build[runtime.GOOS].(string))

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Update successful",
		"version": latestVersion["current"],
		"build":   build,
		"binary":  path.Join(".", build[runtime.GOOS].(string)),
	})

}

type ApplyUpdatesRequest struct {
	Build string `json:"build"`
}

func ApplyUpdates(c *gin.Context) {
	latestVersion, err := cmd.GetRemoteVersionMetadata()

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	build := latestVersion["build"].(map[string]interface{})

	var req ApplyUpdatesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "invalid build:" + err.Error(),
		})
		return
	}

	if req.Build != build[runtime.GOOS].(string) {

		c.JSON(400, gin.H{
			"message": "Invalid build",
		})
		return
	}

	command := exec.Command("./" + build[runtime.GOOS].(string))

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Start()

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error starting command: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Update applied",
	})

	os.Exit(0)
}

func GetVersion(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
		"version": Version,
	})
}
