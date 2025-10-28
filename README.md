<p align="center">
  <a href="https://github.com/bitesinbyte">
    <img src="https://github.com/bitesinbyte/ferret/raw/main/.github/images/logo.png" width="256px" />
  </a>
</p>
<h1 align="center">Ferret</h1>

Automate the syndication of RSS feed posts to various social media platforms seamlessly with Ferret. Simplify your content distribution process and reach your audience effortlessly.

All updated docs and blog posts related to ferret can be found [here](https://blogs.bitesinbyte.com/categories/ferret/).

## Supported Social Sites

| Site       | Status    |
| ---------- | --------- |
| Mastodon   | Supported |
| LinkedIn   | Supported |
| Twitter(X) | Supported |
| Facebook   | Supported |
| Thread     | Supported |
| Instagram  | Supported |

## Configuration

Follow these steps to configure the project for your use:

1.  #### Fork the GitHub Repository

    To make changes and contribute to the project, fork the GitHub repository by following these steps:

    - Visit the GitHub repository you want to fork.
    - Click on the "Fork" button located at the top-right corner of the page.
    - Wait for the forking process to complete.
    - Once forked, you will have your copy of the repository in your GitHub account.

2.  #### Setting Up GitHub Secrets and Variables

    To securely store sensitive information and configure environment variables for your GitHub Actions workflow, follow these steps:

    - Setting Up GitHub Secrets:

      1. Visit your forked repository on GitHub.
      2. Go to the "Settings" tab.
      3. In the left sidebar, click on "Secrets".
      4. Click on "New repository secret".
      5. Add the following secrets:

      ```
         MASTODON_ACCESS_TOKEN
         LINKEDIN_ACCESS_TOKEN
         TWITTER_CONSUMER_KEY
         TWITTER_CONSUMER_SECRET
         TWITTER_ACCESS_TOKEN
         TWITTER_ACCESS_TOKEN_SECRET
         USER_EMAIL: GitHub User Email
         USER_NAME: GitHub User Name
         FACEBOOK_ACCESS_TOKEN
         THREAD_ACCESS_TOKEN
      ```

    - Setting Up GitHub Variables:

      Visit your forked repository on GitHub.

      1. Go to the "Settings" tab.
      2. In the left sidebar, click on "Secrets".
      3. Scroll down to the "Environment Variables" section.
      4. Add the following variables:

      ```
         MASTODON_INSTANCE_URL
         FACEBOOK_PAGE_ID
         THREAD_USER_ID
      ```

    For detailed instructions on how to add secrets and variables in GitHub, refer to the GitHub documentation: Creating and storing encrypted secrets.

3.  #### Updating config.json

    config.json contains the configurations regarding your connectors.

    ```json
    {
    "last_run_time": "2024-02-11T07:07:57.635936+01:00",
    "base_url": "https://blogs.bitesinbyte.com",
    "feed_endpoint": "/feed.xml",
    "does_meta_og_image_has_relative_path": true,
    "socials": [
        "facebook",
        "linkedin",
        "mastodon",
        "twitter",
        "thread"
    ]
    }
    ```
4. #### Enable workflow(GitHub Action)
    By default, GitHub disables workflows when you fork a repository. Since Ferret depends on GitHub Actions/Workflows, you can enable them by following these steps:
    - Navigate to the "Actions" tab.
    - Click "I understand my workflows, go ahead and enable them."
      ![workflow1](.github/images/workflow1.png)
    - Locate and enable the "Social Post" workflow.
      ![workflow2](.github/images/workflow2.png)


<span style='color: red;'>Note</span>

Ensure that you've provided the correct values for each secret and variable according to your setup. These configurations are necessary for the smooth functioning of the project and integration with external services.

## Local Development Steps

Follow these steps to set up and run the project locally on your machine:

Prerequisites
Make sure you have Go version 1.21.5 installed on your system. If not, follow these steps to download and install Go:

1. Visit the official Go website: <https://golang.org/dl/>
2. Download the installer for your operating system.
3. Follow the installation instructions provided on the website.

Setting up Environment Variables

Before running the application locally, ensure you have a .env file in the root directory of the project. This file should contain the following environment variables:

```
MASTODON_INSTANCE_URL=
MASTODON_ACCESS_TOKEN=
TWITTER_CONSUMER_KEY=
TWITTER_CONSUMER_SECRET=
TWITTER_ACCESS_TOKEN=
TWITTER_ACCESS_TOKEN_SECRET=
RSS_FEED_URL=
LINKEDIN_ACCESS_TOKEN=
FACEBOOK_ACCESS_TOKEN=
FACEBOOK_PAGE_ID=
THREAD_USER_ID=
THREAD_ACCESS_TOKEN=
IG_ACCESS_TOKEN=
IG_USER_ID=
IG_GRAPH_VERSION=
IG_FIRST_COMMENT_TEXT=
IG_TRIGGER_WORDS=
IG_MESSAGING_ENABLED=
```

Fill in the values for these variables according to your environment.

Building the Application
To build the application, execute the following command in your terminal:

```bash
go build -o bin/ferret ./cmd/ferret
```

This command will compile the application and generate an executable file named ferret inside the bin directory.

Running the Application
Once the application is built, you can run it using the following command:

```bash
./bin/ferret
```

This command will execute the compiled ferret binary and start the application locally.


<span style='color: red;'>Note</span>

Ensure all required environment variables are correctly set in the .env file before running the application.

## Instagram Automation

- Uses Instagram Graph API container + publish for feed posts, reels (with processing polling), and stories (permissions required).
- Optional first comment: set `IG_FIRST_COMMENT_TEXT` to auto-post right after publish.
- Trigger-to-DM: worker at `pkg/external/instagram/wokers/comments_dm.go` scans comments for words from `IG_TRIGGER_WORDS` and requests a DM message from your AI generator (`pkg/ai/generator`). DM sending requires Instagram Messaging; otherwise it’s skipped/unsupported.

Environment variables:

- `IG_ACCESS_TOKEN`, `IG_USER_ID`, `IG_GRAPH_VERSION` (default v19.0)
- `IG_FIRST_COMMENT_TEXT` – optional first comment content
- `IG_TRIGGER_WORDS` – comma-separated triggers, e.g., `READY,TEMPLATE`
- `IG_MESSAGING_ENABLED` – set `true` only if messaging is configured

## License

Licensed under the [MIT license](https://github.com/bitesinbyte/ferret/blob/main/LICENSE).
