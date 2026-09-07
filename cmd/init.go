package cmd

import (
	"budgetpipe/internal/model"
	"budgetpipe/internal/xlsx"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new budget",
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkDir()
		if err != nil {
			return err
		}

		if err := os.MkdirAll(workDir, 0o755); err != nil {
			return fmt.Errorf("creating work dir %s: %w", workDir, err)
		}

		dir, err := os.Open(workDir)
		if err != nil {
			return fmt.Errorf("opening work dir: %w", err)
		}
		entries, err := dir.Readdirnames(1)
		if err := dir.Close(); err != nil {
			return fmt.Errorf("closing work dir: %w", err)
		}
		if err != nil && err != io.EOF {
			return fmt.Errorf("reading work dir: %w", err)
		}
		if len(entries) > 0 {
			return fmt.Errorf("budget already initialized in %s", workDir)
		}

		if err := createDataCSVFile(); err != nil {
			return fmt.Errorf("creating data csv file: %w", err)
		}

		wb, err := xlsx.NewBudget(template, getSheetName())
		if err != nil {
			return fmt.Errorf("creating budget file: %w", err)
		}

		xlsxPath, err := getXlSXFilePath()
		if err != nil {
			return fmt.Errorf("resolving xlsx path: %w", err)
		}
		if err := wb.SaveAs(xlsxPath); err != nil {
			return fmt.Errorf("saving budget file: %w", err)
		}

		mapperFilePath, err := getMapperFilePath()
		if err != nil {
			return fmt.Errorf("resolving mapper path: %w", err)
		}

		starter := model.NewMapper()
		starterMapper, err := json.MarshalIndent(starter, "", "  ")
		if err != nil {
			return fmt.Errorf("building starter mapper: %w", err)
		}

		if err := os.WriteFile(mapperFilePath, starterMapper, 0o644); err != nil {
			return fmt.Errorf("creating mapper: %w", err)
		}

		fmt.Printf("initialized budget in %s\n", workDir)
		return nil
	},
}

func createDataCSVFile() error {
	months := []string{
		"jan", "feb", "mar", "apr", "may", "jun",
		"jul", "aug", "sep", "okt", "nov", "dec",
	}

	dataDir, err := getDataDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}

	for _, month := range months {
		path := filepath.Join(dataDir, month+".csv")

		if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
