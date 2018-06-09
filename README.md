# CCollage

CCollage generates SVG collages of a Git repository's contributors.

## Usage

    ./ccollage --port 53352

**--port**: Port to listen on. Defaults to 8080.

## Environment Variables

**CCOLLAGE_GH_TOKEN**: A GitHub personal access token (get one [here](https://github.com/settings/tokens)). You'll only need the `public_repo` and `read:user` permissions. _Highly recommended (to increase [rate limits](https://developer.github.com/v3/#rate-limiting))._

**CCOLLAGE_GH_CACHE_SIZE**: The size of the GitHub API client in-memory cache, in megabytes. Defaults to 20MB. _Optional._

**CCOLLAGE_GH_CACHE_MAXAGE**: The maximum age of the GitHub API client cache items, in minutes. Defaults to 7 days. _Optional._

## Licence

Code is licensed under the GNU Affero General Public Licence version 3.