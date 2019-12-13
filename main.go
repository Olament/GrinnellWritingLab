package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gen2brain/beeep"
	"golang.org/x/net/publicsuffix"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func main() {
	/* constant */
	const LOGIN_PAGE string = "https://writinglab.sites.grinnell.edu/schedule/Web/index.php"
	const SCHEDULE_PAGE string = "https://writinglab.sites.grinnell.edu/schedule/Web/schedule.php"
	const USER_NAME string = "YOUR_EMAIL"
	const PASSWORD string = "PASSWORD"
	const INTERVAL time.Duration = 1

	/* initilize cookie jar and http client */
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// get the php session id from website
	//throw away the response since we just want the session id
	_, err = client.Get(LOGIN_PAGE)
	if err != nil {
		log.Fatal(err)
	}

	login(client, LOGIN_PAGE, USER_NAME, PASSWORD) // get login cookie

	/* initialize ticker */
	ticker := time.NewTicker(INTERVAL * time.Minute)

	/* instant first tick */
	for ; true; <-ticker.C {
		fmt.Print(time.Now().Format("2006-Jan-02 15:04:05 "))
		s, err := getSchedule(client, SCHEDULE_PAGE, true)
		if err != nil {
			fmt.Println(err)
		}

		/* only send alert when ... */
		if len(s) > 0 {
			/* send system notification */
			err := beeep.Alert("Writing Lab", "Avaliable slot found", "assets/icon.jpg")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("%v spot(s) found", len(s))
			j, _ := json.MarshalIndent(s, "", "	")
			fmt.Println(string(j))
		} else {
			fmt.Println("No avaliable spot")
		}
	}
}

func login(client *http.Client, login_url string, user_name string, password string) (err error) {
	//TODO better login failure check
	_, err = client.PostForm(login_url, url.Values{
		"email":    {user_name},
		"password": {password},
		"captcha":  {},
		"login":    {"submit"},
		"resume":   {},
		"language": {"en_us"},
	})
	if err != nil {
		return err
	}

	return nil
}

type ReservationStatus struct {
	Date         string
	Time         string
	Name         string
	Availability bool
}

func getSchedule(client *http.Client, schedule_url string, filter bool) (statuses []ReservationStatus, err error) {
	res, err := client.Get(schedule_url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// it's called reservable
	// find reservation status of each day
	doc.Find(".reservations").Each(func(dateIndex int, date *goquery.Selection) {
		resDate := date.Find(".resdate").Text()
		var times []string
		date.Find(".reslabel").Each(func(i int, selection *goquery.Selection) {
			times = append(times, selection.Text())
		})

		// find reservation status of each tutor/mentor
		date.Find(".slots").Each(func(mentorIndex int, mentor *goquery.Selection) {
			name := mentor.Find(".resourceNameSelector").Text()
			// find status of each slot
			mentor.Find(".slot").Each(func(slotIndex int, slot *goquery.Selection) {
				status, _ := slot.Attr("class")

				var availability bool = false
				for _, stat := range strings.Split(status, " ") {
					if stat == "reservable" {
						availability = true
						break
					}
				}

				currReservationStatus := ReservationStatus{
					Date:         resDate,
					Time:         times[slotIndex],
					Name:         name,
					Availability: availability,
				}

				/* filter the result */
				if filter {
					if currReservationStatus.Availability {
						statuses = append(statuses, currReservationStatus)
					}
				} else {
					statuses = append(statuses, currReservationStatus)
				}
			})
		})
	})

	return statuses, nil
}
