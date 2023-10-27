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

To run using gcottom's glfw-js:
```
cd $GOPATH/src/github.com/
mkdir gcottom
cd gcottom
git clone https://github.com/gcottom/glfw-js.git
cd glfw-js
git branch -a # you should see 'clipboard-fix'
git merge origin/clipboard-fix master
```
Then clone this repository
```
cd $GOPATH/src/github.com/
mkdir rmasci
cd rmasci
git clone https://github.com/rmasci/jsontool.git
cd jsontool
```

Edit go.mod (Make sure go version on line 3 is correct, and then make sure you've got the right path to glfw-js

```
replace github.com/fyne-io/glfw-js => /Users/youruser/go/src/github.com/gcottom/glfw-js
```

You should be able to run this:
```
fyne serve --port 8080 --target wasm
```
