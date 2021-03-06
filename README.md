A bunch of Lambda functions that collect data that is then sent to IFTTT where it can be handled however you like. The underlying idea is similar to [Slogger][slogger], but limited to tasks that can run on an external server. Instead of running them on a server however, they will run as scheduled Lambda functions thereby removing the need to maintain a server.

Output sent to IFTTT is limited to a Markdown formatted body. If you want to use it with Day One, you can set tags and journal in your IFTTT trigger.

[slogger]: https://github.com/ttscoff/slogger

# Installation

The easiest way is to pick the collectors you want from the releases section, adjust the configuration file in them and then upload this package to a Lambda function. You should then set the schedule you prefer using CloudWatch schedules. One way to do this, but not required, is using my [Aqua][aqua] application.

On the IFTTT side, you will need to enable the Maker channel and create a trigger with it. For example, you can create a trigger that reads the input and then sends it to Day One. An example IFTTT trigger can be found [here][triggerexample].

[aqua]: https://github.com/ArjenSchwarz/aqua

[triggerexample]: https://ifttt.com/recipes/478666-maker-to-day-one

# Configuration

The configuration for each collector is in its own `config.yml` file, examples (my actual config) can be found in the respective directories. You will need to adjust these before you upload to Lambda. The KMS flag will determine whether the collector will try to decrypt the keys from the configuration.

# Functionalities

## GitHub

The GitHub collector collects your GitHub commits for the period you've set up (for example, the last 24 hours).

## Instapaper

The Instapaper collector collects the articles you (or someone else) saved in the configured feeds.

## Pinboard

The Pinboard collector collects the links you (or someone else) saved in the configured feeds.

# Development

Each collector is designed to run separately as its own Lambda function. The collectors are written in Go and can easily be compiled using the provided Makefile. They will all be compiled for Lambda, and can be deployed with the remaining Makefile functionalities.

## Contribution

Feel free to contribute. Either by improving existing collectors, providing new ones, or creating issues.

For Pull Requests, just follow the standard pattern.

1. Fork the repository
2. Make your changes
3. Make a pull request that explains what it does
