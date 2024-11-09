echo $1 | osascript 3<&0 <<'APPLESCRIPT'
  on run argv
    set stdin to do shell script "cat 0<&3"

    set retList to {}
    using terms from application "Mail"
      tell application "Mail"
        repeat with theMessage in (every message of mailbox "INBOX" of account stdin)
          set dateV to (get date sent of theMessage)
          set subjectV to (get subject of theMessage)
          set contentV to (get content of the theMessage)

          set fieldDelim to "FIELD_DELIM"
          set entry to "Date: " & dateV & fieldDelim & "Subject: " & subjectV & fieldDelim & "Content: " & contentV & fieldDelim
          set entryDelim to "ENTRY_DELIM"
          set entry to entry & entryDelim
          copy entry to the end of the retList
        end repeat
      end tell
    end using terms from

  return retList

  end run
APPLESCRIPT