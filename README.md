# Goa Plugins

This repository contains plugins for 
[Goa v3](https://github.com/goadesign/goa/tree/v3).

Each sub-directory README contains the information and usage for the given
plugin.

Plugins v3.31.0 accompanies Goa v3.31.0 and requires Go 1.26 or later.
Upgrade both modules together:

```shell
go get goa.design/goa/v3@v3.31.0 goa.design/plugins/v3@v3.31.0
go install goa.design/goa/v3/cmd/goa@v3.31.0
```

Regenerate the complete application after upgrading. This release adapts
declaration-producing plugins to Goa's generation plan and updates generated
commands, transport helpers, and test scenarios. Custom plugins that declare
names must plan those declarations before rendering. Read the
[Goa upgrade guide](https://github.com/goadesign/goa/blob/v3.31.0/UPGRADING.md)
for application changes and the generator API migration.

Goa v2 plugins can be found [here](https://github.com/goadesign/plugins/tree/v2).
