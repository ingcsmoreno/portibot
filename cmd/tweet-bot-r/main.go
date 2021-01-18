package main

import (
	"fmt"
	"log"
	"os"
    "os/signal"
    "syscall"
    "math/rand"
    "time"

    "github.com/dghubble/go-twitter/twitter"
    "github.com/dghubble/oauth1"
)

var (
    version   string // version number
    sha1ver   string // sha1 revision used to build the program
    buildTime string // when the executable was built
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

    log.Printf("Logged in as User: %+v\n", user.Name)
    return client, nil
}

func generateTweetAnswer (user string) string {
    rand.Seed(time.Now().UnixNano())
    
    initials := []string{
        "Hola @%s! Listo para algo de sci-fi?", 
        "Hola @%s! Que tal todo? Acá va algo de sci-fi", 
        "Claro que si @%s!", 
        "A la orden @%s!", 
        "Hola @%s, que te pare esto?", 
        "@%s que bueno que preguntas.",
    }
    
    recommendations := []string{
        "El monstruo de sci-fi mundialmente reconocido, Frankenstein de Mary Shelley, libraso. Si no lo leiste, recomendadísimo.",
        "Viste 'Yo, Robot'? entonces el libro Foundación de Isaac Asimov te va a fascinar. Nota de color: pronto se va a estrenar una serie al respecto.",
        "Solaris de Stanislaw Lem, es otra de esas obras que no pueden faltar en una biblioteca sci-fi. Tiene un par de adaptaciones al cine también.",
        "Dune de Frank Herbert, es una novela espectacular, de la cual pronto estrenará una nueva película. Recomendadísima",
        "Te gustan las novelas futuristas? En Neuromancer, de William Gibson, se enfrentan hackers contra una inteligencia artificial... es todo lo que voy a decir",
        "El problema de los tres cuerpos, de Liu Cixin, cuenta la historia de una civilización luchando por sobrevivir al sistema planetario en el que viven. que tal?",
        "El Marciano, de Andy Weir, es una obra genial sobre un astronauta que queda barado en Marte. Su adaptación al cine también fue muy buena!",
    }
    
    return fmt.Sprintf(initials[rand.Intn(len(initials))] + "\n\n" + recommendations[rand.Intn(len(recommendations))], user)
}

func main() {
    log.Printf("Initiating Tweet-bot (Recomm) %s", version)
    log.Printf(" * Commit: %s", sha1ver)
    log.Printf(" * Build Date: %s", buildTime)
    log.Printf("Signing in to Twitter.")
    creds := Credentials{
        AccessToken:       os.Getenv("ACCESS_TOKEN"),
        AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
        ConsumerKey:       os.Getenv("CONSUMER_KEY"),
        ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
    }

    //fmt.Printf("%+v\n", creds)

    client, err := getClient(&creds)
    if err != nil {
        log.Println("Error getting Twitter Client")
        log.Println(err)
    }
    
    // Verify Credentials
    verifyParams := &twitter.AccountVerifyParams{
        SkipStatus:   twitter.Bool(true),
        IncludeEmail: twitter.Bool(true),
    }
    user, _, err := client.Accounts.VerifyCredentials(verifyParams)
    if err != nil {
        log.Printf("Login error %v", err)
    }

    // Print out the pointer to our client
    // for now so it doesn't throw errors
    
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
        // Avoid processing own tweets
        if ( tweet.User.ScreenName == user.ScreenName) { return }

        // Received tweet info
	    log.Println("-----------------")
	    log.Printf("Tweet ID: %d\n", tweet.ID)
	    log.Printf("User: %s\n", tweet.User.ScreenName)
	    log.Printf("Tweet Text: %s\n", tweet.Text)

        // Tweet response text
		answer := generateTweetAnswer(tweet.User.ScreenName)
		log.Printf("Tweet Answer: %s\n", answer)
	    
        // Responding tweet
	    tweetParams := &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID}
        _, resp, err := client.Statuses.Update(answer, tweetParams)
	    if err != nil {
		    log.Printf("Statuses.Tweet error %v", err)
		}
		log.Printf("Tweet Status Code: %d\n", resp.StatusCode)

        // Retweeting the original tweet
        retweetParams := &twitter.StatusRetweetParams{TrimUser: twitter.Bool(true)}
        _, retweetResponse, err := client.Statuses.Retweet(tweet.ID, retweetParams)
        if err != nil {
            log.Printf("Statuses.Retweet error %v", err)
        }
        log.Printf("Retweet Status Code: %d\n\n", retweetResponse.StatusCode)

	}

	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	
	stream.Stop()

}