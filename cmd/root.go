package cmd

import (
	"os"
	"sync"

	"github.com/Minosity-VR/confdump/internal/client"
	"github.com/Minosity-VR/confdump/internal/saver"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "confdump",
	Version: "0.0.1",
	Short:   "Dump a confluence wiki",
	Long: `Dump a confluence wiki to a local directory.
Needs a cookie for authentication and the host of the confluence instance.
THe cookie can be retrieved by opening the developer tools in the browser and copying the value of the cookie named "tenant.session.token"
The host is the host of the confluence instance, without the http/https scheme. For example, if the confluence instance is at https://confluence.example.com, the host is confluence.example.com
The output directory is where the dump will be saved. By default, it is saved in the current directory under the name "confdump"	
`,
	Run: func(cmd *cobra.Command, args []string) {
		cookie, _ := cmd.Flags().GetString("cookie")
		host, _ := cmd.Flags().GetString("host")
		output, _ := cmd.Flags().GetString("output")

		// Create the client
		confClient := client.NewConfluenceClient(host, cookie)

		// Create the saver
		saver := saver.NewFileSaver(output)

		// Create the dumper
		dumper := client.NewDumper(confClient)

		// Communication channels
		fetchChan := make(chan client.ConfluencePage)
		defer close(fetchChan)
		errChan := make(chan error)
		defer close(errChan)

		// Waitgroup to wait for everything to finish
		var wg sync.WaitGroup
		wg.Add(2)

		// Start the dumper
		dumper.StartDumper(&wg, fetchChan, errChan)

		// Start the saver
		saver.StartSaver(&wg, fetchChan, errChan)

		// Start a goroutine that prints errors
		go func() {
			for err := range errChan {
				if err != nil {
					cmd.PrintErrln(err)
				}
			}
		}()

		// Wait for the dumper to finish
		wg.Wait()

		cmd.Print("Done!")
	},
}

func Execute() {
	// Execute the root command
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().String("cookie", "", "The cookie to use for authentication")
	RootCmd.Flags().String("host", "", "The host of the confluence instance. No need to pass the http/https scheme")
	RootCmd.Flags().String("output", "./confdump", "The output directory for the dump")
	RootCmd.MarkFlagRequired("cookie")
	RootCmd.MarkFlagRequired("host")
}
