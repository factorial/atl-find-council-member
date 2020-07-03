# Atlanta City Council Districts

An app to quickly find an Atlanta city council representative.

[Use it!](https://abrie.github.io/atl-find-council-member)

## Technical Overview

This app multiplexes two official Atlanta government sites into one:

- [egis.atlantaga.gov](http://egis.atlantaga.gov/app/home/index.html) to find a district by street address.
- [citycouncil.atlantaga.gov](https://citycouncil.atlantaga.gov/council-members) for information about council members such as full name, photo, and contact info.

## Technical Details

There are three components:

- [/web](web), javascript web interface.
- [/backend](backend), wraps the the egis AP and serves scraped council data.
- [/scraper](scraper), gleans information from the city council web site.

## For Developers

### Frontend development:

The app uses Typescript, Tailwind CSS, JSX, and Snowpack.

1. Fork this repo: [click here](https://github.com/abrie/atl-find-council-member/fork)
2. Clone the fork onto your development machine.
3. Install dependencies: `yarn install`
4. Change to front end's folder: `cd web`
5. Start the devserver: `yarn start`

To contribute changes:

1. Add this repo as a remote: `git remote add upstream https://github.com/abrie/atl-find-council-member`
2. Follow the [Github guide for syncing and merging](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/working-with-forks) to make a PR.

### Backend development:

The backend uses Go and the Chi router framework.

_(todo)_

### Scraper development:

The scraper uses Python3 with Beautifulsoup.

_(todo)_
