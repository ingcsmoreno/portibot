import requests
import os
import json
import dbmanager

# To set your environment variables in your terminal run the following line:
# export 'BEARER_TOKEN'='<your_bearer_token>'

BEARER_TOKEN = 'AAAAAAAAAAAAAAAAAAAAAFntLAEAAAAARIy7WWhUCNjrkPz6BVs4eCEXzM8%3D0AyN7C2vAzd2MsUpjaDSxhXGb4JmCXxR2om6x1MWeyDN3S179N'
def auth():
    return BEARER_TOKEN #os.environ.get("BEARER_TOKEN")

def get_userid_by_usarname(username):
    url = "https://api.twitter.com/2/users/by?usernames={}".format(username)
    headers = {"Authorization": "Bearer {}".format(BEARER_TOKEN)}
    response = requests.request("GET", url, headers=headers)
    if response.status_code != 200:
        raise Exception(
            "Request returned an error: {} {}".format(
                response.status_code, response.text
            )
        )
    json_resp = response.json()
    return json_resp['data'][0]['id']

def create_url(user_id):
    # Replace with user ID below
    # user_id = 2244994945
    return "https://api.twitter.com/2/users/{}/tweets".format(user_id)


def get_params():
    # Tweet fields are adjustable.
    # Options include:
    # attachments, author_id, context_annotations,
    # conversation_id, created_at, entities, geo, id,
    # in_reply_to_user_id, lang, non_public_metrics, organic_metrics,
    # possibly_sensitive, promoted_metrics, public_metrics, referenced_tweets,
    # source, text, and withheld
    return {"tweet.fields": "created_at,author_id,conversation_id,in_reply_to_user_id", "expansions":"referenced_tweets.id"}


def create_headers(bearer_token):
    headers = {"Authorization": "Bearer {}".format(bearer_token)}
    return headers


def connect_to_endpoint(url, headers, params):
    response = requests.request("GET", url, headers=headers, params=params)
    print(response.status_code)
    if response.status_code != 200:
        raise Exception(
            "Request returned an error: {} {}".format(
                response.status_code, response.text
            )
        )
    return response.json()


def main():
    user_id = get_userid_by_usarname('martincasatti')
    bearer_token = auth()
    url = create_url(user_id=user_id)
    headers = create_headers(bearer_token)
    params = get_params()
    json_response = connect_to_endpoint(url, headers, params)
    #print(json.dumps(json_response, indent=4, sort_keys=True))

    
    db = dbmanager.DBManager()
    '''
    print(len(json_response['data']))
    for twitt in json_response['data']:
        print ("ID: {} -> texto: {} -> referenced: {} ({})".format(twitt['id'],twitt['text'],twitt['referenced_tweets'][0]['id'],twitt['referenced_tweets'][0]['type']))
        try:
            db.insertTwitt(
                twitt['id'],
                twitt['text'],
                twitt['author_id'],
                twitt['conversation_id'],
                twitt.get('in_reply_to_user_id',"")
            )
            print ("Insertado el Twitt:",twitt['id'])
        except Exception as e:
            print ("Se produjo un error:",e)    
    '''
    db.insertTwittRelation('1359949376844660738','1359951287383654400','replied_to')


if __name__ == "__main__":
    main()