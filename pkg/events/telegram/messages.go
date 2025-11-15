package telegram

const msgHello = `🪩 Welcome to NAMNADA LINK - your personal reading list assistant!

Just send me a link - I'll save it.

When you're ready to read:
🕹️ Use /random to get a surprise link  
🕹️ Use /list to view everything saved

You can also manage your links 🔗

🔎To view detailed information use /help

Happy reading! 📬`

const msgHelp = `🧿 NAMNADA LINK can help you save and manage links to read later.
Here’s what you can do:

/random - Get a random unread article  
/read - Mark an article as read  
/remove - Delete an article  
/list - Show all saved articles  
/help - Show this help message

Just send me any link, and I’ll save it automatically! 💾`

const (
	msgSaved          = "💾 Saved to your reading list!"
	msgNoSavedPages   = "🕰️ You have no saved pages yet.\nJust send me a link to get started!"
	msgAlreadyExists  = "📰 This page is already in your list"
	msgMarkedAsRead   = "🧮 Marked as read!"
	msgRemoved        = "🗑️ Page removed!"
	msgUnknownCommand = "🥡 I didn't understand that command.\nTry /help to see what I can do!"
	msgURLRequired    = "🔗 Please provide a valid URL"
)
