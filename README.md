# Treasure Coin

Treasure coin is a Proof of Concept for the OST Kit alpha III challenge. We built a treasure hunting proof of concept platform that uses printed QR codes as treasures and cryptocurrency as rewards.

A better explanation can be found in the [blog post](https://medium.com/@pdomingos.costa/tresurecoin-ost-kit-alpha-phase-iii-f66705ddfd84) describing the project.

## Features

- Integration with the OST Kit API
- Authentication and Authorization into the platform
- Managing user wallets through the OST kit API
- Allows the creation of games with an arbitrary number of treasures
- Creating games costs branded tokens
- Fingding a treasure rewards branded tokens
- Historical data of the results of playing events
- Transaction history

## Requirements

The application is fully implemented in Golang and uses BoltDB for persistence. To compile the project simply get the required dependencies, and compile `cmd/main.go`. The Makefile also contains the most common operations.

In order to connect to the OST APIs, a `.env` file is required in the root of the repository. This file can be created by using the `.env.sample` file as a template.

## Issues

All issues found and discussion about the technical aspects of the project, can be done through the Issues section of the Github Repository.

## Contributing

During the challenge we will not accept external contributions (since this repository is part of the on-going challenge).

After this period and the results published, anyone is welcome to send their contributions and patches through as Pull Requests.
