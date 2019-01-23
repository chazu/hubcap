package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gdamore/tcell"
	"github.com/google/go-github/github"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
)

func writeDefaultConfig(username *string) {
	configStr := []byte(fmt.Sprintf("---\nusername: %s\n...", *username))
	err := ioutil.WriteFile(fmt.Sprintf("%s/.hubcap.yaml", os.Getenv("HOME")), configStr, 0666)
	if err != nil {
		panic(fmt.Errorf("Fatal: %s \n", err))
	}
}

func main() {
	// Manage config

	viper.SetConfigName(".hubcap")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(os.Getenv("HOME"))

	_, err := os.Stat(fmt.
		Sprintf("%s/.hubcap.yaml", os.Getenv("HOME")))

	if os.IsNotExist(err) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Whats your github username: ")
		text, _ := reader.ReadString('\n')
		writeDefaultConfig(&text)
	}

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	client := github.NewClient(nil)
	options := &github.ActivityListStarredOptions{Sort: "created"}

	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false)

	starred, _, err := client.Activity.ListStarred(context.Background(), viper.Get("username").(string), options)

	if err != nil {
		fmt.Printf("ERR: %s", err.Error())
	}
	for i, star := range starred {
		table.SetCell(i, 0, &tview.TableCell{
			Text:  *star.Repository.Name,
			Color: tcell.ColorWhite,
		})
		table.SetCell(i, 1, &tview.TableCell{
			Text:  *star.Repository.Description,
			Color: tcell.ColorWhite,
		})
		table.SetCell(i, 2, &tview.TableCell{
			Text:  *star.Repository.Language,
			Color: tcell.ColorWhite,
		})
	}

	if err := tview.NewApplication().SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
