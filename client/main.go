// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	pb "github.com/paulvollmer/robotstxt-datastore/proto/robotstxt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	var flagGrpcServerAddr string

	var rootCmd = &cobra.Command{Use: "robotstxt"}
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&flagGrpcServerAddr, "grpc", "localhost:5000", "the grpc server address")

	var flagCheckRefresh bool
	cmdCheck := &cobra.Command{
		Use:   "check",
		Short: "Check if the robotstxt data of the given host exist at the database",
		Long:  `The check command will query the database and check if the data is older than the refresh_after setting.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("missing uri")
				os.Exit(2)
			}
			client := initGrpcClient(flagGrpcServerAddr)
			data, err := client.CheckRobotstxt(context.Background(), &pb.Check{Host: args[0], Refresh: flagCheckRefresh})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			printData(data)
		},
	}
	cmdCheck.Flags().BoolVarP(&flagCheckRefresh, "refresh", "r", false, "refresh the database entry")

	cmdGet := &cobra.Command{
		Use:   "get",
		Short: "Get the robotstxt data of the given host",
		Long:  `The get command will query the database and print the data to stdout.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("missing uri")
				os.Exit(2)
			}
			client := initGrpcClient(flagGrpcServerAddr)
			data, err := client.GetRobotstxt(context.Background(), &pb.GetRobotstxtRequest{Host: args[0]})
			if err != nil {
				fmt.Printf("check error: %v\n", err)
				os.Exit(1)
			}
			printData(data)
		},
	}

	var flagListOffset int64
	var flagListLimit int32
	var flagListSearch string
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "Get a list of robotstxt beginning at the given cursor",
		Long:  `The list command will query the database at the begin of the given limit/offset and print the data to stdout.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := initGrpcClient(flagGrpcServerAddr)
			data, err := client.ListRobotstxts(context.Background(), &pb.ListRobotstxtsRequest{Limit: flagListLimit, Offset: flagListOffset, Search: flagListSearch})
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			fmt.Printf("Total: %v\n\n", data.TotalCount)

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
			fmt.Fprintln(w, "Scheme\t Host\t StatusCode\t Response Time\t Response URL\t Total Sitemaps\t Total Rules")
			fmt.Fprintln(w, "------\t ----\t ----------\t -------------\t ------------\t --------------\t -----------")
			for i := 0; i < len(data.Items); i++ {
				fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\t %v\t %v\n", data.Items[i].Scheme, data.Items[i].Host, data.Items[i].Statuscode, data.Items[i].ResponseTime, data.Items[i].ResponseUrl, len(data.Items[i].Sitemaps), len(data.Items[i].Rules))
			}
			w.Flush()
		},
	}
	cmdList.Flags().Int32VarP(&flagListLimit, "limit", "l", 25, "the limit")
	cmdList.Flags().Int64VarP(&flagListOffset, "offset", "o", 0, "the offset")
	cmdList.Flags().StringVarP(&flagListSearch, "search", "s", "", "the search string")

	rootCmd.AddCommand(cmdCheck, cmdGet, cmdList)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// read in environment variables that match
	viper.AutomaticEnv()
}

func initGrpcClient(addr string) pb.RobotstxtServiceClient {
	// Set up a connection to the server.
	fmt.Println("GRPS ADDR:", addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		os.Exit(1)
	}
	//defer conn.Close()
	client := pb.NewRobotstxtServiceClient(conn)
	return client
}

func printData(data *pb.Robotstxt) {
	fmt.Printf("Scheme       : %v\n", data.Scheme)
	fmt.Printf("Host         : %v\n", data.Host)
	fmt.Printf("StatusCode   : %v\n", data.Statuscode)
	fmt.Printf("ResponseTime : %v\n", data.ResponseTime)
	fmt.Printf("ResponseURL  : %v\n", data.ResponseUrl)
	fmt.Printf("\nSitemaps     : (total %v)\n", len(data.Sitemaps))
	for i := 0; i < len(data.Sitemaps); i++ {
		fmt.Printf("- %v\n", data.Sitemaps[i])
	}
	//fmt.Printf("\nRules     : (total %v)\n", len(data.Rules))
	//for i := 0; i < len(data.Rules); i++ {
	//	fmt.Printf("- Agent: %v \tPaths: %v\n", data.Rules[i].Agent, data.Rules[i].Paths)
	//}
}
