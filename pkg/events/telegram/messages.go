package telegram

const msgHelp = `I can save links to pages for you and help you manage them.

Here’s what I can do:

/start — show welcome message
/random — send you random unread pages
/read — mark page as read
/remove — delete page
/list — show all saved pages
/help — show this help message

You can also just send me a link, and I'll save it for you automatically.`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgSaved          = "Saved!"
	msgNoSavedPages   = "You have no saved pages"
	msgAlreadyExists  = "You already have this page in your list"
	msgMarkedAsRead   = "Marked as read!"
	msgRemoved        = "Page removed!"
	msgUnknownCommand = "Unknown command"
	msgURLRequired    = "Please provide a URL"
)
