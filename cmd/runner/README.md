
# runner

Publish Cloud Run applications.


## Build new containers

Build a new container with the application `myname`:

```shell
runner build myname --project google-project-foo
```

Inside our normal Jenkins scripts where a variable is defined to configure gcloud previously:

```shell
runner build myname
```

Dockerfile must be organized inside a folder with the name of the application: `myname/Dockerfile`. Container will build from the directory where this application runs to allow cross-applications package imports.

You can build multiple containers at the same time:

```shell
runner build foo bar baz --project $GOOGLE_PROJECT
```


## Deploy to Cloud Run

Generic execution in any environment:

```shell
runner deploy myname --project google-project-foo --sentry foo-name
```

Inside our normal Jenkins scripts where a variable is defined to configure gcloud previously:

```shell
runner deploy myname --sentry foo-name
```

You can deploy multiple containers at the same time:

```shell
runner deploy foo bar baz --project $GOOGLE_PROJECT --sentry foo-name
```
