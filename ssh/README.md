# 🌳 SSH

## Summary

edit your ssh config file

### Usage 

ssh -n example -h example.com -u root

## Description

The ssh command provides the ability to update your ssh config file
through the command line.

Options:

-n	 	:	name to identify this ssh section or hostname that should be used to establish the connection

-host	:	hostname that should be used to establish the connection 

-u 		:	username to be used for the connection

-i 		:	private key that the client should use for authentication when connecting to the ssh server

-p 		:	port that the remote SSH daemon is running on. only necessary if the remote SSH instance is not running on the default port 22
