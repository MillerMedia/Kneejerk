# Kneejerk

Kneejerk is a Python-based tool for scanning environment variables from React websites.

## Features
* Scans JavaScript files of a provided URL for environment variables.
* Outputs found environment variables to the console or to a specified file.

## Usage
```angular2html
kneejerk -u https://www.example.com [-o output.txt]
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

## License

This project is licensed under the MIT License.