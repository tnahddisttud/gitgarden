# gitgarden
Create a contribution graph of local git commits.

gitgarden is a cli tool lets you create github like commit graph for your local git projects. It utilizes the go-git package to to interact with git for fetching the commit details and is inspired by @flaviocope's work. The commit graph is displayed in gruvbox theme, it may support multiple themes in future.

### installation
---
To install `gitgarden` ensure that you have go installed in your computer. There are two ways to install it: 
1. Build it on your own using : `go build -o gitgarden ./main.go`

2. If you use linux, you can install it by running the following command:
    ```bash
    curl -sSL https://raw.githubusercontent.com/tnahddisttud/gitgarden/main/install.sh | bash
    ```

### usage
---
gitgarden supports two main operations (as of now): 

- add <path>: add a folder to be tracked by gitgarden for any git repos it contains.
- email <email>: specify the email of the user to be scanned for statistics.


To track the repos in a particular working directory, use the following command:
```bash
gitgarden -add "absolute/path/to/folder"
```


To print the commit graph on your terminal, use the following command:
```bash
gitgarden -email "your@github.email"
```
### future plans
---
Would love to implement the following features:

- adding an `-ignore` flag to ignore certain repos from being tracked
- supporting multiple themes
- adding aliases for emails

### thanks
Stay updated on my projects by connecting with me on X (formerly Twitter): https://x.com/tnahddisttud
