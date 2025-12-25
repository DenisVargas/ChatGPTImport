package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "main.go [source] [optional output]",
	Short: "Parses ChatGPT conversation exports into markdown files",
	Long: `This tool processes JSON exports of ChatGPT conversations and converts them into well-formatted markdown files,
exporting them as individual conversation files to the specified output location or the current directory if no output is provided.`,
	Args: cobra.RangeArgs(0, 2),
	Run:  rootCommand,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntP("limit", "l", 0, "Limit the number of conversations to process (0 for no limit, Default)")
}

var conversationsFile string = "scrap.json"
var outputDir string = "./output"

func rootCommand(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		cmd.Help()
		return
	}

	if len(args) >= 1 {
		conversationsFile = args[0]
	}

	if len(args) == 2 {
		outputDir = args[1]
	}

	fileInfo, err := os.Stat(conversationsFile)
	if err != nil {
		panic(err)
	}

	if fileInfo.IsDir() {
		conversationsFile = conversationsFile + "/conversations.json"
	}

	conversations, err := loadConversations(conversationsFile)
	if err != nil {
		panic(err)
	}

	if len(conversations) == 0 {
		fmt.Println("No conversations found in the provided file.")
		return
	}

	limit, _ := cmd.Flags().GetInt("limit")
	if limit == 0 || (limit != 0 && len(conversations) < limit) {
		limit = len(conversations)
	}

	for _, conv := range conversations[:limit] {

		md := composeMarkdown(conv)
		safeTitle := strings.ReplaceAll(conv.Title, " ", "_")
		safeTitle = strings.ReplaceAll(safeTitle, "/", "_")
		filename := fmt.Sprintf("%s/%s.md", outputDir, safeTitle)

		if err := renderMarkdown(md, filename); err != nil {
			fmt.Printf("Error rendering %s: %v\n", filename, err)
		} else {
			fmt.Printf("✅ Generated: %s\n", filename)
		}
	}
}
