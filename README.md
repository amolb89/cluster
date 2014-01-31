GoTalk

This library is written in golang. It can be used to send point to point or broadcast messages among a cluster of servers. Configuration about cluster is provided by .json config file. Configuration includes unique id of the server and its address.

Usage

To test the point to point and broadcast functionality use,

go test github.com/amolb89/cluster/server

Test module tests send and receive of large number of point to point and broadcast messages. Test module uses cluster information from serverConfig1.json, it has server id and address on which the server would run.

How it works
Each server has outbox and inbox channels on which it sends messages and receives messages from. Outbox channel can be accessed with Outbox() go routine, similarly inbox channel can be accessed via Inbox() go routine.

