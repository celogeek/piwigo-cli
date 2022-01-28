# Piwigo Cli

## Installation

This tools allow you to interact with your Piwigo Instance.
Installation:
```
$ go install github.com/celogeek/piwigo-cli/cmd/piwigo-cli@latest
```

## Help

To get help:
```
$ piwigo-cli -h

Usage:
  piwigo-cli [OPTIONS] <command>

Help Options:
  -h, --help  Show this help message

Available commands:
  categories  Categories management
  getinfos    Get general information
  images      Images management
  method      Reflexion management
  session     Session management
```

## QuickStart

### Login

First connect to your instance:
```
$ piwigo-cli session login -u URL -l USER -p PASSWORD
```

### Status

Then check your status
```
$ piwigo-cli session status
```

## Commands

### General commands

```
$ piwigo-cli getinfos
```

### GetInfos Command

General information of your instance.

```
$ piwigo-cli getinfos

┌───────────────────┬─────────────────────┐
│ KEY               │ VALUE               │
├───────────────────┼─────────────────────┤
│ version           │ 12.2.0              │
│ nb_elements       │ 39664               │
│ nb_categories     │ 816                 │
│ nb_virtual        │ 816                 │
│ nb_physical       │ 0                   │
│ nb_image_category │ 39714               │
│ nb_tags           │ 73                  │
│ nb_image_tag      │ 24024               │
│ nb_users          │ 3                   │
│ nb_groups         │ 1                   │
│ nb_comments       │ 0                   │
│ first_date        │ 2021-08-27 20:15:15 │
│ cache_size        │ 4242                │
└───────────────────┴─────────────────────┘
```

### Categories List Command

List the categories.

```
$ piwigo-cli categories list

┌──────┬───────────────────────────┬────────┬──────────────┬──────────────────────────────────────────────────────┐
│   ID │ NAME                      │ IMAGES │ TOTAL IMAGES │ URL                                                  │
├──────┼───────────────────────────┼────────┼──────────────┼──────────────────────────────────────────────────────┤
│    4 │ 2021                      │      0 │         1520 │ https://yourphotos/index.php?/category/4             │
│  677 │ 2021/Animals              │      0 │           49 │ https://yourphotos/index.php?/category/677           │
│   24 │ 2021/Animals/Cats         │     29 │           32 │ https://yourphotos/index.php?/category/24            │
│  760 │ 2021/Animals/Cats/Videos  │      3 │            3 │ https://yourphotos/index.php?/category/760           │
└──────┴───────────────────────────┴────────┴──────────────┴──────────────────────────────────────────────────────┘
```

### Images List Command

List the images of a category.

Recursive list:

```
$ piwigo-cli images list -r

Category1/SubCategory1/IMG_00001.jpeg
Category1/SubCategory1/IMG_00002.jpeg
Category1/SubCategory1/IMG_00003.jpeg
Category1/SubCategory1/IMG_00004.jpeg
Category1/SubCategory2/IMG_00005.jpeg
Category1/SubCategory2/IMG_00006.jpeg
Category2/SubCategory1/IMG_00007.jpeg
```

Specify a category:
```
$ piwigo-cli images list -r -c 2

Category2/SubCategory1/IMG_00007.jpeg
```

Tree view:

```
$ piwigo-cli images list -r -c 1 -t

.
├── SubCategory1
│   ├── IMG_00001.jpeg
│   ├── IMG_00002.jpeg
│   ├── IMG_00003.jpeg
│   └── IMG_00004.jpeg
└── SubCategory2
    ├── IMG_00005.jpeg
    └── IMG_00006.jpeg
```

### Images Details Command

Get details of an image. It supports the birthday plugin.

```
$ piwigo-cli images details -i 38062

┌───────────────┬───────────────────────────────────────────────────────────────────────────┐
│ KEY           │ VALUE                                                                     │
├───────────────┼───────────────────────────────────────────────────────────────────────────┤
│ Id            │ 38062                                                                     │
│ Md5           │ 6ad2abade6d5460181890e2bad671002                                          │
│ Name          │ 2006 04 14 015                                                            │
│ DateAvailable │ 2021-11-25 20:25:05                                                       │
│ DateCreation  │ 2006-04-14 04:14:00                                                       │
│ LastModified  │ 2022-01-01 23:11:48                                                       │
│ Width         │ 1984                                                                      │
│ Height        │ 1488                                                                      │
│ Url           │ https://yourphotos/picture.php?/38062                                     │
│ ImageUrl      │ https://yourphotos/upload/2021/11/25/20211125202505-6ad2abad.jpg          │
│ Filename      │ 2006_04_14_015.jpeg                                                       │
│ Filesize      │ 513                                                                       │
│ Categories    │ 2007                                                                      │
│ Tags          │ User Tag 1 (46 years old)                                                 │
│               │ User Tag 2 (8 months old)                                                 │
│               │ User Tag 3 (48 years old)                                                 │
└───────────────┴───────────────────────────────────────────────────────────────────────────┘
Derivatives:
┌─────────┬───────┬────────┬──────────────────────────────────────────────────────────────────────────────────────┐
│ NAME    │ WIDTH │ HEIGHT │ URL                                                                                  │
├─────────┼───────┼────────┼──────────────────────────────────────────────────────────────────────────────────────┤
│ thumb   │   144 │    108 │ https://yourphotos/i.php?/upload/2021/11/25/20211125202505-6ad2abad-th.jpg           │
│ xsmall  │   432 │    324 │ https://yourphotos/i.php?/upload/2021/11/25/20211125202505-6ad2abad-xs.jpg           │
│ xxlarge │  1656 │   1242 │ https://yourphotos/_data/i/upload/2021/11/25/20211125202505-6ad2abad-xx.jpg          │
│ square  │   120 │    120 │ https://yourphotos/_data/i/upload/2021/11/25/20211125202505-6ad2abad-sq.jpg          │
│ small   │   576 │    432 │ https://yourphotos/_data/i/upload/2021/11/25/20211125202505-6ad2abad-sm.jpg          │
│ medium  │   792 │    594 │ https://yourphotos/_data/i/upload/2021/11/25/20211125202505-6ad2abad-me.jpg          │
│ large   │  1008 │    756 │ https://yourphotos/i.php?/upload/2021/11/25/20211125202505-6ad2abad-la.jpg           │
│ xlarge  │  1224 │    918 │ https://yourphotos/_data/i/upload/2021/11/25/20211125202505-6ad2abad-xl.jpg          │
│ 2small  │   240 │    180 │ https://yourphotos/i.php?/upload/2021/11/25/20211125202505-6ad2abad-2s.jpg           │
└─────────┴───────┴────────┴──────────────────────────────────────────────────────────────────────────────────────┘
```

### Images Upload Command

Upload an image or a video in chunks. It will skip existing file on the server.

```
$ piwigo-cli images upload -f ~/Downloads/IMG_4886.jpeg -j 4
```

### Images Upload Tree Command

Upload a tree of images and videos in chunks, skipping existing file on the server.

```
$ piwigo-cli images upload-tree -d ~/Downloads/2021 -j4 -c 4
```

### Images Tag Command

Massive tag your image in your terminal.
```
$ piwigo-cli images tag -h
Usage:
  piwigo-cli [OPTIONS] images tag [tag-OPTIONS]

Help Options:
  -h, --help                      Show this help message

[tag command options]
      -i, --id=                   image id to tag
      -t, --tag-id=               look up for the first image of this tagId
      -T, --tag=                  look up for the first image of this tagName
      -x, --exclude=              exclude tag from selection
      -m, --max=                  loop on a maximum number of images (default: 1)
      -k, --keep                  keep survey filter
      -K, --keep-previous-answer  Preserve previous answer
```

Example:
Retag the image you mark as "todo:todo-2006", 50 max at a time, by keeping the previous selection between images.
```
$ piwigo-cli images tag -x ^todo -T todo:todo-2006 -m 50 -K
```

It display in well on iTerm:

- image
- some details with the previous tag
- list of tags selection

You can use "SPACE" for selection of a tag, "LEFT" to remove the current selection, type words to lookup.

![visutag](https://user-images.githubusercontent.com/65178/151510534-1c029d2b-f3a4-4fef-a4ae-8d831077a387.jpg)


## License

MIT
