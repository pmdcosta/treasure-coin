# Treasure Coin

# Setup
Create a `.env` file in the root of the project, using the `.env.sample` as a base. Simply provide the
missing data.

# Start
Run `make start` to start the server.

# TODO - JC
If you have any idea to make it a bit prettier go ahead, otherwise the only thing missing is the script
to make the 1k requests to the API.
I suggest making it in Go, to get a bit more exp, and you can also use the package we already built.
Make this script as a command line file that gets its parameters as flags (check the main.go for clues).

The script should create an admin user, airdrop him some tokens, and he should create a bunch of games.
The script should also create a number of users, and have these users find the treasures from the created games.
We should also create a new action for dumping tokens, so that we can easily empty out the user accounts.

## Bonus points
You should be able to supply as a parameter the number of games to create, the number of treasure
per game and the number of users to create.
Try to make user of go routines and channels to find treasures in parallel.
