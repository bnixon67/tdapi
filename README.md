# tdapi

tdapi provides a Go interface for the [Todoist REST API](https://developer.todoist.com/rest/v1/).

**This package and any related application is not created by, affiliated with, or supported by Doist.**

**This is still a work in progress, but does have some working examples.**

In order to use this package, you must register an application in the [Todoist App Management console](https://developer.todoist.com/appconsole.html). Use `https://example.com/redirect` as the **OAuth redirect URL** or adjust the source code to match the URL you entered.

In order to run the examples, you must set the TDCLIENTID and TDCLIENTSECRET environmental variables to the **Client ID** and **Client secret** from the app registration.

The current approach assumes the client runs on a host without a browser.

1. The user is instructed to vist a URL to login and authorize the client.

2. Once the login is successful, the user must copy the response URL and provide to the client program.

   Note that the default redirect URL `https://example.com/redirect` use the special-use example domain and the browser will display a generic message. You will need to copy the generated URL from the browser location bar into the command line program. The url should look something like `https://example.com/redirect?state={characters}&code={characters}`. The state and code are used to complete the OAuth access token exchange process.

By default, a .token.json file is created to store the OAuth2 Access Bearer token.

If you don't want to use OAuth and register an application, you can manually create a .json.token file, which looks like the json file below, with your API token from https://todoist.com/prefs/integrations:
```json
{"access_token":"Replace with your API token","token_type":"Bearer","expiry":"0001-01-01T00:00:00Z"}
```

**DO NOT SHARE YOUR API TOKEN OR THE TOKEN FILE CREATED BY THE APP. ANYONE WITH THE TOKEN HAS ACCESS TO YOUR TODOIST ACCOUNT.**
