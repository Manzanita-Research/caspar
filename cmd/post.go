package cmd

func init() {
	rootCmd.AddCommand(buildResourceCommands(kindPost))
}
