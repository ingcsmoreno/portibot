module github.com/ingcsmoreno/tweetbot-r

go 1.15

// This replace is a temporal workaround to get media upload features
// Should be removed once this PR gets merged: https://github.com/dghubble/go-twitter/pull/148
replace github.com/dghubble/go-twitter => github.com/janisz/go-twitter v0.0.0-20201206102041-3fe237ed29f3

require (
	github.com/dghubble/go-twitter v0.0.0-20201011215211-4b180d0cc78d
	github.com/dghubble/oauth1 v0.6.0
	github.com/dghubble/sling v1.3.0
	github.com/go-resty/resty/v2 v2.5.0
	github.com/tidwall/gjson v1.6.8
)
