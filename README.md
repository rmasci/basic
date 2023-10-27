# JSON Edit

This tool is just a sample tool. What is interesting is that it fixes json on the fly.
Once run in the right side if you type invalid json, it fixes it in the right.
```
{"First" "Sherlock" "Last" "Holmes" "Detective" true "Friend" "Watson"}
```
This turns it into:
```
{
  "First": "Sherlock",
  "Last": "Holmes",
  "Detective": true,
  "Friend": "Watson"
}
```

Experimental is the package used by gcottom so that Fyne's wasm supports cut and paste.  On a mac cmd-v doesn't work, but you should be able to right mouse click and select copy / pase. The menu items for 'Copy' in the app work.  It's a sample, work in progress. The file open / save doesn't work. If you want to give me a pull request for that please do.
