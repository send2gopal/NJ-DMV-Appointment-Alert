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
func appointmentBot(p string, ch chan bool) {
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		result := strings.Contains(string(r.Body[:]), "Next Available")
		if result {
			playSound()
		}
		ch <- result
	})

	//Command to visit the website
	c.Visit(p)
}

func main() {
	if len(os.Args) < 2 {
		panic("Need appointment type")
	}
	quit := make(chan bool)

	i := 0
	url := "https://telegov.njportal.com/njmvc/AppointmentWizard/" + os.Args[1]
	fmt.Printf("Call to: %v", url)
	for {
		go appointmentBot(url, quit)
		r := <-quit
		if r {
			fmt.Print(r)
			return
		}
		fmt.Printf("\nLoop %v ", i)
		time.Sleep(120 * time.Second)
		i++
	}
}
