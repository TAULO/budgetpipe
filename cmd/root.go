package cmd

import (
	"budgetpipe/filenames"
	"budgetpipe/flags"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// for testing:
// budget  --csv "data/csv/august.csv" --xlsx "data/test.xlsx" --month august
// TODO:
// budget add  --csv "data/csv/august.csv" --xlsx "data/test.xlsx" --month august

var template []byte

func getWorkDir() (string, error) {
	dir := viper.GetString(flags.Dir)
	if dir == "" {
		return os.Getwd()
	}

	return filepath.Abs(dir)
}

func getDataDir() (string, error) {
	dir, err := getWorkDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "data"), nil
}

func getMapperFilePath() (string, error) {
	dir, err := getWorkDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, filenames.Mapper), nil
}

func getBudgetFilePath() (string, error) {
	dir, err := getWorkDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filenames.XLSX), nil
}

func getSheetName() string {
	sheetName := viper.GetString(flags.Sheet)
	currYear := strconv.Itoa(time.Now().Year())
	if sheetName == "" {
		return currYear
	}

	return sheetName
}

func getDataCSVFilePath(month string) (string, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dataDir, month+".csv"), nil
}

var rootCmd = &cobra.Command{
	Use:   "budget",
	Short: "Import transactions into your budget",
}

func Execute(templateBytes []byte) error {
	template = templateBytes
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP(flags.Dir, "C", "", "Working directory (default: current directory)")
	rootCmd.PersistentFlags().StringP(flags.Sheet, "S", "", "The name of the sheet (default: current year)")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		panic(err)
	}
}
