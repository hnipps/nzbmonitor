package cmd

import (
	"bufio"
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [command]",
	Short: "Send a command to the service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sendCommand(args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func sendCommand(command string) {
	conn, err := net.Dial("tcp", "localhost:6666")
	if err != nil {
		fmt.Println("Error connecting to service:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s\n", command)
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Print("Response from service: ", response)
}
