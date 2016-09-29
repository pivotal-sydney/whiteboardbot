# Whiteboard Slack Bot

WhiteboardBot is a slack bot that allows users to easily create new entries into whiteboard directly from Slack!

To improve the user experience of creating entries during your commute to work, by talking to the bot with a set of keywords, it will allow you to create and update new faces, events, interestings and helps.

# Usage
## Register A Standup
Whiteboard Bot needs to know which standup to post to. It keeps a mapping between Slack channel and standup. To register a standup for your current channel, execute:
```
/wb register <id>

# example
/wb r 1
```
where <id> refers to the integer ID of your standup provided by Whiteboard. You're now ready to create entries to your standup!

## <a name="create">Creating a new Whiteboard Entry
To create a whiteboard entry
```
/wb <command> <title/name>
```
where 
* `<command>` can be one of [faces, events, interestings, helps]
* `<title/name>` is a title/name of an entry

Examples:
```
/wb faces My new face!
/wb interestings Some intersting title
/wb events The new event title happening soon!
/wb helps How does the whiteboard bot work!?
```

If the command is accepted, the bot will upload your entry to the Whiteboard and respond to you with the content of your entry. It also remembers your recent entry to allow you to update your entry and set other details (see below).

##  Setting details on Whiteboard Entry
* To update a name/title (they do the same thing)
```
/wb name My Name
```
```
/wb title My Title
```

* To set/update a body
```
/wb body My Body
```

* To update the date [YYYY-MM-DD] (defaults to today)
```
/wb date 2015-12-01   // December 1st, 2015
```

You can continue to edit the entry until you begin [creating a new entry](#create).

## Having The Bot Speak In Your Channel

Whiteboard Bot also supports a mode where it speaks to everyone in the channel. After inviting the bot into your channel, use the exact same commands you're familiar with, but without the leading /. Whiteboard Bot will then reply to commands aloud in the channel.

# Additional Features
## Abbreviations
Whiteboardbot recognizes abbreviations of each command.  It can recognize the best match to each command.  
For example:
```
/wb f New Face
```
will be recognized as creating a new faces entry.
```
/wb i Something interesting
```
will be recognized as creating an interestings entry.
```
/wb b Some body
```
will be recognized as setting the body (description) of the current entry.

## Presentation
You can now use the bot in presentation mode! The bot will show you the list of all the items for today so you can run standup directly in Slack!
```
/wb present
/wb p
```

## Help/Usage
You can ask bot for help by typing
```
/wb ?
```
It will respond with the information on how to use it, containing all supported commands and some examples.

## Command Case Insensitivity
Most users will probably be adding entries on their phones. Most mobile phones will capitalize the first letter you type.
Luckily, Whiteboardbot commands are case insensitive! So even if your phone starts capitalizing a command, it will still work!
```
/Wb InTeReStInG This will still work!
```

# Setting up and running bot
In order to have the bot work correctly, you need to have several ENV variables configured.

```
SLACK_TOKEN=payloadtoken              // The token Slack sends with slash command payloads
WB_HOST_URL=http://localhost:3000     // The host url of the Whiteboard App
WB_AUTH_TOKEN=rails-auth-token        // Visit Whiteboard in your browser and copy the contents of `csrf-token` meta tag
WB_BOT_API_TOKEN=someapitoken         // The API token of your bot.  See Slack docs to create a bot, and get API token
WB_DB_HOST=localhost:6379             // The Redis IP address with port
WB_DB_PASSWORD=password               // The Redis password
```

Set these in the `manifest.yml` file you copied from `manifest.yml.example` if you're deploying to Cloud Foundry.

## Building
* Set GOPATH env variable
* Check out whiteboardbot project from github using go get: `go get github.com/pivotal-sydney/whiteboardbot`
* Go to the project directory: `cd $GOPATH/src/github.com/pivotal-sydney/whiteboardbot/`
* Install godep (for dependency managment): `go get github.com/tools/godep`
* Fetch all dependencies using godep: `godep restore`
* Update all dependencies using godep: `godep update ./...`
* Now you're ready to build the project: `go build` This will create a whiteboardbot binary which can be run from the command line.
* To run the test execute this command: `go test ./...`

## Deploying To Cloud Foundry
* Set GOPATH env variable
* Check out whiteboardbot project from github using go get: `go get github.com/pivotal-sydney/whiteboardbot`
* Go to the project directory: `cd $GOPATH/src/github.com/pivotal-sydney/whiteboardbot/`
* Install godep (for dependency managment): `go get github.com/tools/godep`
* Fetch all dependencies using godep: `godep restore`
* Update all dependencies using godep: `godep update ./...` (this creates the `vendor` dir)
* Copy the sample manifest and fill it in with your details: `cp manifest.yml.sample manifest.yml`
* Push the app using the Cloud Foundry CLI: `cf push`
