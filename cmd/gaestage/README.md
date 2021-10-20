
# gaestage

The command `gcloud app deploy` stages all files it needs to a temporary folder. You can restrict what files will be uploaded through a `.gcloudignore` file, though that mechanism does not influences the initial copy. If you want to exclude big files, or an uncopiable `node_modules` folder gcloud will try to stage them anyway.

This tool will help you pre-stage all the files ignoring those in the `.gcloudignore` configuration. It comes at a cost of a double file copy when later gcloud does it again, but it only copies files that will be deployed later.


## Usage

Pass the App Engine command as extra arguments and the tool will run it in the correct folder after copying the files:

```shell
gaestage -- gcloud app deploy module.yml
```

You can also control the stage folder:

```shell
gaestage -s /tmp/myfolder
cd /tmp/myfolder
gcloud app deploy module.yml
```

There are more options you can explore with `gaestage -h`


## Official bug

There is an official bug filled in Google Issue Tracker that you can track: [173717530](https://issuetracker.google.com/issues/173717530)
