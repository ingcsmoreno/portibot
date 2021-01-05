package main

import (
	"fmt"
	"log"
	"os"
    "github.com/dghubble/go-twitter/twitter"
    "github.com/dghubble/oauth1"
    "os/signal"
    "syscall"
)

// Credentials stores all of our access/consumer tokens
// and secret keys needed for authentication against
// the twitter REST API.
type Credentials struct {
    ConsumerKey       string
    ConsumerSecret    string
    AccessToken       string
    AccessTokenSecret string
}

// getClient is a helper function that will return a twitter client
// that we can subsequently use to send tweets, or to stream new tweets
// this will take in a pointer to a Credential struct which will contain
// everything needed to authenticate and return a pointer to a twitter Client
// or an error
func getClient(creds *Credentials) (*twitter.Client, error) {
    // Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
    config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
    // Pass in your Access Token and your Access Token Secret
    token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

    httpClient := config.Client(oauth1.NoContext, token)
    client := twitter.NewClient(httpClient)

    // Verify Credentials
    verifyParams := &twitter.AccountVerifyParams{
        SkipStatus:   twitter.Bool(true),
        IncludeEmail: twitter.Bool(true),
    }

    // we can retrieve the user and verify if the credentials
    // we have used successfully allow us to log in!
    user, _, err := client.Accounts.VerifyCredentials(verifyParams)
    if err != nil {
        return nil, err
    }

    log.Printf("User's ACCOUNT:\n%+v\n", user)
    return client, nil
}

func main() {
    fmt.Println("Go-Twitter Bot v0.01")
    creds := Credentials{
        AccessToken:       os.Getenv("ACCESS_TOKEN"),
        AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
        ConsumerKey:       os.Getenv("CONSUMER_KEY"),
        ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
    }

    fmt.Printf("%+v\n", creds)

    client, err := getClient(&creds)
    if err != nil {
        log.Println("Error getting Twitter Client")
        log.Println(err)
    }

     Print out the pointer to our client
     for now so it doesn't throw errors
    params := &twitter.StreamFilterParams{
	    Track: []string{"#tiramescifi"},
	    StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(params)

	//params := &twitter.StreamSampleParams{
	//    StallWarnings: twitter.Bool(true),
	//}
	//stream, err := client.Streams.Sample(params)

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
	    fmt.Println("-----------------")
	    fmt.Println(tweet.ID)
	    fmt.Println(tweet.User.ScreenName)
	    fmt.Println(tweet.Text)

	    tweetParams := &twitter.StatusUpdateParams{
		    InReplyToStatusID: tweet.ID,
		}
		status := fmt.Sprintf("Hola @%s! listo para algo de buena ciencia ficción? (Probando de nuevo)", tweet.User.ScreenName)
		//log.Println(status)
	    tweet, resp, err := client.Statuses.Update(status, tweetParams)
	    if err != nil {
		    log.Println(err)
		}
		log.Printf("%+v\n", resp)

	}
	demux.DM = func(dm *twitter.DirectMessage) {
	    fmt.Println(dm.SenderID)
	}

	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	
	stream.Stop()

}