
# block-listener

Hyperledger Fabric Block Listener and Parser on Huawei BCS

## Setup

Open setup/setup.go to check the config variables (Default should work):
- configPath is the channel sdk config location 
- orgId and userName follows the one provided by BCS

Run retup.go. This should do the following:
- producing vars.json, used by the /config functions to setup the configuration variables
- replacing the directory in channel sdk config to your current working directory.
  change the oldDir variable accordingly everytime you move/rename the repo
  
## Listener

Open listener/listener.go to check the config variables (Default should work):
- seekType dictates the listening behavior. There are three types:
  - seek.Oldest goes from first block
  - seek.Newest goes from the last block
  - seek.FromBlock goes from a given block, dictated by startBlock variable
- startBlock dictates the starting block when seekType is seek.FromBlock

With BCS service on, running listener.go will:
- listen to blockEvents according to the given config variables
- parse the blockEvents into json files in the /block-event-parses folder:
  - the prefix in- is given to indentMarhsal-ed json
These occur continuously until the process is interrupted
