package main

import (
    "fmt"
    "io/ioutil"
	"log"
	"os"
    "os/signal"
    "strconv"
    "syscall"
    "math/rand"
    "net/http"
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
func getClient(creds *Credentials) (*twitter.Client, *twitter.User, error) {
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
        return nil, nil, err
    }

    log.Printf("Logged in as User: %+v\n", user.Name)
    return client, user, nil
}

func generateTweetAnswer (dbAcc DBAccess, user string) (string, string)  {
    rand.Seed(time.Now().UnixNano())
    
    initials := []string{
        "Hola @%s! Listo para algo de sci-fi?", 
        "Hola @%s! Que tal todo? Acá va algo de sci-fi", 
        "Claro que si @%s!", 
        "A la orden @%s!", 
        "Hola @%s", 
        "@%s que bueno que preguntas.",
        "Hay tanto de donde elegir @%s",
        "@%s pide, Pórtico contesta",
        "Pediste Sci-Fi @%s? No se diga mas",
    }

    middle := []string{
        "que te parece %s (%s) de %s?", 
        "%s (%s) de %s es realmente genial!", 
        "definitivamente %s (%s) de %s es de esas obras infaltables", 
        "%s (%s) de %s, totalmente recomendable", 
        "que tal %s (%s) de %s? Si no está en tu repertorio, debería.", 
        "segurísimo que mas de uno te recomendaría %s (%s) de %s, no vamos a ser la excepción XD",
        "%s (%s) de %s es de esas obras que no pueden faltar",
        "te recomendamos %s (%s) de %s, es sobre... tiene eso que... en fin, te va a encantar.",
    }
    
    //recommendations := []string{
    //    "El monstruo de sci-fi mundialmente reconocido, Frankenstein de Mary Shelley, libraso. Si no lo leiste, recomendadísimo.",
    //    "Viste 'Yo, Robot'? entonces el libro Foundación de Isaac Asimov te va a fascinar. Nota de color: pronto se va a estrenar una serie al respecto.",
    //    "Solaris de Stanislaw Lem, es otra de esas obras que no pueden faltar en una biblioteca sci-fi. Tiene un par de adaptaciones al cine también.",
    //    "Dune de Frank Herbert, es una novela espectacular, de la cual pronto estrenará una nueva película. Recomendadísima",
    //    "Te gustan las novelas futuristas? En Neuromancer, de William Gibson, se enfrentan hackers contra una inteligencia artificial... es todo lo que voy a decir",
    //    "El problema de los tres cuerpos, de Liu Cixin, cuenta la historia de una civilización luchando por sobrevivir al sistema planetario en el que viven. que tal?",
    //    "El Marciano, de Andy Weir, es una obra genial sobre un astronauta que queda barado en Marte. Su adaptación al cine también fue muy buena!",
    //}

    result, _, _ := getRandomBook(dbAcc)
    
    message := fmt.Sprintf(middle[rand.Intn(len(middle))], result.Titulo, strconv.Itoa(result.Publicado), result.Autor)

    //return fmt.Sprintf(initials[rand.Intn(len(initials))] + "\n\n" + recommendations[rand.Intn(len(recommendations))], user)
    return fmt.Sprintf(initials[rand.Intn(len(initials))] + "\n\n" + message , user), result.URLPortada
}

func getImage (imageURL string) []byte {
    resp, err := http.Get(imageURL)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    image, _ := ioutil.ReadAll(resp.Body)

    return image
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

    client, user, err := getClient(&creds)
    if err != nil {
        log.Println("Error getting Twitter Client")
        log.Println(err)
    }
    
    // Estructura con los parámetros fijos de acceso al servidor
    acc := DBAccess{
        user:     "admin",
        password: "admin",
        protocol: "http",
        host:     "sibila.website",
        //host:     "localhost",
        port:     "2480",
        database: "portico",
    }
    // Print out the pointer to our client
    // for now so it doesn't throw errors

    log.Println("Stream reading started")
    
    params := &twitter.StreamFilterParams{
	    Track: []string{"#quieroscifi"},
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

        receivedTweet := Twitt{
            Class:           "Twitt",
            ID:              strconv.FormatInt(tweet.ID, 10),
            Text:            tweet.Text,
            AuthorID:        strconv.FormatInt(tweet.User.ID, 10),
            AuthorName:      tweet.User.ScreenName,
            ConversationID:  "001", // Don't know which ID is this
            InReplyToUserID: strconv.FormatInt(tweet.InReplyToUserID, 10) }

        insertTwittDirect(acc, receivedTweet)
        //result, tweetStatusCode, status := insertTwitt(acc, t)
        //fmt.Println("Response Info (insertTwitt):")
        //fmt.Println(result)
        //fmt.Println(tweetStatusCode)
        //fmt.Println(status)


        // Tweet response text
	    
        answer, mediaURL := generateTweetAnswer(acc, tweet.User.ScreenName)
		log.Printf("Tweet Answer: %s\n", answer)
        log.Printf("Tweet Pic: %s\n", mediaURL)
	    tweetParams := &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID}
        
        image := getImage(mediaURL)
        imgRes, _, imgErr := client.Media.Upload(image, "IMAGE")
        if imgErr == nil {
            log.Printf("Media ID: %d", imgRes.MediaID)
            tweetParams.MediaIds = []int64{imgRes.MediaID}
        }

        // Responding tweet
        answerTweet, resp, err := client.Statuses.Update(answer, tweetParams)
	    if err != nil {
		    log.Printf("Statuses.Tweet error %v", err)
		}
		log.Printf("Tweet Status Code: %d\n", resp.StatusCode)

        insertTwittRelation(acc, strconv.FormatInt(answerTweet.ID, 10), strconv.FormatInt(tweet.ID, 10), "replied_to")
        //_, deployStatusCode, _ := insertTwittRelation(acc, strconv.FormatInt(answerTweet.ID, 10), strconv.FormatInt(tweet.ID, 10), "replied_to")
        //fmt.Println(deployStatusCode)
        
        // Retweeting the original tweet
        retweetParams := &twitter.StatusRetweetParams{TrimUser: twitter.Bool(true)}
        retweet, retweetResponse, err := client.Statuses.Retweet(tweet.ID, retweetParams)
        if err != nil {
            log.Printf("Statuses.Retweet error %v", err)
        }
        log.Printf("Retweet Status Code: %d\n\n", retweetResponse.StatusCode)
  
        insertTwittRelation(acc, strconv.FormatInt(retweet.ID, 10), strconv.FormatInt(tweet.ID, 10), "retweeted")
        //_, retweetStatusCode, _ := insertTwittRelation(acc, strconv.FormatInt(retweet.ID, 10), strconv.FormatInt(tweet.ID, 10), "retweeted")
        //fmt.Println("Response Info (insertTwittRelation):")
        //fmt.Println(retweetStatusCode)
        
	}

	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	
    log.Println("Stream reading stoped")
	stream.Stop()

}