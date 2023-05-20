import re
import requests
from bs4 import BeautifulSoup
from urllib.parse import urljoin
import argparse
from pkg_resources import get_distribution

# ASCII Banner
banner = f"""
 _  __                _           _    
| |/ /               (_)         | |   
| ' / _ __   ___  ___ _  ___ _ __| | __
|  < | '_ \ / _ \/ _ | |/ _ | '__| |/ /
| . \| | | |  __|  __| |  __| |  |   < 
|_|\_|_| |_|\___|\___| |\___|_|  |_|\_\              
                    |__/                
                               v0.0.1
"""
print(banner)

# Pattern for .js files
js_file_pattern = re.compile(r'.*\.js')

# Regex to find environment variables in both formats
env_var_pattern = re.compile(r'(\b(?:NODE|REACT|AWS)[A-Z_]*\b\s*:\s*".*?")|(process\.env\.[A-Z_][A-Z0-9_]*)')


def scrape_js_files(url, found_vars=set(), debug=False):
    response = requests.get(url)

    # Parse the HTML content
    soup = BeautifulSoup(response.text, 'html.parser')

    # Find all script and link tags
    script_tags = soup.find_all(['script', 'link'])

    # Check each tag
    for tag in script_tags:
        src = tag.get('src') or tag.get('href')

        # Check if tag has a source and if it's under '/static/' and it's a .js file
        if src and '/static/' in src and js_file_pattern.match(src):
            js_url = urljoin(url, src)
            js_response = requests.get(js_url)

            # If the response is HTML, it may be a directory listing
            if 'html' in js_response.headers.get('Content-Type'):
                scrape_js_files(js_url, found_vars, debug)
            else:
                # Search for environment variables
                matches = env_var_pattern.findall(js_response.text)
                for match in matches:
                    match = match[0] if match[0] else match[1]  # Choose the match from the correct group
                    if match not in found_vars:
                        found_vars.add(match)
                        if debug:
                            print(f'[kneejerk] [js] [debug] {js_url} [{match}]')
                        else:
                            print(f'[kneejerk] [js] [info] {js_url} [{match}]')


def main():
    parser = argparse.ArgumentParser(description='Kneejerk - A tool for scanning environment variables in .js files')
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument('-u', '--url', help='URL of the website to scan')
    group.add_argument('-l', '--list', help='Path to a file containing a list of URLs to scan')
    parser.add_argument('-o', '--output', help='Path to output file')
    parser.add_argument('-debug', action='store_true', help='Print debugging statements')

    args = parser.parse_args()

    found_vars = set()
    if args.url:
        scrape_js_files(args.url, found_vars, args.debug)
    else:
        with open(args.list, 'r') as file:
            urls = file.readlines()
            for url in urls:
                url = url.strip()
                scrape_js_files(url, found_vars, args.debug)

    if args.output:
        with open(args.output, 'w') as f:
            for var in found_vars:
                f.write(f'[kneejerk] [js] [info] {var}\n')
            print(f'Results saved to {args.output}')


if __name__ == "__main__":
    main()
