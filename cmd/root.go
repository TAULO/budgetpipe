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

func getWorkDir() (string, error) {
	dir := viper.GetString(flags.Dir)
	if dir == "" {
		return os.Getwd()
	}

	return filepath.Abs(dir)
}

func getCSVFilePath() (string, error) {
	dir, err := getWorkDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, filenames.CSV), nil
}

func getMapperFilePath() (string, error) {
	dir, err := getWorkDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, filenames.Mapper), nil
}

func getXlSXFilePath() (string, error) {
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

var rootCmd = &cobra.Command{
	Use:   "budget",
	Short: "Import transactions into your budget",
	//RunE: func(cmd *cobra.Command, args []string) error {
	//	csvPath := cmd.Flags().Lookup("csv").Value.String()
	//	xlsxPath := cmd.Flags().Lookup("xlsx").Value.String()
	//	month := cmd.Flags().Lookup("month").Value.String()
	//
	//	mapperPath := defaultMapperPath()
	//
	//	return budget.Run(csvPath, xlsxPath, mapperPath, month)
	//},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	//rootCmd.Flags().String("csv", "", "path to the bank data CSV file")
	//rootCmd.Flags().String("xlsx", "", "path to the budget Excel file")
	//rootCmd.Flags().String("month", "", "month to import")
	//_ = rootCmd.MarkFlagRequired("csv")
	//_ = rootCmd.MarkFlagRequired("xlsx")
	//_ = rootCmd.MarkFlagRequired("month")

	rootCmd.PersistentFlags().StringP(flags.Dir, "C", "", "Working directory (default: current directory)")
	rootCmd.PersistentFlags().StringP(flags.Sheet, "S", "", "The name of the sheet (default: current year)")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		panic(err)
	}
}

func defaultMapperPath() string {
	path, _ := os.Getwd()

	return filepath.Join(
		path,
		"config",
		"mapper.json",
	)
}
