import sys
import re
import requests
from bs4 import BeautifulSoup
from urllib.parse import urljoin

def get_links(url):
    try:
        response = requests.get(url)
        response.raise_for_status()  # Raise HTTPError for bad resp.
        soup = BeautifulSoup(response.content, 'html.parser')
        links = []
        for a_tag in soup.find_all('a', href=True):
            absolute_url = urljoin(url, a_tag['href'])
            links.append(absolute_url)
        return list(set(links))  # Remove duplicate links
    except requests.exceptions.RequestException as e:
        print(f"Error fetching {url}: {e}")
        return []
    except Exception as e:
        print(f"Error parsing {url}: {e}")
        return []

if __name__ == "__main__":
    base_url = sys.argv[1]
    # Updated regex: capture only the domain
    match = re.search(r'^https?://([^/]+)', base_url)
    if match:
        domain = match.group(1)
        print(f"Domain extracted: {domain}")
    else:
        print("Invalid URL format. Using default filename 'output.txt'.")
        domain = "output"
    
    file_name = domain + ".txt"
    link_list = get_links(base_url)
    with open(file_name, "w") as f:
        for link in link_list:
            f.write(f"{link}\n")
    print(f"Extracted links saved to {file_name}")
