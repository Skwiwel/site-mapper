# Site Mapper
A tool for creating map of sites. Given a site it scrapes for links and crawls further to create a directional graph of the website.

Currently the tool only displays the unique page adresses it found in no particular order. Proper graph data is internally acquired anyway, so the project may be expanded in the future with a cooler graph display.

Nothing too original here. It was mainly a project of mine aimed at better understanding the Golang's concurrency, http and unit test mechanisms.

## Usage
Site mapper currently offers a command line interface. To instruct the algorithm one can use the following flags:
- `--address` - required; the address of the site to map
- `--depth` - crawl depth. Represents the maximum "height" of the crawling tree. default: 2
- `--fast` - enables the fast mode of mapping. In this mode no javascript is parsed so only bare html is taken into consideration. This may lead to an incomplete mapping, but is 10s of times faster.

## Inner workings
For each link (`<a href=...>`) a new recursive goroutine is created with its depth lowered in comparison to the calling goroutine. The scraped pages are stored in a concurrency safe map.

In the fast mode the pages are fetched using the Golang http library. This is fast, but insufficient for the purposes of extensive scrapig of any modern website. For that an actual web browser must be simulated and to this end the mapper uses an actual chrome process in the default mode. With the use of [chromedp](https://github.com/chromedp/chromedp) one can create a headless chrome instance and utilize it to allow a website to run as designed. Chromedp is used in the default mode. I found chromedp somewhat unreliable and slow, but generally once it starts it gets the job done.

## Known issues
Will get fixed if I'll ever bother to work on this project some more.
- [ ] Once in about 4 runs the headless chrome instance doesn't start. The temporary fix is to just rerun the executable.
- [ ] Mapped pages are sometimes marked with an incorrect status code
- [ ] The program is currently not displaying the links acquired at the end of the search depth. This is a little inefficient, although it can show a status code for each link this way.
