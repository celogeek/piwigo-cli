# piwigo-cli

This tools allow you to interact with your Piwigo Instance.

## Install

```
go install github.com/celogeek/piwigo-cli/cmd/piwigo-cli@latest
```

## Quickstart

### Login

```
piwigo-cli session login -u URL -l USER -p PASSWORD
```

### Check your status

```
piwigo-cli session status
```

### List Categories

```
piwigo-cli categories list
```

### List images in a category and sub categories

```
piwigo-cli images list -c 4 -r
```

With a tree style

```
piwigo-cli images list -c 4 -r -t
```

### Upload a tree of images

This will also create all the categories using the directory name.

In this example, the category 4 (-c 4) is 2021
```
images upload-tree -d ~/Downloads/2021 -j4 -c 4
```

This will create the categories based on your local directories, and upload only the images that doen't already exists somewhere else.

The check of the images existance use MD5 checksum. If you change the metadata of the images, it will reupload the image as a new one.

You can remove the duplicates in piwigo by looking for similar photo with Name & Date & Size. Just pickup the first one so when you try to upload again, it won't reupload the images.

## License

MIT