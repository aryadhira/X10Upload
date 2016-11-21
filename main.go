package main

import (
	. "eaciit/x10consoleapps/x10upload/helpers"
	"strings"

	"time"

	tk "github.com/eaciit/toolkit"
	"github.com/howeyc/fsnotify"
)

func main() {
	config := ReadConfig()
	processcom := config["processcompany"]
	failedcom := config["failedcompany"]
	successcom := config["successcompany"]
	inboxcom := config["inboxcompany"]
	processind := config["processindividual"]
	failedind := config["failedindividual"]
	successind := config["successindividual"]
	inboxind := config["inboxindividual"]
	webapps := config["webapps"]
	webapps2 := config["webapps2"]

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		tk.Println(err.Error())
	}

	watcher2, err2 := fsnotify.NewWatcher()
	if err2 != nil {
		tk.Println(err2.Error())
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				strEV := strings.Split(""+ev.String(), ":")
				action := strings.Trim(strEV[len(strEV)-1], " ")
				if action == "CREATE" {
					time.Sleep(5 * time.Second)
					ProcessFile(inboxcom, processcom, failedcom, successcom, "Company", webapps2)
				}
			case err := <-watcher.Error:
				tk.Println("error:", err)
			case ev2 := <-watcher2.Event:
				strEV2 := strings.Split(""+ev2.String(), ":")
				action2 := strings.Trim(strEV2[len(strEV2)-1], " ")
				if action2 == "CREATE" {
					time.Sleep(5 * time.Second)
					ProcessFile(inboxind, processind, failedind, successind, "Individual", webapps)
				}
			case err2 := <-watcher.Error:
				tk.Println("error:", err2)
			}

		}
	}()

	err = watcher.Watch(inboxcom)
	if err != nil {
		tk.Println(err.Error())
	}

	err2 = watcher2.Watch(inboxind)
	if err2 != nil {
		tk.Println(err2.Error())
	}

	tk.Println("Watcher Started...")

	<-done

	watcher.Close()
}
