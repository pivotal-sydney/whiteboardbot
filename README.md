# Whiteboard Slack Bot

WhiteboardBot is a slack bot that allows users to easily create new entries into whiteboard directly from Slack!

To improve the user experience of creating entries during your commute to work, by talking to the bot with a set of keywords,
it will allow you to create and update new faces, events, interestings and helps.

# Usage
## <a name="create">Creating a new Whiteboard Entry
To begin creating a whiteboard entry
```
wb <command>
```
where `<command>` can be one of [faces, events, interestings, helps]

The bot remembers the context of what type of entry you are working with. If the command is accepted, it will respond with
the current state of the entry you are working with.

##  <a name="detail">Setting details on Whiteboard Entry
* To set/update a name/title (they do the same thing)
```
wb name My Name
```
```
wb title My Title
```

* To set/update a body
```
wb body My Body
```

* To set/update the date [YYYY-MM-DD] `(defaults to today)`
```
wb date 2015-12-01   // December 1st, 2015
```

Once all mandatory fields (fields denoted with a *) have been set on an entry, the entry will be saved and uploaded to Whiteboard.
You can continue to edit the entry until you begin [creating a new entry](#create)

## Setting up and running bot
In order to have the bot work correctly, you need to have several ENV variables configured.

```
WB_HOST_URL=http://localhost:3000     // The host url of the Whiteboard App
WB_BOT_API_TOKEN=someapitoken         // The API token of your bot.  See Slack docs to create a bot, and get API token
WB_DB_HOST=localhost:6379             // The Redis IP address with port 
WB_DB_PASSWORD=password               // The Redis password 
```

```
go build
```
will create a whiteboardbot binary which can be run from the command line.

## Add bot to Slack channel
Once the bot is running, mention `@whiteboardbot` in a channel, and register a standup ID
```
wb register <id>
wb r 1
```
where <id> refers to the integer ID of your standup provided by Whiteboard.  You're now ready to create entries to your standup!

# Additional Features
## Abbreviations
Whiteboardbot recognizes abbreviations of each command.  It can recognize the best match to each command.  
For example:
```
wb f
```
will be recognized as creating a new faces entry.
```
wb i
```
will be recognized as creating an interestings entry.
```
wb b Some body
```
will be recognized as setting the body (description) of the current entry.

## Presentation
You can now use the bot in presentation mode!  The bot will show you the list of all the items for today so you can run standup directly in Slack!
```
wb present
wb p
```

## Creating and setting title together
Whiteboard recognizes new entry creation commands with titles!  In order to reduce the amount of messages, you can actually set the title/name
of the new entry by including the title along with the [creating a new entry](#create) command.

The title/name is the only mandatory field, so this one liner is the minimum required in order to create a new entry.
```
wb f My new face!
wb i Some intersting title
wb e The new event title happening soon!
wb h How does the whiteboard bot work!?
```

## Uploading images via Mobile
You can now upload an image to Whiteboard (via any Slack client like iOS, Android, or Desktop).  Simply share an image, use a [create command](#create) via the title section of the upload, and add any text into the body via the comment section.  Making changes remains the same as before by [setting details](#detail) on the entry.
```
[Share an image via Slack]
[Title]
wb i Here's the title of the image I'm uploading
[Comment]
More context about the image I'm uploading
```

## Command Case Insensitivity
Most users will probably be adding entries on their phones.  Most mobile phones will capitalize the first letter you type.
Luckily, Whiteboardbot commands are case insensitive!  So even if your phone starts capitalizing a command, it will still work!
```
Wb InTeReStIng This will still work!
```
