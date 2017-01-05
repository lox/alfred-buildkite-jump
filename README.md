# Alfred Buildkite Jump

A [workflow for Alfred 3](https://www.alfredapp.com/help/workflows/) for indexing your Buildkite pipelines and quickly opening them in your default browser.

![](https://lachlan.me/s/WVFrDoEw.png)

## Installation

1. Generate a [Buildkite API Token](https://buildkite.com/user/api-access-tokens) with the `read_organisations` and `read_pipelines` permission.

2. Install the latest version from the [releases page on Github](https://github.com/lox/alfred-buildkite-jump/releases). The configuration screen will open, click the [configure workflow and variables](https://lachlan.me/s/w768Cvri.png) `[X]` icon in the top right hand corner.
 
3. Add an environmental variable called `BUILDKITE_API_TOKEN` with the value from step 1. 

4. Open your Alfred prompt, type `bk` and select `Update Buildkite Pipelines`. This will take a few seconds. You can repeat this whenever you add a new pipeline with `bk > update`.

5. Now you can type `bk ...` and it will autocomplete to your pipelines. 





