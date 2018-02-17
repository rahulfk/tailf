# tailf

You need to have go installed in your system as a prerequisite.

## Installation

You can install the client and agent package using following command.

go get github.com/rahulfk/tailf/tagent

go get github.com/rahulfk/tailf/tclient

## Execution

For running agent, execute 

tagent -addr=localhost:8080

For running client, execute 

tclient -addr=localhost:8080 -path=\<absolute path of the file to tail\>
