package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/CyCoreSystems/dg-whitelist/list"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd, addCmd, delCmd)

	rootCmd.PersistentFlags().String("list", "white", "list on which to operate: one of black, white, or grey")
	rootCmd.PersistentFlags().String("s", "http://localhost:3000", "Service root URL")
}

func main() {
	rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "dgctl",
	Short: "dgctl is a simple control tool for the dg-whitelist service",
	Long:  `A simple CLI adapter for the HTTP API of the dg-whitelist service.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		return
	},
}

func verifyList(listName string, err error) (string, error) {
	if err != nil {
		return "", err
	}

	switch listName {
	case list.ListBlack:
	case list.ListGrey:
	case list.ListWhite:
	default:
		return "", fmt.Errorf("unknown list type %s", listName)
	}

	return listName, nil
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get the addresses from a list",
	Run: func(cmd *cobra.Command, args []string) {

		out := []list.Item{}

		svc, err := cmd.Flags().GetString("s")
		if err != nil {
			log.Fatalln("failed to parse service url")
		}

		listName, err := verifyList(cmd.Flags().GetString("list"))
		if err != nil {
			log.Fatalln("failed to parse list name")
		}

		resp, err := http.Get(fmt.Sprintf("%s/%s", svc, listName))
		if err != nil {
			log.Fatalln("failed to get list from server:", err)
		}

		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			log.Fatalln("failed to parse server response:", err)
		}

		for _, e := range out {
			fmt.Println(e.Address)
		}

		return
	},
}

var addCmd = &cobra.Command{
	Use:   "add <address> [<address>...]",
	Short: "add an address to a list",
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := cmd.Flags().GetString("s")
		if err != nil {
			log.Fatalln("failed to parse service url")
		}

		listName, err := verifyList(cmd.Flags().GetString("list"))
		if err != nil {
			log.Fatalln("failed to parse list name")
		}

		for _, addr := range args {
			buf := new(bytes.Buffer)

			if err := json.NewEncoder(buf).Encode(&list.Item{
				List:    listName,
				Address: addr,
			}); err != nil {
				log.Fatalln("failed to encode address for addition:", err)
			}

			resp, err := http.Post(svc, "application/json", buf)
			if err != nil {
				log.Fatalln("failed to make request to server:", err)
			}

			if resp.StatusCode > 299 {
				log.Fatalf("failed to add %q: %s", addr, resp.Status)
			}
		}
	},
}

var delCmd = &cobra.Command{
	Use:   "del <address> [<address>...]",
	Short: "remove an address from a list",
	Run: func(cmd *cobra.Command, args []string) {

		svc, err := cmd.Flags().GetString("s")
		if err != nil {
			log.Fatalln("failed to parse service url")
		}

		listName, err := verifyList(cmd.Flags().GetString("list"))
		if err != nil {
			log.Fatalln("failed to parse list name")
		}

		for _, addr := range args {
			buf := new(bytes.Buffer)

			if err := json.NewEncoder(buf).Encode(&list.Item{
				List:    listName,
				Address: addr,
			}); err != nil {
				log.Fatalln("failed to encode address for deletion:", err)
			}

			req, err := http.NewRequest(http.MethodDelete, svc, buf)
			if err != nil {
				log.Fatalln("failed to construct delete request:", err)
			}
			req.Header.Add("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatalln("failed to make request to server:", err)
			}

			if resp.StatusCode > 299 {
				log.Fatalf("failed to delete %q: %s", addr, resp.Status)
			}
		}
	},
}
