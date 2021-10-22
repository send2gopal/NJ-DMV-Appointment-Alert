package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func playSound() {
	f, err := os.Open("./appointment.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
func appointmentBot(p string) bool {
	c := colly.NewCollector()
	returnval := false
	c.OnResponse(func(r *colly.Response) {
		returnval = strings.Contains(string(r.Body[:]), "Next Available")
	})

	//Command to visit the website
	c.Visit("https://telegov.njportal.com/njmvc/AppointmentWizard/" + p)

	return returnval
}

func main() {
	if len(os.Args) < 2 {
		panic("Need appointment type")
	}
	ticker := time.NewTicker(300 * time.Second)
	quit := make(chan bool)

	go func() {

		for {
			select {
			case <-ticker.C:
				r := appointmentBot(os.Args[1])
				if r {
					fmt.Print(r)
					playSound()
					quit <- r
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-quit
}
