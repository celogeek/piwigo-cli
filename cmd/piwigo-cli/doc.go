/*
Piwigo Cli

This tools allow you to interact with your Piwigo Instance.
Installation:
 go install github.com/celogeek/piwigo-cli/cmd/piwigo-cli@latest

Help

To get help:
 piwigo-cli -h

QuickStart

Login

First connect to your instance:
 piwigo-cli session login -u URL -l USER -p PASSWORD

Then check your status
 piwigo-cli session status
*/
package main
