# flowbang
Support a flow with github projects and issues (manage issues, branches and pull requests)

This repository is a playground for me to understand and learn Go.
If you have any issues, concerns or simply think something should be done better,
don't hasitate to write to me or create an issue.

## How to run

```bash
go get github.com/wrutek/flowbang
go build github.com/wrutek/flowbang
```

No you have only one command that is working `configure`. 
This command creates ~/.config/flowbang/flowbang.conf file and stores there
your answares to 3 simple questions:
 - github oauth token
 - repository on which you want to workon
 - repository in which a project board is configured (sometimes it is different from above repo)

 ## KUDOS

 I would like to say some warm words to all contributors of [termbox-go](https://github.com/nsf/termbox-go).
 You are awesome guys!
 