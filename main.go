package main

import (
	"github.com/levelabs/level-go/cmd"
)

func main() {
	cmd.Execute()
}

// func main() {
// 	app = NewApp()
//
// 	http.HandleFunc("/api", apiHandler)
//
// 	NewScheduler()
//
// 	fmt.Printf("Starting server at port 8080\n")
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatal(err)
// 	}
// }
