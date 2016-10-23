A bunch of Lambda functions that collect data that is then sent to IFTTT where it can be handled however you like. The underlying idea is similar to [Slogger][slogger], but limited to tasks that can run on an external server. Instead of running them on a server however, they will run as scheduled Lambda functions thereby removing the need to maintain a server.

# Installation

The easiest way is to pick the collectors you want from the releases section, adjust the configuration file in them and then upload this package to a Lambda function. You should then set the schedule you prefer using CloudWatch schedules. One way to do this, but not required, is using my [Aqua][aqua] application.

On the IFTTT side, you will need to enable the Maker channel and create a trigger with it. For example, you can create a trigger that reads the input and then sends it to Day One.

[aqua]: https://github.com/ArjenSchwarz/aqua

# Configuration

The configuration for each collector is in its own `config.yml` file, examples (my actual config) can be found in the respective directories. You will need to adjust these before you upload to Lambda. The KMS flag will determine whether the collector will try to

# Functionalities

## GitHub

The GitHub collector

# Development

Each collector is designed to run separately as its own Lambda function. The collectors are written in Go and can easily be compiled for use with Lambda using `bin/build.sh` which will create a Lambda ready zip file with the same name containing the (unfortunately) required index.js, the config file, and the binary. For example, the below command will create a zipfile called `github.zip` in the lambdapackage directory.

```bash
$ bin/build.sh github
```

##  

[slogger]: https://github.com/ttscoff/slogger
