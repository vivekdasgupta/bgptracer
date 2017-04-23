package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
  "github.com/ashwanthkumar/slack-go-webhook"
)

// Direct connect path 1

var awsdc1map map[string]string
var awsip1map map[string]string

func init_awsdc1map() {
	awsdc1map = make(map[string]string)

	awsdc1map["169.254.xxx.xx"] = "AWSAccount1-DX1 Our Peer IP"
	awsdc1map["169.254.xxx.xx"] = "AWSAccount1-DX1 AWS Peer IP"
	awsdc1map["169.254.xxx.xx"] = "AWSAccount2-DX2 Our Peer IP"
	awsdc1map["169.254.xxx.xx"] = "AWSAccount2-DX2 AWS Peer IP"

}

func init_awsip1map() {
	awsip1map = make(map[string]string)

	awsip1map["aws1-our"] = "169.254.xxx.xx"
	awsip1map["aws1-peer"] = "169.254.xxx.xx"
	awsip1map["aws2-our"] = "169.254.xxx.xx"
	awsip1map["aws2-peer"] = "169.254.xxx.xx"

}

// Direct connect path 2

var awsdc2map map[string]string
var awsip2map map[string]string

func init_awsdc2map() {
	awsdc2map = make(map[string]string)

	awsdc2map["169.254.xxx.xx"] = "AWSAccount1-DX2 Our Peer IP"
	awsdc2map["169.254.xxx.xx"] = "AWSAccount1-DX2 AWS Peer IP"
	awsdc2map["169.254.xxx.xx"] = "AWSAccount2-DX1 Our Peer IP"
	awsdc2map["169.254.xxx.xx"] = "AWSAccount2-DX1 AWS Peer IP"

}

func init_awsip2map() {
	awsip2map = make(map[string]string)

	awsip2map["aws1-our"] = "169.254.xxx.xx"
	awsip2map["aws1-peer"] = "169.254.xxx.xx"
	awsip2map["aws2-our"] = "169.254.xxx.xx"
	awsip2map["aws2-peer"] = "169.254.xxx.xx"

}

func gettime() (timestr string) {
	t := time.Now()
	return t.String()
}

func main() {

	var logmsg1, logmsg2, logmsg3, logmsg4 string
	src := "aws1"
	dst := "aws2"
	success_val1 := "Host is up"
	success_val2 := "80/tcp open  http"
	bgppath := 0
	prevpath := 0

	port := os.Args[1]
	targethost := os.Args[2]
	dlogfile := os.Args[3]
	slogfile := os.Args[4]
	delay, _ := strconv.Atoi(os.Args[5])

	t := time.Now()

	filesuffix := t.Format("2006-01-02_15-04-05")

	dfileName := dlogfile + filesuffix + ".log"
	dfileHandle, errDfile := os.Create(dfileName)
	if errDfile != nil {
		log.Fatal(errDfile)
	}

	sfileName := slogfile + filesuffix + ".log"

	sfileHandle, errSfile := os.Create(sfileName)
	if errSfile != nil {
		log.Fatal(errSfile)
	}

	dwriter := bufio.NewWriter(dfileHandle)
	swriter := bufio.NewWriter(sfileHandle)

	defer dfileHandle.Close()
	defer sfileHandle.Close()

	init_awsdc1map()
	init_awsip1map()

	init_awsdc2map()
	init_awsip2map()

	log.Println(" ****** BGP tracer initialized ...  ******  ", gettime())
	log.Println("Summary Log file  ", slogfile)
	log.Println("Detailed Log file    ", dlogfile)
	log.Println("Port:  ", port, " Target Host: ", targethost, " Delay in seconds : ", delay)

	fmt.Fprintln(swriter, "****** BGP tracer initialized ...  ******  \n\n")
	fmt.Fprintln(dwriter, "****** BGP tracer initialized ...  ******  \n\n")

	// main loop

	for { // infinite loop

		log.Printf("\n\n")
		fmt.Fprintln(dwriter, " ******   ****** \n\n")

		cmndOut, err := exec.Command("nmap", "-p", port, "--traceroute", targethost).Output()
		if err != nil {
			log.Println("Error :", err)
		}
		output := string(cmndOut)

		if strings.Contains(output, success_val1) || strings.Contains(output, success_val2) {

			if strings.Contains(output, awsip1map["aws2-peer"]) {

				bgppath = 1

				logmsg1 = " Network UP : From -> " + src + " to " + dst + " Using: " + awsdc1map[awsip1map["aws2-peer"]] + " IP: " + awsip1map["aws2-peer"]
				log.Println(logmsg1)
				slackpost(logmsg1)
				fmt.Fprintln(swriter, gettime()+logmsg1)
				fmt.Fprintln(dwriter, gettime()+logmsg1)
				fmt.Fprintln(dwriter, gettime()+" BGP Endpoints -> ", awsdc1map[awsip1map["aws2-our"]], " to ", awsdc1map[awsip1map["aws2-peer"]])
				fmt.Fprintln(dwriter, output)
				if prevpath != bgppath && prevpath != 0 {
					log.Println(" BGP path change detected... For -> ", src, " to ", dst)
					slackpost(" BGP path change detected... ")
					fmt.Fprintln(swriter, gettime()+" BGP path change detected... ")
					fmt.Fprintln(dwriter, gettime()+" BGP path change detected... ")

				}
				prevpath = bgppath

			} else if strings.Contains(output, awsip2map["aws2-peer"]) {

				bgppath = 2
				logmsg2 = " Network UP : From -> " + src + " to " + dst + " Using: " + awsdc2map[awsip2map["aws2-peer"]] + " IP: " + awsip2map["aws2-peer"]
				log.Println(logmsg2)
				slackpost(logmsg2)
				fmt.Fprintln(swriter, gettime()+logmsg2)
				fmt.Fprintln(dwriter, gettime()+logmsg2)
				fmt.Fprintln(dwriter, gettime()+" BGP Endpoints -> ", awsdc2map[awsip2map["aws2-our"]], " to ", awsdc2map[awsip2map["aws2-peer"]])
				fmt.Fprintln(dwriter, output)
				if prevpath != bgppath && prevpath != 0 {
					log.Println(" BGP path change detected... For -> ", src, " to ", dst)
					slackpost(" BGP path change detected... ")
					fmt.Fprintln(swriter, gettime()+" BGP path change detected... ")
					fmt.Fprintln(dwriter, gettime()+" BGP path change detected... ")

				}
				prevpath = bgppath

			} else {
				logmsg3 = " Network connectivity up: Path Unknown/New : From -> " + src + " to " + dst
				slackpost(logmsg3)
				log.Println(" Network connectivity up: Path Unknown/New : From -> ", src, " to ", dst)
				fmt.Fprintln(swriter, gettime()+" Network connectivity up: Path Unknown/New : From -> ", src, " to ", dst)
				fmt.Fprintln(dwriter, output)
			}

		} else {
			logmsg4 = " Network connectivity down: From -> " + src + " to " + dst
			slackpost(logmsg4)
			log.Println(" Network connectivity down: From -> ", src, " to ", dst)
			fmt.Fprintln(swriter, gettime()+" Network connectivity down : From -> ", src, " to ", dst)
			fmt.Fprintln(dwriter, output)
		}

		dwriter.Flush()
		swriter.Flush()

		log.Println(" Wait a few seconds before next try ... ")
		log.Printf("\n\n")
		fmt.Fprintln(dwriter, " Wait a few seconds before next try ... \n\n")
		time.Sleep(time.Duration(delay) * 1000 * time.Millisecond) // wait delay secs

	} // end for (infinite)
	// main loop end

}

func slackpost(text string) {
	webhookUrl := "https://hooks.slack.com/services/XXXXXXX/XXXXXXX/UseYourOwnWebHookURL"

	attachment1 := slack.Attachment{}
	//attachment1.AddField(slack.Field{Title: "Author", Value: "Vivek Dasgupta"}).AddField(slack.Field{Title: "Status", Value: "Completed"})
	payload := slack.Payload{
		Text:        text,
		Username:    "bgptracer-bot",
		Channel:     "#alerts",
		IconEmoji:   ":bird:",
		Attachments: []slack.Attachment{attachment1},
	}
	err := slack.Send(webhookUrl, "", payload)
	if len(err) > 0 {
		fmt.Printf("error: %s\n", err)
	}
}
