// Package gmodLoading2 is a webserver for Garry's mod server owners,
// It serves a customizable Bootstrap loading page and ULX motd page.
//
// Useage (example script found in tools/start.bat)
//		gmodLoading2.exe -ip 0.0.0.0 -port 8080 -apikey XXX
// Don't forget to add this line to server.cfg!
//		sv_loadingurl "http://gmod.example/loading.html?mapname=%m&steamid=%s
// And for ULX add this to data/ulx/config.txt
//		set showMotd "http://gmod.example/motd.html"
//
// It also makes use of Valve's Steam Web API,
// So you need your own Steam API key!
//		http://steamcommunity.com/dev/apikey
package main
