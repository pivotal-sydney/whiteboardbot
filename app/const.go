package app

const (
	THUMBS_UP   = ":+1:"
	THUMBS_DOWN = ":-1:"
	USAGE       = "*Usage*:\n" +
		"        `/wb [command] [text...]`\n" +
		"    where commands include:\n" +
		"*Registration Command*\n" +
		"        `register`, `r` - followed by <standup_id>, registers current channel to Whiteboard's standup id\n" +
		"\n" +
		"*Presentation Command*\n" +
		"		 `present`, `p` - presents today's standup. Follow with number of days to limit the entries shown by date (i.e. `/wb p 2` will only return entries for the next 2 days)\n" +
		"\n" +
		"*Create Commands*\n" +
		"        `faces`, `f` - followed by a title, creates a new faces entry\n" +
		"        `interestings`, `i` - followed by a title, creates a new interestings entry\n" +
		"        `helps`, `h` - followed by a title, creates a new helps entry\n" +
		"        `events`, `e` - followed by a title, creates a new events entry\n" +
		"\n" +
		"*Detail Commands* (updates details of a started entry)\n" +
		"        `title`, `t`, `name`, `n` - updates a name/title detail to a started entry\n" +
		"        `body`, `b` - updates a body detail to a started entry\n" +
		"        `date`, `d` - updates a date detail to a started entry (YYYY-MM-DD)\n" +
		"\n" +
		"Example:\n" +
		"        `/wb f New Face!` - will create a new face with the name 'New Face!'\n" +
		"        `/wb d 2015-01-02` - will update the new face date to 02 Jan 2015"

	NEW_ENTRY_HELP_TEXT = "_Now go update the details. Need help?_ `wb ?`"
	MISSING_STANDUP     = "You haven't registered your standup yet. `/wb r <id>` first!"
	MISSING_ENTRY       = "Hey, you forgot to start new entry. Start with one of `/wb <command> [title]` first!\nNeed help? Try `/wb ?`"
	MISSING_INPUT       = "Hey, next time add a title along with your entry!\nLike this: `/wb <command> My title`\nNeed help? Try `/wb ?`"
)
