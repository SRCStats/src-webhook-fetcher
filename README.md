# SRC Webhook Fetcher
A relatively basic program designed to search the API of [speedrun.com](https://speedrun.com) for new runs in the new, verified, and rejected categories, then subsequently store them in a MongoDB database. Designed for a serverless model running on an interval, working alongside another function to respond to webhooks for each run.

# Configuration
## **IMPORTANT**
This application is not intended to be used yet, as it currently misses features such as pagination. These instructions are primarily for collaborators developing the project.

Run the following commands in the desired location (other than $GOPATH):
```bat
git clone https://github.com/SRCStats/src-webhook-fetcher.git
cd src-webhook-fetcher
```
Configure a new MongoDB instance, or use an existing one. ([Azure Cosmos DB](https://azure.microsoft.com/en-us/services/cosmos-db/) w/ MongoDB API is used in production, but any will do.) Run the following commands to set the environment variables:

### Windows
```bat
set SRC_WEBHOOK_MONGODB_CONNECTION_STRING={Connection_String}
set SRC_WEBHOOK_DATABASE={Database_Name} &:: "srcstats" Recommended
set SRC_WEBHOOK_COLLECTION={Collection_Name} &:: "webhook-last-runs" Recommended
```
### Mac/Linux
```bash
# This is temporary for your current terminal session, add these commands to your .bashrc or .bash_profile to persist across sessions
export SRC_WEBHOOK_MONGODB_CONNECTION_STRING={Connection_String}
export SRC_WEBHOOK_DATABASE={Database_Name} # "srcstats" Recommended
export SRC_WEBHOOK_COLLECTION={Collection_Name} # "webhook-last-runs" Recommended
```
Finally, run the following go commands to finish the configuration and run the project:
```
go get ./...
go run cmd/main.go
```

# License
This project is licensed under the Anti-Capitalist Software License. See [LICENSE](LICENSE) for more information.