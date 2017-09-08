# A not too fancy telnet server in golang.  All substance, no bells and whistles.  Standard stuff.

## Steps to setup on systems using unix / which have make
- Ensure you have an appropriate config.json modeled after the one in the repo.
  - Note: If you change the logFile to something else and want the make command to clean it, change it in the Makefile as well.
- Type make in the home directory.  It should start on its own.
- Connect to the server through telnet (on Unix, the command should be telnet IPHostAddr Port)

## Limitations

Although currently the chat server can block messages from users who are in a users blockedUsers map, it does not actually have the ability to block a user.

I thought a possible way to implement this is to check the input sent by each user and add a conditional which determines whether or not it is a message or a command.

Reminiscent of early Battle.net days from my childhood, I want to invoke commands by prefacing the message with a /.

## Approach

I began by scouring the web to look for writing basic chat servers in golang.  I wanted to see multiple working versions of a chat server before I actually started to make implementation decisions on my own,
and I wanted the luxury of comparing and seeing why some people chose one method of implementation compared to the other.

After finding a couple I felt were solid, I read through their source and looked up any methods they used in the standard go library I didn't understand.  


## Learning
Having written the prototypes to a couple of CLI tools before in Golang, I had an idea of how to use the language to interact with the file system.

Channels were new to me however.  Although it is working, I believe I could edit the current number of channels and ways I write to the channels to do better in the future. 

I also think it would help the chatroom in general to have Room structs which have their own individual channels.  I think in the future I could subscribe my individual client's channels to listen to a room's channels and implement some type of room switching that way.
