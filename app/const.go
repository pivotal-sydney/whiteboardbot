package app

const(
	THUMBS_UP = ":+1:\n"
	THUMBS_DOWN = ":-1:\n"
	USAGE =
	"*Usage*:\n" +
	"        `wb [command] [text...]`\n" +
	"    where commands include:\n" +
	"*Registration Command*\n" +
	"        `register`, `r` - followed by <standup_id>, registers current channel to Whiteboard's standup id\n" +
	"\n" +
	"*Presentation Command*\n" +
	"		 `present`, `p` - presents today's standup. Follow with number of days to limit the entries shown by date (i.e. `wb p 2` will only return entries for the next 2 days)\n" +
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
	"        `wb f New Face!` - will create a new face with the name 'New Face!'\n" +
	"        `wb d 2015-01-02` - will update the new face date to 02 Jan 2015"
)
