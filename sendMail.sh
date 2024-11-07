echo $1 | osascript 3<&0 <<'APPLESCRIPT'
to splitString(someString)
    	try
    		set tempTID to AppleScript's text item delimiters -- save current delimiters
    		set AppleScript's text item delimiters to "|"
    		set pieces to text items of someString -- split the string
    		set AppleScript's text item delimiters to tempTID -- restore old delimiters
    		set firstPart to item 1 of pieces
    		set secondPart to item 2 of pieces
    	on error errmess -- delimiter not found
    		log errmess
    		return {firstPart, ""} -- empty string for missing item
    	end try
    	return {firstPart, secondPart}
    end splitString

  on run argv
    set stdin to do shell script "cat 0<&3"
    set {toAddy, mySubject} to splitString(stdin)

    using terms from application "Mail"
    	tell application "Mail"
    		--create new message with subject and content
    		set newMessage to make new outgoing message with properties {subject:mySubject, content:""}
    		-- add the To: addresses to the new message
    		tell newMessage
    			make new to recipient with properties {address:toAddy}
    			send
    		end tell
    	end tell
    end using terms from
  end run
APPLESCRIPT