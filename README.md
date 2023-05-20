# Kneejerk

Kneejerk is a pentesting command line tool for scanning environment variables from React websites.

## Features
* Scans JavaScript files of a provided URL for environment variables.
* Outputs found environment variables to the console or to a specified file.

## Usage
```angular2html
kneejerk -u https://www.example.com -o output.txt
```

```angular2html
Kneejerk - A tool for scanning environment variables in .js files

optional arguments:
  -h, --help            show this help message and exit
  -u URL, --url URL     URL of the website to scan
  -l LIST, --list LIST  Path to a file containing a list of URLs to scan
  -o OUTPUT, --output OUTPUT
                        Path to output file
  -debug                Print debugging statements
```
## Installation

### Homebrew (Recommended)

You can install Kneejerk using Homebrew:

```bash
brew tap MillerMedia/kneejerk
brew install kneejerk
```

### From Source 

Alternatively, you can install Kneejerk by cloning this repository and running setup.py:

```bash
git clone https://github.com/MillerMedia/Kneejerk
cd Kneejerk
python setup.py install
```

## Contributing

I welcome contributions from the community! If you have any suggestions, bug reports, or ideas for improvement, feel free to open an issue or submit a pull request.

## Support the project

If you find this project helpful and would like to support its development, please consider donating:  
  
[![Buy me a coffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/yOd1JU9MQe)

## License

This project is licensed under the MIT License.