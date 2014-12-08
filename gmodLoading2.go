// Copyright 2014 Johan "tazaar" Englund
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Flag declarations
var port_Ptr = flag.String("port", "8080", "Specify port to listen on")
var ip_Ptr = flag.String("ip", "0.0.0.0", "Specify IP to listen on")
var apiKey_Ptr = flag.String("apikey", "", "Get your own at http://steamcommunity.com/dev/apikey")
var help_Ptr = flag.Bool("help", false, "Show help")
var log_Ptr = flag.Bool("log", false, "Log to file")

// Steam's GetPlayerSummaries(v0002) api struct
type steam_GetPlayerSummaries struct {
	Response struct {
		Players []struct {
			Steamid                  string `json:"steamid"`
			Communityvisibilitystate int    `json:"communityvisibilitystate"`
			Profilestate             int    `json:"profilestate"`
			Personaname              string `json:"personaname"`
			Lastlogoff               int    `json:"lastlogoff"`
			Profileurl               string `json:"profileurl"`
			Avatar                   string `json:"avatar"`
			Avatarmedium             string `json:"avatarmedium"`
			Avatarfull               string `json:"avatarfull"`
			Personastate             int    `json:"personastate"`
			Realname                 string `json:"realname"`
			Primaryclanid            string `json:"primaryclanid"`
			Timecreated              int    `json:"timecreated"`
			Personastateflags        int    `json:"personastateflags"`
			Loccountrycode           string `json:"loccountrycode"`
			Locstatecode             string `json:"locstatecode"`
			Loccityid                int    `json:"loccityid"`
		} `json:"players"`
	} `json:"response"`
}

// loading_h should be seen when joining/loading the server.
// sv_loadingurl "http://gmod.example/loading.html?mapname=%m&steamid=%s"
// Steam may not respond, A better solution could be AJAX'ing needed data
func loading_h(w http.ResponseWriter, r *http.Request) {
	// TODO; Implement gmod javascript callbacks
	log.Println("serving loading.html")

	// Struct to be passed to template
	type BasicInfo struct {
		MapName    string
		PlayerName string
	}
	// Set some fake data just in case
	b := BasicInfo{
		"de_myst3ry",
		"Anon",
	}

	// Fetch player summary & populate BasicInfo
	g := r.URL.Query()
	// Make the steam api call
	url := "http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=" + *apiKey_Ptr + "&format=json&steamids=" + g["steamid"][0]
	res, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		log.Println("Could not reach Steam servers!")
	} else {
		defer res.Body.Close()
		// Unmarchal json into a steam_GetPlayerSummaries struct
		jsonData, err := ioutil.ReadAll(res.Body)
		var d steam_GetPlayerSummaries
		err = json.Unmarshal(jsonData, &d)
		if err != nil {
			log.Println(err.Error())
		}
		// Populate template data
		b.PlayerName = d.Response.Players[0].Personaname
		b.MapName = g["mapname"][0]
	}

	// Combine template with data and send to client
	t := template.New("loading")
	t, err = template.ParseFiles("templates/loading.tmpl")
	if err != nil {
		log.Println(err.Error())
	}
	err = t.Execute(w, b)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("done")

}

// ulxMotd_h should be seen when joining an ULX server.
func ulxMotd_h(w http.ResponseWriter, r *http.Request) {
	log.Println("serving motd.html")
	// Combine template with data and send to client
	t := template.New("motd")
	t, err := template.ParseFiles("templates/motd.tmpl")
	if err != nil {
		log.Println(err.Error())
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("done")
}

func main() {
	// Parse flags and respond to a few
	flag.Parse()
	if *help_Ptr {
		fmt.Println("Default settings")
		fmt.Println("gmodLoading2.exe -port 8080 -ip 0.0.0.0 -apikey xxx")
		os.Exit(0)
	}
	if *log_Ptr {
		LogFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", ":", err)
		}
		log.SetOutput(LogFile)
	}
	if *apiKey_Ptr == "" {
		log.Println("API Key empty! get one at http://steamcommunity.com/dev/apikey")
		os.Exit(1)
	}

	log.Println("Listening at " + *ip_Ptr + ":" + *port_Ptr)
	// HTTP Handlers
	http.HandleFunc("/loading.html", loading_h)
	http.HandleFunc("/motd.html", ulxMotd_h)
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(*ip_Ptr+":"+*port_Ptr, nil)
}
