package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Communinst/CWDB6Sem/backend/config"
	entities "github.com/Communinst/CWDB6Sem/backend/entity"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createDump(c *gin.Context) {
	dumpCfg := config.MustLoadDumpConfig()
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02_15-04-05")
	filePath := filepath.Join(dumpCfg.Dir, fmt.Sprintf("%s_%s", dumpCfg.Prefix, timeString))

	fmt.Printf("%s\n", filePath)

	cmd := exec.Command("docker", "exec", dumpCfg.ContainerName, "pg_dump", "-U", dumpCfg.Username, "-F", "c", dumpCfg.DbName)
	outputFile, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer outputFile.Close()

	cmd.Stdout = outputFile

	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileInfo, err := outputFile.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fileSize := fileInfo.Size()

	err = h.service.DumpServiceInterface.InsertDump(c.Request.Context(), filePath, fileSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dump created successfully"})
}

func (h *Handler) restoreDump(c *gin.Context) {
	var fileName entities.Dump

	if err := c.ShouldBindJSON(&fileName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	dumpCfg := config.MustLoadDumpConfig()

	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02_15-04-05")
	filePath := fmt.Sprintf("%s_%s", dumpCfg.RestorePrefix, timeString)

	copyCmd := exec.Command("docker", "cp", fileName.Filename, fmt.Sprintf("%s:%s", dumpCfg.ContainerName, filePath))
	if err := copyCmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	restoreCmd := exec.Command("docker", "exec", dumpCfg.ContainerName, "pg_restore", "-U", dumpCfg.Username, "--clean", "-d", dumpCfg.DbName, filePath)
	if err := restoreCmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dump restored successfully"})
}

func (h *Handler) getAllDumps(c *gin.Context) {
	dumps, err := h.service.DumpServiceInterface.GetAllDumps(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("%v", dumps)
	c.JSON(http.StatusOK, dumps)
}
