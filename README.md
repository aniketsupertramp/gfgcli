# gfgcli
a command-line tool that use https://www.geeksforgeeks.org interviews-corner section to fetch companies and their corresponding interview-experiences articles. 

Using this, you can pretend to be working in the office while simaltaneously, preparing for the next job. :p

This tool uses github.com/gocolly/colly for fecthing html content, github.com/manifoldco/promptui for interactive terminal and github.com/gernest/wow for the spinner. (I have modified github.com/manifoldco/promptui library to add one more key as "End" to terminate the session.)

Prerequisite: <br />
1.)Go

Setup: <br />
1.)***cd $GOPATH*** <br />
2.) RUN ***go get github.com/aniketsupertramp/gfgcli/src*** <br />
3.) ***go build -o gfgc src/github.com/aniketsupertramp/gfgcli/src/*.go*** <br />
4.) ***cp gfgc /usr/local/bin*** <br />

***Done !!!***

![](gfgExample.gif)
