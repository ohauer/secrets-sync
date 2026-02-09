package filewriter

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// CleanupOrphanedTempFiles removes stale .tmp files from output directories
func CleanupOrphanedTempFiles(outputDirs []string, logger *zap.Logger) error {
	cleaned := 0

	for _, dir := range outputDirs {
		if dir == "" {
			continue
		}

		// Check if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		// Find all .tmp files in directory
		pattern := filepath.Join(dir, "*.tmp")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			logger.Warn("failed to glob temp files",
				zap.String("dir", dir),
				zap.Error(err),
			)
			continue
		}

		// Remove each .tmp file
		for _, tmpFile := range matches {
			if err := os.Remove(tmpFile); err != nil {
				logger.Warn("failed to remove orphaned temp file",
					zap.String("file", tmpFile),
					zap.Error(err),
				)
			} else {
				logger.Info("removed orphaned temp file",
					zap.String("file", tmpFile),
				)
				cleaned++
			}
		}
	}

	if cleaned > 0 {
		logger.Info("cleanup complete",
			zap.Int("files_removed", cleaned),
		)
	}

	return nil
}

// GetOutputDirectories extracts unique output directories from file paths
func GetOutputDirectories(filePaths []string) []string {
	dirMap := make(map[string]bool)

	for _, path := range filePaths {
		if path == "" {
			continue
		}
		dir := filepath.Dir(path)
		dirMap[dir] = true
	}

	dirs := make([]string, 0, len(dirMap))
	for dir := range dirMap {
		dirs = append(dirs, dir)
	}

	return dirs
}
