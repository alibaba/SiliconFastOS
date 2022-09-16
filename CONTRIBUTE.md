# How to contribute

You can open issues for bugs you've found or features you think are missing. You can also submit pull requests to this repository.
You can contribute in the following ways:

- Finding and reporting bugs
- Contributing code to SiliconFastOS by fixing bugs or implementing features
- Improving the documentation

## Submit a patch

1. You need to start by opening a new issue describing the bug or feature you're intending to fix. Even if you think it's relatively minor, it's helpful to know what people are working on. Please use the search function to make sure that you are not submitting duplicates, and that a similar report or request has not already been resolved or rejected.

2. Forking the project, and setup a new branch to work in. It's important that each group of changes be done in separate branches in order to ensure that a pull request only includes the commits related to that bug or feature.

3. Do your best to have formatted commit messages for each change. This provides consistency throughout the project, and ensures that commit messages are able to be formatted properly by various git tools.

4. Finally, push the commits to your fork and submit a pull request. Please use clean, concise titles for your pull requests. Splitting tasks into multiple smaller pull requests is often preferable.

## Commit message rules

Please use Aug to edit your commit message, so we can use some git tools for project.

    <type>: [optional scope]: <description>

    [optional body]

    [optional footer(s)]

The type can be :

- build
- chore
- docs
- feat
- fix
- perf
- refactor
- test

### example

    chore: some chore

    fix(package): fix bug

    feat: add new feature
