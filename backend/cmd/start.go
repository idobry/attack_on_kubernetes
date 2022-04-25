/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/attack_on_kubernetes/controller"
)

var (
	//NAMESPACE = os.Getenv("NAMESPACE") //man
	HOST = os.Getenv("HOST") //man
	IMAGE = os.Getenv("IMAGE") //man
	SSHPASS = os.Getenv("SSHPASS") //man
	SSHUSER = os.Getenv("SSHUSER") //man
	INGRESS_CLASS = os.Getenv("INGRESS_CLASS") //man
	SAN = os.Getenv("SERVICE_ACCOUNT_NAME") //optional 
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		controller.NAMESPACE, _ = cmd.Flags().GetString("namespace")
		controller.HOST, _ = cmd.Flags().GetString("host")
		controller.IMAGE, _ = cmd.Flags().GetString("image")
		controller.SSHPASS, _ = cmd.Flags().GetString("sshpass")
		controller.SSHUSER, _ = cmd.Flags().GetString("sshuser")
		controller.INGRESS_CLASS, _ = cmd.Flags().GetString("ingress-class")
		controller.SAN, _ = cmd.Flags().GetString("service-account")

		serve()
	},
}

func serve() {
	r := mux.NewRouter()
	r.HandleFunc("/create", controller.CreateNewWetty).Methods("GET")
	r.HandleFunc("/delete", controller.DeleteWetty).Methods("POST")
	r.HandleFunc("/health-check", HealthCheck).Methods("GET")
	fmt.Println("ListenAndServe...")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().String("namespace","default", "A help for namespace")

	startCmd.PersistentFlags().String("host","", "A help for host")
	startCmd.MarkPersistentFlagRequired("host")

	startCmd.PersistentFlags().String("image","idobry/k8tty_backend", "A help for image")

	startCmd.PersistentFlags().String("sshuser","k8tty", "A help for sshuser")

	startCmd.PersistentFlags().String("sshpass","k8tty", "A help for sshuser")

	startCmd.PersistentFlags().String("ingress-class", "", "A help for foo")
	startCmd.MarkPersistentFlagRequired("ingress-class")

	startCmd.PersistentFlags().String("service-account", "default", "A help for foo")
}
