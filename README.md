nfssh
=====

Simple utility to manage my ssh-tunneled NFS mounts (in OSX).

## Brief
1. Create a config.json with your details.
2. Ensure that you have added your private key to your ssh-agent (e.g. run `# ssh-add`)
3. `# ./nfssh` will then create the tunnels and mount the NFS share. When you give it the kill/interrupt signal it gracefully un-mounts the share.

## Details
It opens up a ssh-client (connection), utilizing the ssh-agent for authentication, and starts listening for two connections on localhost of the local machine (one port for nfsd, and one for mountd). Next we tell the mount utility to mount our NFS from localhost specifying that it should use our custom ports. When the mount utility opens new connections to our tunnel listeners they are piped to sockets dialed remotely through the ssh-client, reaching nfsd and mountd respectively. Following this the program waits for either the ssh-socket failing (in which case we shutdown the local listeners before attempting to reconnect the ssh-client every five seconds until we get back up, gracefully only remounting if needed), or a kill/interrupt signal (which means we un-mount the NFS share and tear down all the connections).

## LICENSE
Copyright (c) 2016, Kristoffer A. Iversen
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

The views and conclusions contained in the software and documentation are those
of the authors and should not be interpreted as representing official policies,
either expressed or implied, of the FreeBSD Project.
