# Note

This project originally had over 100 commits. This is a public version of one of my private repos. The history of the private repo has been squashed into a single commit. That is the reason why there are so few commits.

# Introduction

This is a tool that can be used to simulate digital circuits on the web.

# Required Softwares

The softwares that are required to run this software are as follows:

- [MySQL v8.0](https://dev.mysql.com/downloads/installer/)
- [Node.js v20](https://nodejs.org/en/download/prebuilt-installer/current)
- [Golang v1.22](https://go.dev/dl/#go1.22.6)

# Before Running

Make sure that the mysql service is running on your device. On windows you can run the mysql service as follows:

- Open a powershell window as administrator
- Run the following command:

  `net start mysql80`

# How to Run

There are two major steps to running this software and they are as follows:

## Server Side

### Configuration

#### Database

You have to configure your database before running the server. In order to do so, run the following command in the terminal:

`go run . --configdb`

#### DotEnv

Create a file named `.env` in the project root directory and add the following line to it:
`TOKEN_SECRET=<a-random-string-of-characters>`

### Running the Server

After the database has been configured, you can launch the server with the following command:

`go run .`

## Client Side

### Installing Dependencies

Run the following command to install dependencies:

`npm i`

### Configuration

- Run `ipconfig` on Windows or `ifconfig` for Linux and MacOS and look for your machine's network ip address.
- Make a file in the client directory and name it `.env`.
- In it paste the network ip address in the following format:
  `VITE_API_SERVER="http://<network-address>:8080"`

### Running the Client

In order to run the client run the following command:

`npm run dev-host`

And go to the link that is printed in the console
