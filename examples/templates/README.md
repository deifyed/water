
## Templates directory

The templates directory contains a collection of directories, one for each scaffoldable item (be it file or folder). Each
scaffoldable item MUST contain a [metadata.json](../metadata.json) file which is used to contextualize multiple variants of
an item.

## Item metadata

The [metadata.json](../metadata.json) file contains a list of metadata entries.

Each metadata entry has a `target` which points to a file or a directory inside the scaffoldable item, and a list of tags.

## Definitions

| Word              | Definition                                                                       |
|-------------------|----------------------------------------------------------------------------------|
| Scaffoldable item | A directory containing definitions for a file or a folder that can be scaffolded |
