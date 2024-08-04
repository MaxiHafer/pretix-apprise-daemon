# pretix-apprise-daemon
This daemon aims to provide notification capabilities for pretix. It consumes pretix webhooks and sends out apprise notifications.

## Currently WIP
The project is currently a work in progress, expect no working builds or complete features. Feel free to open your PR I am more than happy to integrate good changes

## Future features

In the future, it should be possible to customize the notification body using go templates, currently this is only a POC and the templates are compiled in using embed.FS, so any changes to the templates mandate a recompilation of the binary. The cleanest solution would probably be to read in the templates at runtime ¯\\_(ツ)_/¯
